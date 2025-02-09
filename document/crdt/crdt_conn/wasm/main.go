//go:build js && wasm

// +build: js,wasm
package main

import (
	"math/rand"
	"strconv"
	"syscall/js"

	crdtconn "github.com/AbdelrahmanWM/SyncVerse/document/crdt/crdt_conn"
	. "github.com/AbdelrahmanWM/SyncVerse/document/crdt/crdt_conn/internal/utils"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/global"
)

func main() {
	Log("New peer!")
	js.Global().Set("joinSession", js.FuncOf(JoinSession))
	select {}
}
func JoinSession(v js.Value, p []js.Value) any {
	clientID := getRandomID()
	peerConnectionManager := crdtconn.GetPeerConnectionManager()
	peerConnectionManager.AddNewPeer(clientID)
	return nil
}
func getRandomID() global.ReplicaID { // temp
	return global.ReplicaID(strconv.Itoa(rand.Int()))
}
