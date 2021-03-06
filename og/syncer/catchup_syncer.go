package syncer

import (
	"sync"
	"time"

	// "github.com/annchain/OG/ffchan"
	"github.com/annchain/OG/og"
	"github.com/annchain/OG/og/downloader"
	"github.com/annchain/OG/types"
	"github.com/sirupsen/logrus"
)

const (
	//forceSyncCycle      = 10 * time.Second // Time interval to force syncs, even if few peers are available
	minDesiredPeerCount = 5 // Amount of peers desired to start syncing

	// This is the target size for the packs of transactions sent by pushPendingTxLoop.
	// A pack can get larger than this if a single transactions exceeds this size.
	txsyncPackSize = 100 * 1024

	// TODO: this value will be set to optimal value in the future.
	// If generating sequencer is very fast with few transactions, it should be bigger,
	// otherwise it should be smaller
	SyncerCheckTime = time.Second * 6

	// when to stop sync once started
	stopSyncHeightDiffThreashold uint64 = 0
	// when to start sync
	startSyncHeightDiffThreashold uint64 = 2
)

type CatchupSyncerStatus int

func (m CatchupSyncerStatus) String() string {
	switch m {
	case Started:
		return "CSSStarted"
	case Stopped:
		return "CSSStopped"
	default:
		return "Unknown"
	}
}

const (
	Started CatchupSyncerStatus = iota
	Stopped
)

type CatchupSyncer struct {
	NodeStatusDataProvider og.NodeStatusDataProvider
	PeerProvider           og.PeerProvider
	Hub                    *og.Hub

	Downloader *downloader.Downloader
	SyncMode   downloader.SyncMode

	// should be enabled until quit
	EnableEvent chan bool
	Enabled     bool

	quitLoopEvent chan bool
	quit          chan bool

	OnWorkingStateChanged         []chan CatchupSyncerStatus
	OnNewTxiReceived              []chan types.Txi
	NewPeerConnectedEventListener chan string
	syncFlag                      bool
	WorkState                     CatchupSyncerStatus
	mu                            sync.RWMutex
	BootStrapNode                 bool
	currentBestHeight             uint64 //maybe incorect
}

func (c *CatchupSyncer) Init() {
	c.EnableEvent = make(chan bool)
	c.quitLoopEvent = make(chan bool)
	c.NewPeerConnectedEventListener = make(chan string)
	c.quit = make(chan bool)
}

func (c *CatchupSyncer) Start() {
	go c.eventLoop()
	go c.loopSync()
}

func (c *CatchupSyncer) Stop() {
	c.quit <- true
	c.quitLoopEvent <- true
	// <-ffchan.NewTimeoutSender(c.quit, true, "catchupSyncerQuit", 1000).C
	// <-ffchan.NewTimeoutSender(c.quitLoopEvent, true, "catchupSyncerQuitLoopEvent", 1000).C
}

func (CatchupSyncer) Name() string {
	return "CatchupSyncer"
}

func (c *CatchupSyncer) isUpToDate(maxDiff uint64) bool {
	_, bpHash, seqId, err := c.PeerProvider.BestPeerInfo()
	if err != nil {
		logrus.WithError(err).Debug("get best peer")
		if c.BootStrapNode {
			return true
		}
		return false
	}
	ourId := c.NodeStatusDataProvider.GetCurrentNodeStatus().CurrentId
	logrus.WithField("bestPeer SeqId", seqId).
		WithField("bestPeerHash", bpHash).
		WithField("our SeqId", ourId).
		Trace("checking uptodate")

	if seqId <= ourId+maxDiff {
		log.WithField("bestPeer SeqId", seqId).
			WithField("bestPeerHash", bpHash).
			WithField("our SeqId", ourId).
			Debug("we are now up to date")
		return true
	} else {
		logrus.WithField("bestPeer SeqId", seqId).
			WithField("bestPeerHash", bpHash).
			WithField("our SeqId", ourId).
			Debug("we are yet up to date")
	}
	return false
}

func (c *CatchupSyncer) loopSync() {
	defer c.Downloader.Terminate()
	for {
		select {
		case <-c.quit:
			log.Info("CatchupSyncer loopSync received quit message. Quitting...")
			return
		case peer := <-c.NewPeerConnectedEventListener:
			log.WithField("peer", peer).Info("new peer connected")
			if !c.Enabled {

				log.Debug("catchupSyncer not enabled")
				continue
			}
			go c.syncToLatest()
		case <-time.After(SyncerCheckTime):
			if !c.Enabled {
				log.Debug("catchup syncer not enabled")
				continue
			}
			go c.syncToLatest()
		}

	}
}

