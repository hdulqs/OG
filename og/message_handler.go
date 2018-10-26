package og

import (
	"github.com/annchain/OG/common"
	"github.com/annchain/OG/og/downloader"
	"github.com/annchain/OG/types"
	"github.com/sirupsen/logrus"
)

// IncomingMessageHandler is the default handler of all incoming messages for OG
type IncomingMessageHandler struct {
	Og  *Og
	Hub *Hub
}

func (h *IncomingMessageHandler) HandleFetchByHashRequest(syncRequest types.MessageSyncRequest, peerId string) {
	logrus.WithField("q", syncRequest.String()).Debug("received MessageSyncRequest")

	if len(syncRequest.Hashes) == 0 {
		logrus.Debug("empty MessageSyncRequest")
		return
	}

	var txs []*types.Tx
	var seqs []*types.Sequencer

	for _, hash := range syncRequest.Hashes {
		txi := h.Og.TxPool.Get(hash)
		if txi == nil {
			txi = h.Og.Dag.GetTx(hash)
		}
		switch tx := txi.(type) {
		case *types.Sequencer:
			seqs = append(seqs, tx)
		case *types.Tx:
			txs = append(txs, tx)
		}

	}
	syncResponse := types.MessageSyncResponse{
		Txs:        txs,
		Sequencers: seqs,
	}
	data, err := syncResponse.MarshalMsg(nil)
	if err != nil {
		logrus.Error("failed to marshall MessageSyncResponse message")
		return
	}

	logrus.WithField("p", syncResponse.String()).Debug("sending MessageSyncResponse")
	h.Hub.SendToPeer(peerId, MessageTypeFetchByHashResponse, data)
}

func (h *IncomingMessageHandler) HandleHeaderResponse(headerMsg types.MessageHeaderResponse, peerId string) {
	logrus.WithField("q", headerMsg).Debug("received MessageHeaderResponse")
	headers := headerMsg.Sequencers
	// Filter out any explicitly requested headers, deliver the rest to the downloader
	seqHeaders := types.SeqsToHeaders(headers)
	filter := len(seqHeaders) == 1

	// TODO: verify fetcher
	//if filter {
	//	// Irrelevant of the fork checks, send the header to the fetcher just in case
	//	seqHeaders = h.Og.fetcher.FilterHeaders(peerId, seqHeaders, time.Now())
	//}
	if len(seqHeaders) > 0 || !filter {
		err := h.Hub.Downloader.DeliverHeaders(peerId, seqHeaders)
		if err != nil {
			logrus.Debug("Failed to deliver headers", "err", err)
		}
	}
	logrus.WithField("header lens", len(seqHeaders)).Debug("heandle MessageTypeHeaderResponse")
}

