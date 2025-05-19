//go:build js && wasm

// +build: js,wasm
package main

import (
	"math/rand"
	"strconv"
	"syscall/js"

	. "github.com/AbdelrahmanWM/SyncVerse/document/global/utils"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/crdt_conn/webrtc_peer"
	"github.com/AbdelrahmanWM/SyncVerse/document/global"
)

func main() {
	Log("New peer!")
	replicaID := getRandomID() //temp
	peer := webrtc_peer.NewWebRTCPeer(replicaID)
	js.Global().Set("joinSession", js.FuncOf(peer.JoinSession))
	js.Global().Set("leaveSession", js.FuncOf(peer.LeaveSession))
	js.Global().Set("getAllPeers", js.FuncOf(peer.GetAllPeersJS))
	js.Global().Set("getPeerMode", js.FuncOf(peer.GetPeerModeJS))
	// js.Global().Set("getAllPeers",js.FuncOf(peer.GetAllPeers))
	select {}
}

func getRandomID() global.ReplicaID { // temp
	return global.ReplicaID(strconv.Itoa(rand.Int()))
}
