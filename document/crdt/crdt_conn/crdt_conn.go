//go:build js && wasm

// +build: js,wasm
package crdtconn

import (
	"sync"

	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/crdt_conn/webrtc_peer"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/global"
)
type PeerStatus string

const (
	SignalingInitiated PeerStatus = "SIGNALING_INITIATED"
)

type Peer interface {
	
}
type PeerEvent struct {
	PeerID string
	State  PeerStatus
}



type PeerEvents chan PeerEvent

type PeerConnectionManager struct {
	peerEvents chan PeerEvent
	wg         *sync.WaitGroup
	peerStatuses *sync.Map
	peerMap map[string]*webrtc_peer.WebRTCPeer
}

func NewPeerConnectionManager()*PeerConnectionManager{
	var peerEvents = make(chan PeerEvent)
	var wg sync.WaitGroup
	var peerStatuses sync.Map
	var peerMap = make(map[string]*webrtc_peer.WebRTCPeer)

	return &PeerConnectionManager{peerEvents,&wg,&peerStatuses,peerMap}
}

func (pcm *PeerConnectionManager)AddPeer(replicaID global.ReplicaID){
	peer:=webrtc_peer.NewWebRTCPeer(replicaID)
	pcm.peerMap[string(replicaID)]=peer
	pcm.peerEvents<-PeerEvent{PeerID: string(replicaID),State:SignalingInitiated}
}




