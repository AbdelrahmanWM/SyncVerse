//go:build js && wasm

// +build: js,wasm
package signalingserverconn

import (
	"encoding/json"
	"fmt"

	"syscall/js"

	. "github.com/AbdelrahmanWM/SyncVerse/document/crdt/crdt_conn/conn_types"
	. "github.com/AbdelrahmanWM/SyncVerse/document/global/utils"
	"github.com/AbdelrahmanWM/signalingserver/signalingserver/message"
)

const signalingServerURL = "ws://127.0.0.1:8090/signalingserver"

type Connection interface {
	SetLocalDescription(input json.RawMessage) error
	SetRemoteDescription(input json.RawMessage) error
	AddICECandidate(input json.RawMessage) error
	SendPendingICECandidates() error
	PushEvent(status PeerConnectionStatus, content PEMetadata)
}
type SignalingServerConn struct {
	socket               js.Value
	peerID               string
	peerNotificationChan *PeerEvents
	peerConns            map[string]Connection
}

func NewSignalingServerConn(peerNotificationChan *PeerEvents, peerConns map[string]Connection) *SignalingServerConn {
	return &SignalingServerConn{peerNotificationChan: peerNotificationChan, peerConns: peerConns}
}
func (conn *SignalingServerConn) SindIdentifySelfMessage() error {
	identifySelfMsgContent := message.IdentifySelfContent{ID: ""}
	identifySelfMsgContentJson, err := json.Marshal(identifySelfMsgContent)
	if err != nil {
		return err
	}
	identifySelfMsg := message.Message{
		Kind:    message.IdentifySelf,
		Reach:   message.Self,
		PeerID:  "",
		Content: identifySelfMsgContentJson,
	}
	identifySelfMsgJson, err := json.Marshal(identifySelfMsg)
	if err != nil {
		return err
	}
	err = conn.Send(string(identifySelfMsgJson))
	return err
}
func (conn *SignalingServerConn) Connect() error {
	socket := js.Global().Get("WebSocket").New(signalingServerURL)
	if socket.IsUndefined() {
		Log("Failed to create WebSocket")
		return fmt.Errorf("failed to create WebSocket")
	}

	socket.Set("onopen", js.FuncOf(conn.handleSocketOnOpen))
	socket.Set("onmessage", js.FuncOf(conn.handleSocketOnMessage))
	socket.Set("onclose", js.FuncOf(conn.handleSocketOnClose))
	socket.Set("onError", js.FuncOf(conn.handleSocketOnError))

	conn.socket = socket
	Log("Connected to WebSocket server")
	return nil
}

func (conn *SignalingServerConn) Disconnect() error {
	if conn.socket.Get("readyState").Int() == 1 {
		conn.socket.Call("close")
		Log("WebSocket connection closed.")
	} else {
		Log("WebSocket is not open, cannot disconnect.")
	}
	return nil
}

func (conn *SignalingServerConn) Socket() js.Value {
	return conn.socket
}
func (conn *SignalingServerConn) Send(message string) error {
	if conn.socket.Get("readyState").Int() != 1 {
		return fmt.Errorf("Websocket is not open")
	}
	conn.socket.Call("send", message)
	return nil
}
func (conn *SignalingServerConn) handleSocketOnOpen(v js.Value, p []js.Value) any {
	Log("Websocket connected!")
	conn.peerNotificationChan.PushEvent(SIGNALING_INITIATED, nil)

	return nil
}
func (conn *SignalingServerConn) handleSocketOnMessage(v js.Value, p []js.Value) any {
	event := p[0]
	messageData := event.Get("data").String() // Get message from event
	var msg message.Message
	err := json.Unmarshal([]byte(messageData), &msg)
	if err != nil {
		Log("Error on unmarshaling message: " + err.Error())
		return nil
	}
	Log("(" + msg.Sender + ") " + string(msg.Content))
	switch msg.Kind {
	case message.IdentifySelf:
		var identifyMsgContent message.IdentifySelfContent
		err := json.Unmarshal(msg.Content, &identifyMsgContent)
		if err != nil {
			Log("Error unmarshaling message content " + err.Error())
			break
		}
		conn.peerID = identifyMsgContent.ID

	case message.Offer:
		targetPeer := msg.Sender
		conn.peerNotificationChan.PushEvent(GOT_OFFER, GotAnOfferMetadata{targetPeer, msg.Content})

	case message.Answer:
		targetPeer := msg.Sender
		targetPeerConn, ok := conn.peerConns[targetPeer]
		if !ok {
			Log("Error setting remote description")
			break
		}
		err := targetPeerConn.SetRemoteDescription(msg.Content)
		if err != nil {
			Log(fmt.Sprintf("Error setting remote description: %v", err))
		}
		// err = targetPeerConn.SendPendingICECandidates()
		// if err != nil {
		// 	Log("Error sending ICE candidates: " + err.Error())
		// }
	case message.ICECandidate:
		targetPeer := msg.Sender
		targetPeerConn, ok := conn.peerConns[targetPeer]
		if !ok {
			Log("Error setting remote description")
			break
		}
		if candidateErr := targetPeerConn.AddICECandidate(msg.Content); candidateErr != nil {
			Log("Error adding ICE candidate: " + candidateErr.Error())
		}
	case message.GetAllPeerIDs:
		var peers message.GetAllPeerIDsContent
		err = json.Unmarshal(msg.Content, &peers)
		if err != nil {
			Log("Error unmarshaling GetAllPeerIDsContent")
			break
		}
		var connectToPeers = false
		for _, peerID := range peers.PeersIDs {
			if _, ok := conn.peerConns[peerID]; !ok {
				connectToPeers = true
				break
			}
		}
		if len(peers.PeersIDs) == 0 && len(conn.peerConns) == 0 { // first peer connecting
			connectToPeers = true
		}
		if connectToPeers {
			conn.peerNotificationChan.PushEvent(GOT_ALL_PEER_IDS, GetAllPeersMetadata{peers.PeersIDs})
		}
	case message.DisconnectionNotification:

		var disconnectionNotificationContent message.DisconnectionNotificationContent
		err := json.Unmarshal(msg.Content, &disconnectionNotificationContent)
		if err != nil {
			Log("Error unmarshaling disconnectionNotificationContent")
			break
		}
		targetPC := conn.peerConns[disconnectionNotificationContent.DisconnectedPeerID]
		if targetPC != nil {
			targetPC.PushEvent(DISCONNECTION_INITIATED, nil)
			// conn.RemovePeerConnection(disconnectionNotificationContent.DisconnectedPeerID)
		}
	}

	return nil
}
func (conn *SignalingServerConn) handleSocketOnClose(v js.Value, p []js.Value) any {
	Log("Connection closed.")
	return nil
}
func (conn *SignalingServerConn) handleSocketOnError(v js.Value, p []js.Value) any {
	Log("Error with websocket connection")
	return nil
}

func (conn *SignalingServerConn) AddNewPeerConnection(id string, peerConn Connection) {
	conn.peerConns[id] = peerConn
}
func (conn *SignalingServerConn) RemovePeerConnection(id string) {
	delete(conn.peerConns, id)
}
func (conn *SignalingServerConn) PeerID() string {
	return conn.peerID
}
