package signalingserverconn

import (
	"encoding/json"
	"fmt"
	. "github.com/AbdelrahmanWM/SyncVerse/document/crdt/crdt_conn/internal/utils"
    "github.com/gorilla/websocket"
	"github.com/AbdelrahmanWM/signalingserver/signalingserver/message"
)

const signalingServerURL = "ws://localhost:8090/signalingserver"

type Connection interface {
	SetLocalDescription(input json.RawMessage) error
	SetRemoteDescription(input json.RawMessage) error
	AddICECandidate(input json.RawMessage) error
	SendPendingICECandidates() error
}
type SignalingServerConn struct {
	socket    *websocket.Conn
	peerID    string
	peerConns map[string]Connection
}

func NewSignalingServerConn(peerConns map[string]Connection) *SignalingServerConn {
	return &SignalingServerConn{peerConns: peerConns}
}
func (conn *SignalingServerConn) SindIdentifySelfMessage()error {
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
	err=conn.Send(string(identifySelfMsgJson))
	return err
}
func (conn *SignalingServerConn) Connect()error{
	dialer:=websocket.DefaultDialer
	socket,_,err := dialer.Dial(signalingServerURL,nil)
	if err!=nil{
		return err
	}
	go conn.handleSocketOnMessage()
	conn.socket = socket
	return nil
}

func(conn *SignalingServerConn) Disconnect()error{
	err:=conn.socket.Close()
	return err
}

func (conn *SignalingServerConn) Socket() *websocket.Conn {
	return conn.socket
}
func (conn *SignalingServerConn) Send(message string) error {
	// if conn.socket.Get("readyState").Int() != 1 {
	// 	return fmt.Errorf("Websocket is not open")
	// }
	// conn.socket.Call("send", message)
	err:=conn.socket.WriteMessage(0,[]byte(message))
	if err!=nil {
		Log(fmt.Sprintf("Error sending message: %v",err))
	}
	return err
}
func (conn *SignalingServerConn) handleSocketOnMessage() {
	for {

	_, messageData, err := conn.socket.ReadMessage()
	if err!=nil{
		Log(fmt.Sprintf("Error reading message: %v",err))
		break
	}
	var msg message.Message
	err = json.Unmarshal(messageData, &msg)
	if err != nil {
		Log("Error on unmarshaling message: " + err.Error())
	}
	Log("(" + msg.Sender + ") " + string(msg.Content))
	switch msg.Kind {
	case message.IdentifySelf:
		var identifyMsgContent message.IdentifySelfContent
		err := json.Unmarshal(msg.Content, &identifyMsgContent)
		if err != nil {
			Log("Error unmarshaling message content " + err.Error())
		}
		conn.peerID = identifyMsgContent.ID

	case message.Offer:
		targetPeer := msg.Sender
		targetPeerConn, ok := conn.peerConns[targetPeer]
		if !ok {
			Log("Error setting remote description")
			Log(fmt.Sprintf("%s->%#v", targetPeer, conn.peerConns))
			break
		}
		err := targetPeerConn.SetRemoteDescription(msg.Content)
		if err != nil {
			Log(fmt.Sprintf("Error setting remote description: %v", err))
		}
		err = targetPeerConn.SendPendingICECandidates()
		if err != nil {
			Log("Error sending ICE candidates: " + err.Error())
		}
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
		err = targetPeerConn.SendPendingICECandidates()
		if err != nil {
			Log("Error sending ICE candidates: " + err.Error())
		}
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
	}

	}
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
