//go:build js && wasm

// +build: js,wasm
package main

import (
	// "syscall/js"

	. "github.com/AbdelrahmanWM/SyncVerse/document/crdt/crdt_conn/internal/utils"
	// "github.com/AbdelrahmanWM/SyncVerse/document/crdt/crdt_conn/webrtc_peer"
)

func main() {

	// peer := webrtc_peer.NewWebRTCPeer("A")
	Log("New peer!")
	// js.Global().Set("connectToSignalingServer", js.FuncOf(peer.ConnectToSignalingServer))
	// js.Global().Set("disconnectFromSignalingServer", js.FuncOf(peer.DisconnectFromSignalingServer))
	// js.Global().Set("getAllPeerIDs", js.FuncOf(peer.GetAllPeerIDs))
	// js.Global().Set("newPeerConnection", js.FuncOf(peer.NewPeerConnectionJS))
	// js.Global().Set("clearLog", js.FuncOf(ClearLog))
	// js.Global().Set("sendToAll", js.FuncOf(peer.SendToAll))
	// js.Global().Set("sindIdentifySelfMessage", js.FuncOf(peer.SindIdentifySelfMessageJS))
	select {}
}
