package webrtc_peer

import (
	"encoding/json"
	"errors"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/crdt_conn/internal/signalingserverconn"
	. "github.com/AbdelrahmanWM/SyncVerse/document/crdt/crdt_conn/internal/utils"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/crdt_conn/internal/webrtcpeerconn"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/global"

	"github.com/AbdelrahmanWM/signalingserver/signalingserver/message"
	"github.com/pion/webrtc/v4"
)

type WebRTCPeer struct {
	replicaID            global.ReplicaID
	replicaIDToServerIDs map[global.ReplicaID]string
	signalingServerConn  *signalingserverconn.SignalingServerConn
	peerConnections      map[string]*webrtcpeerconn.PeerConnection
}

func NewWebRTCPeer(replicaID global.ReplicaID) *WebRTCPeer {
	peerConnections := make(map[string]*webrtcpeerconn.PeerConnection)
	replicaIDToServerIDs := make(map[global.ReplicaID]string)
	connectionMap := make(map[string]signalingserverconn.Connection)
	for key, pc := range peerConnections {
		connectionMap[key] = pc
	}
	signalingServerConn := signalingserverconn.NewSignalingServerConn(connectionMap)

	return &WebRTCPeer{replicaID, replicaIDToServerIDs, signalingServerConn, peerConnections}
}
func (p *WebRTCPeer) ConnectToSignalingServer() {
	p.signalingServerConn.Connect()
}
func (p *WebRTCPeer) DisconnectFromSignalingServer() {
	p.signalingServerConn.Disconnect()
}

func (p *WebRTCPeer) NewPeerConnection(peerConnectionID string) error {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}
	conn, err := webrtcpeerconn.NewPeerConnection(&config, p.signalingServerConn, peerConnectionID)
	if err != nil {
		return err
	}
	p.AddNewPeerConnection(peerConnectionID, conn)
	return nil
}
func (peer *WebRTCPeer) SendOffer(peerID string) error {
	return peer.peerConnections[peerID].SendOffer()
}
func (peer *WebRTCPeer) SendAnswer(peerID string) error {
	return peer.peerConnections[peerID].SendAnswer()
}

func (peer *WebRTCPeer) GetAllPeerIDs() error {
	if peer.signalingServerConn.Socket() == nil {
		Log("Socket connection not found.")
		return errors.New("[ERROR] socket connection not found") // temp
	}
	getAllPeerIDsMsg := message.Message{
		Kind:    message.GetAllPeerIDs,
		PeerID:  "",
		Content: nil,
		Reach:   message.Self,
		Sender:  "",
	}
	msgJSON, err := json.Marshal(getAllPeerIDsMsg)
	if err != nil {
		Log("Error marshalling message:" + err.Error())
		return err
	}
	Log("Sending message: " + string(msgJSON))
	peer.signalingServerConn.Send(string(msgJSON))
	return nil
}

func (p *WebRTCPeer) SendToAll(message string) error {////
	for _, pc := range p.peerConnections {
		err := pc.SendMessage([]byte(message))
		if err != nil {
			return err
		}
	}
	return nil
}
func (p *WebRTCPeer) SindIdentifySelfMessageJS() any {
	p.signalingServerConn.SindIdentifySelfMessage()
	return nil
}
func (p *WebRTCPeer) AddNewPeerConnection(id string, pc *webrtcpeerconn.PeerConnection) {
	p.peerConnections[id] = pc
	p.signalingServerConn.AddNewPeerConnection(id, pc)
}
func (p *WebRTCPeer) RemovePeerConnection(id string) {
	delete(p.peerConnections, id)
	p.signalingServerConn.RemovePeerConnection(id)
}

// ////////////////////////////////////////////
func (p *WebRTCPeer) JoinSession() {
	// connect to signaling server
	p.signalingServerConn.Connect()
	p.signalingServerConn.SindIdentifySelfMessage()

}