func (h *IncomingMessageHandler) HandleHeaderRequest(query types.MessageHeaderRequest, peerId string) {
	logrus.WithField("q", query).Debug("received MessageHeaderRequest")

	hashMode := !query.Origin.Hash.Empty()
	first := true
	logrus.WithField("hash", query.Origin.Hash).WithField("number", query.Origin.Number).WithField(
		"hashmode", hashMode).WithField("amount", query.Amount).WithField("skip", query.Skip).Debug("requests")
	// Gather headers until the fetch or network limits is reached
	var (
		bytes   common.StorageSize
		headers []*types.Sequencer
		unknown bool
	)
	for !unknown && len(headers) < int(query.Amount) && bytes < softResponseLimit && len(headers) < downloader.MaxHeaderFetch {
		// Retrieve the next header satisfying the query
		var origin *types.Sequencer
		if hashMode {
			if first {
				first = false
				origin = h.Og.Dag.GetSequencerByHash(query.Origin.Hash)
				if origin != nil {
					query.Origin.Number = origin.Number()
				}
			} else {
				origin = h.Og.Dag.GetSequencer(query.Origin.Hash, query.Origin.Number)
			}
		} else {
			origin = h.Og.Dag.GetSequencerById(query.Origin.Number)
		}
		if origin == nil {
			break
		}
		headers = append(headers, origin)
		bytes += estHeaderRlpSize

		// Advance to the next header of the query
		switch {
		case hashMode && query.Reverse:
			// Hash based traversal towards the genesis block
			ancestor := query.Skip + 1
			if ancestor == 0 {
				unknown = true
			} else {
				seq := h.Og.Dag.GetSequencerById(query.Origin.Number - ancestor)
				query.Origin.Hash, query.Origin.Number = seq.GetTxHash(), seq.Number()
				unknown = query.Origin.Hash.Empty()
			}
		case hashMode && !query.Reverse:
			// Hash based traversal towards the leaf block
			var (
				current = origin.Number()
				next    = current + query.Skip + 1
			)
			if next <= current {
				logrus.Warn("GetBlockHeaders skip overflow attack", "current", current, "skip", query.Skip, "next", next, "attacker", peerId)
				unknown = true
			} else {
				if header := h.Og.Dag.GetSequencerById(next); header != nil {
					nextHash := header.GetTxHash()
					oldSeq := h.Og.Dag.GetSequencerById(next - (query.Skip + 1))
					expOldHash := oldSeq.GetTxHash()
					if expOldHash == query.Origin.Hash {
						query.Origin.Hash, query.Origin.Number = nextHash, next
					} else {
						unknown = true
					}
				} else {
					unknown = true
				}
			}
		case query.Reverse:
			// Number based traversal towards the genesis block
			if query.Origin.Number >= query.Skip+1 {
				query.Origin.Number -= query.Skip + 1
			} else {
				unknown = true
			}

		case !query.Reverse:
			// Number based traversal towards the leaf block
			query.Origin.Number += query.Skip + 1
		}
	}

	msgRes := types.MessageHeaderResponse{
		Sequencers: headers,
	}
	data, _ := msgRes.MarshalMsg(nil)
	logrus.WithField("len ", len(msgRes.Sequencers)).Debug("send MessageTypeGetHeader")

	h.Hub.SendToPeer(peerId, MessageTypeHeaderResponse, data)
}

func (h *IncomingMessageHandler) HandleTxsResponse(request types.MessageTxsResponse) {
	logrus.Debug("received MessageTxsResponse")
	if request.Sequencer != nil {
		logrus.WithField("len", len(request.Txs)).WithField("seq id", request.Sequencer.Id).Debug("got response txs ")
	} else {
		logrus.Warn("got nil sequencer")
		return
	}

	lseq := h.Og.Dag.LatestSequencer()
	//todo need more condition
	if lseq.Number() < request.Sequencer.Number() {
		h.Og.TxBuffer.AddRemoteTxs(request.Sequencer, request.Txs)
	}
	return
}

func (h *IncomingMessageHandler) HandleTxsRequest(msgReq types.MessageTxsRequest, peerId string) {
	logrus.WithField("q", msgReq).Debug("received MessageTxsRequest")

	var msgRes types.MessageTxsResponse

	var seq *types.Sequencer
	if msgReq.SeqHash != nil && msgReq.Id != 0 {
		seq = h.Og.Dag.GetSequencer(*msgReq.SeqHash, msgReq.Id)
	} else {
		seq = h.Og.Dag.GetSequencerById(msgReq.Id)
	}
	msgRes.Sequencer = seq
	if seq != nil {
		msgRes.Txs = h.Og.Dag.GetTxsByNumber(seq.Id)
	} else {
		logrus.WithField("id ", msgReq.Id).WithField("hash", msgReq.SeqHash).Warn("seq was not found for request ")
	}
	data, _ := msgRes.MarshalMsg(nil)
	logrus.WithField("txs num ", len(msgRes.Txs)).Debug("send MessageTypeGetTxs")
	h.Hub.SendToPeer(peerId, MessageTypeTxsResponse, data)
}

