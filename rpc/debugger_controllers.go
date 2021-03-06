package rpc

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type SyncStatus struct {
	Id                       string `json:"id"`
	SyncMode                 string `json:"syncMode"`
	CatchupSyncerStatus      string `json:"catchupSyncerStatus"`
	CatchupSyncerEnabled     bool   `json:"catchupSyncerEnabled"`
	IncrementalSyncerEnabled bool   `json:"incrementalSyncerEnabled"`
	Height                   uint64 `json:"height"`
	LatestHeight             uint64 `json:"latestHeight"`
	BestPeer                 string `json:"bestPeer"`
	Error                    string `json:"error"`
	Txid                     uint32 `json:"txid"`
}

//Status node status
func (r *RpcController) SyncStatus(c *gin.Context) {
	status := r.syncStatus()
	cors(c)
	c.JSON(http.StatusOK, status)
}

func (r *RpcController) syncStatus() SyncStatus {
	var status SyncStatus

	status = SyncStatus{
		Id:                       r.P2pServer.Self().ID.TerminalString(),
		SyncMode:                 r.SyncerManager.Status.String(),
		CatchupSyncerStatus:      r.SyncerManager.CatchupSyncer.WorkState.String(),
		CatchupSyncerEnabled:     r.SyncerManager.CatchupSyncer.Enabled,
		IncrementalSyncerEnabled: r.SyncerManager.IncrementalSyncer.Enabled,
		Height: r.SyncerManager.NodeStatusDataProvider.GetCurrentNodeStatus().CurrentId,
	}

	peerId, _, seqId, err := r.SyncerManager.Hub.BestPeerInfo()
	if err != nil {
		status.Error = err.Error()
	} else {
		status.LatestHeight = seqId
		status.BestPeer = peerId

	}
	status.Txid = r.Og.TxPool.GetTxNum()
	return status
}

func (r *RpcController) Performance(c *gin.Context) {
	cd := r.PerformanceMonitor.CollectData()
	cors(c)
	c.JSON(http.StatusOK, cd)
}
