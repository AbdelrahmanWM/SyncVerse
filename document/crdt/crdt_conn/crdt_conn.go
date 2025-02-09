package crdtconn

import (
	"fmt"
	"sync"

	. "github.com/AbdelrahmanWM/SyncVerse/document/crdt/crdt_conn/internal/utils"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/crdt_conn/webrtc_peer"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/global"
)

// type definitions
type PeerStatus string

const (
	SignalingInitiated PeerStatus = "SIGNALING_INITIATED"
)

// type Peer interface { // later

// }

type PeerEvent struct {
	PeerID global.ReplicaID
	State  PeerStatus
}

type PeerEvents chan PeerEvent

type PeerConnectionManager struct {
	peerEvents   chan PeerEvent
	wg           *sync.WaitGroup
	mux          *sync.Mutex
	peerStatuses *sync.Map
	peerMap      map[global.ReplicaID]*webrtc_peer.WebRTCPeer
}

// package lvl variables
var peerConnectionManager *PeerConnectionManager // singleton instance shared across all peers

// function definitions

func GetPeerConnectionManager() *PeerConnectionManager {
	if peerConnectionManager == nil {
		peerConnectionManager = NewPeerConnectionManager()
		go peerConnectionManager.peerEventListener()
	}
	return peerConnectionManager
}

func NewPeerConnectionManager() *PeerConnectionManager {
	var peerEvents = make(chan PeerEvent)
	var wg sync.WaitGroup
	var mux sync.Mutex
	var peerStatuses sync.Map
	var peerMap = make(map[global.ReplicaID]*webrtc_peer.WebRTCPeer)

	return &PeerConnectionManager{peerEvents, &wg, &mux, &peerStatuses, peerMap}
}

func (pcm *PeerConnectionManager) AddNewPeer(replicaID global.ReplicaID) { //when a user joins the session
	peer := webrtc_peer.NewWebRTCPeer(replicaID)
	pcm.addPeerToPeerMap(replicaID, peer)
	Log("new webrtc peer")
	pcm.peerEvents <- PeerEvent{PeerID: replicaID, State: SignalingInitiated}
	Log("Event added")
}
func (pcm *PeerConnectionManager) addPeerToPeerMap(replicaID global.ReplicaID, peer *webrtc_peer.WebRTCPeer) {
	pcm.mux.Lock()
	defer pcm.mux.Unlock()
	pcm.peerMap[replicaID] = peer
}
func (pcm *PeerConnectionManager) peerEventListener() {
	for event := range pcm.peerEvents {
		pcm.peerStatuses.Store(event.PeerID, event.State)
		Log(fmt.Sprintf("Peer %s moved to state %s:", event.PeerID, event.State))
		pcm.handlePeerEvent(event)
	}
}
func (pcm *PeerConnectionManager) handlePeerEvent(event PeerEvent) {
	switch event.State {
	case SignalingInitiated:
		Log(fmt.Sprintf("Peer %s have initialized signaling", event.PeerID))
	}
}