func (h *IncomingMessageHandler) HandleBodiesResponse(request types.MessageBodiesResponse, peerId string) {
	logrus.WithField("q", request).Debug("received MessageBodiesResponse")

	// Deliver them all to the downloader for queuing
	transactions := make([][]*types.Tx, len(request.Bodies))
	sequencers := make([]*types.Sequencer, len(request.Bodies))
	for i, bodyData := range request.Bodies {
		var body types.MessageTxsResponse
		_, err := body.UnmarshalMsg(bodyData)
		if err != nil {
			logrus.WithError(err).Warn("decode error")
			break
		}
		if body.Sequencer == nil {
			logrus.Warn(" body.Sequencer is nil")
			break
		}
		transactions[i] = body.Txs
		sequencers[i] = body.Sequencer
	}
	logrus.WithField("bodies len", len(request.Bodies)).Debug("got bodies")

	// Filter out any explicitly requested bodies, deliver the rest to the downloader
	filter := len(transactions) > 0 || len(sequencers) > 0
	// TODO: verify fetcher
	//if filter {
	//	transactions = h.Og.fetcher.FilterBodies(peerId, transactions, sequencers, time.Now())
	//}
	if len(transactions) > 0 || len(sequencers) > 0 || !filter {
		logrus.WithField("txs len", len(transactions)).WithField("seq len", len(sequencers)).Debug("deliver bodies ")
		err := h.Hub.Downloader.DeliverBodies(peerId, transactions, sequencers)
		if err != nil {
			logrus.Debug("Failed to deliver bodies", "err", err)
		}
	}
	logrus.Debug("handle MessageTypeBodiesResponse")
	return
}

func (h *IncomingMessageHandler) HandleBodiesRequest(msgReq types.MessageBodiesRequest, peerId string) {
	logrus.WithField("q", msgReq).Debug("received MessageBodiesRequest")
	var msgRes types.MessageBodiesResponse
	var bytes int

	for i := 0; i < len(msgReq.SeqHashes); i++ {
		seq := h.Og.Dag.GetSequencerByHash(msgReq.SeqHashes[i])
		if seq == nil {
			logrus.Warn("seq is nil")
			break
		}
		if bytes >= softResponseLimit {
			logrus.Debug("reached softResponseLimit ")
			break
		}
		if len(msgRes.Bodies) >= downloader.MaxBlockFetch {
			logrus.Debug("reached MaxBlockFetch 128 ")
			break
		}
		var body types.MessageTxsResponse
		body.Sequencer = seq
		body.Txs = h.Og.Dag.GetTxsByNumber(seq.Id)
		bodyData, _ := body.MarshalMsg(nil)
		bytes += len(bodyData)
		msgRes.Bodies = append(msgRes.Bodies, types.RawData(bodyData))
	}
	data, _ := msgRes.MarshalMsg(nil)
	logrus.WithField("bodies num ", len(msgRes.Bodies)).Debug("send MessageTypeBodiesResponse")
	h.Hub.SendToPeer(peerId, MessageTypeBodiesResponse, data)
}

func (h *IncomingMessageHandler) HandleSequencerHeader(msgHeader types.MessageSequencerHeader, peerId string) {
	if msgHeader.Hash == nil {
		return
	}

	//no need to broadcast again ,just all our peers need know this ,not all network
	//set peer's head
	h.Hub.SetPeerHead(peerId, *msgHeader.Hash, msgHeader.Number)

	//if h.SyncManager.Status != syncer.SyncStatusIncremental{
	//	return
	//}
	lseq := h.Og.Dag.LatestSequencer()
	for i := lseq.Number(); i < msgHeader.Number; i++ {
		go func(i uint64) {
			//p.RequestTxsById(i + 1)
		}(i)
	}
	return
}

func (h *IncomingMessageHandler) HandlePing(peerId string) {
	logrus.Debug("received your ping. Respond you a pong")
	h.Hub.SendToPeer(peerId, MessageTypePong, []byte{1})
}

func (h *IncomingMessageHandler) HandlePong() {
	logrus.Debug("received your pong.")
}

func (h *IncomingMessageHandler) HandleFetchByHashResponse(syncResponse types.MessageSyncResponse, sourceId string) {
	logrus.WithField("q", syncResponse.String()).Debug("received MessageSyncResponse")
	if (syncResponse.Txs == nil || len(syncResponse.Txs) == 0) &&
		(syncResponse.Sequencers == nil || len(syncResponse.Sequencers) == 0) {
		logrus.Debug("empty MessageSyncResponse")
		return
	}

	for _, v := range syncResponse.Txs {
		logrus.WithField("tx", v).WithField("peer", sourceId).Debugf("received sync response Tx")
		h.Og.TxBuffer.AddRemoteTx(v)
	}
	for _, v := range syncResponse.Sequencers {
		logrus.WithField("seq", v).WithField("peer", sourceId).Debugf("received sync response seq")
		h.Og.TxBuffer.AddRemoteTx(v)
	}
}