func (c *CatchupSyncer) isSyncing() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.syncFlag
}

//getWorkState
func (c *CatchupSyncer) getWorkState() CatchupSyncerStatus {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.WorkState
}

func (c *CatchupSyncer) setSyncFlag() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.syncFlag = true
}

func (c *CatchupSyncer) unsetSyncFlag() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.syncFlag = false
}

func (c *CatchupSyncer) CacheNewTxEnabled() bool {
	if c.getWorkState() == Stopped {
		return true
	}
	ourSeqId := c.NodeStatusDataProvider.GetHeight()
	if ourSeqId+startSyncHeightDiffThreashold*3 < c.currentBestHeight {
		return false
	}
	return true

}

func (c *CatchupSyncer) syncToLatest() error {
	if c.isSyncing() {
		log.Trace("catchup syncing task is busy")
		return nil
	}
	c.setSyncFlag()
	defer c.unsetSyncFlag()
	//get best peer ,and sync with this peer until we catchup

	// must sync to latest at the beginning
	var diff uint64 = 0
	for !c.isUpToDate(diff) {
		c.NotifyWorkingStateChanged(Started)
		diff = stopSyncHeightDiffThreashold

		bpId, bpHash, seqId, err := c.PeerProvider.BestPeerInfo()
		if err != nil {
			log.WithError(err).Warn("picking up best peer")
			return err
		}
		ourId := c.NodeStatusDataProvider.GetCurrentNodeStatus().CurrentId

		log.WithField("peerId", bpId).WithField("seq", seqId).WithField("ourId", ourId).
			Debug("catchup sync with best peer")
		c.currentBestHeight = seqId
		// Run the sync cycle, and disable fast sync if we've went past the pivot block
		if err := c.Downloader.Synchronise(bpId, bpHash, seqId, c.SyncMode); err != nil {
			log.WithError(err).Warn("catchup sync failed")
			return err
		}
		logrus.WithField("seqId", seqId).Debug("finished downloader synchronize")
		//bpHash, seqId, err = c.PeerProvider.GetPeerHead(bpId)
		//if err != nil {
		//	logrus.WithError(err).Warn("sync failed")
		//	return err
		//}
	}
	c.NotifyWorkingStateChanged(Stopped)
	// allow a maximum of startSyncHeightDiffThreashold behind
	diff = startSyncHeightDiffThreashold
	return nil
}

//we don't need to broadcast this ,because we broadcast all of our latest sequencer head when it change
/*
func (c *CatchupSyncer) notifyProgress() {
	nodeStatus := c.NodeStatusDataProvider.GetCurrentNodeStatus()
	if nodeStatus.CurrentId > 0 {
		// We've completed a sync cycle, notify all peers of new state. This path is
		// essential in star-topology networks where a gateway node needs to notify
		// all its out-of-date peers of the availability of a new block. This failure
		// scenario will most often crop up in private and hackathon networks with
		// degenerate connectivity, but it should be healthy for the mainnet too to
		// more reliably update peers or the local TD state.
		//go h.BroadcastBlock(head, false)
		hash := nodeStatus.CurrentBlock
		msg := types.MessageSequencerHeader{Hash: &hash, Number: nodeStatus.CurrentId}
		data, _ := msg.MarshalMsg(nil)
		c.Hub.BroadcastMessage(og.MessageTypeSequencerHeader, data)
	}
}

*/
func (c *CatchupSyncer) eventLoop() {
	for {
		select {
		case v := <-c.EnableEvent:
			log.WithField("enable", v).Info("catchup syncer got enable event ")
			c.Enabled = v
		case <-c.quitLoopEvent:
			log.Debug("catchup syncer eventLoop received quit message. Quitting...")
			return
		}
	}
}

//NotifyWorkingStateChanged if starts status is true ,stops status is  false
func (c *CatchupSyncer) NotifyWorkingStateChanged(status CatchupSyncerStatus) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.WorkState == status {
		return
	}
	c.WorkState = status
	for _, ch := range c.OnWorkingStateChanged {
		ch <- status
		// <-ffchan.NewTimeoutSender(ch, status, "NotifyWorkingStateChanged", 1000).C
	}
}
