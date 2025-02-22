//go:build js && wasm

// +build: js,wasm
package webrtc_peer

// TODO: refactor the module
import (
	"encoding/json"
	"errors"
	"fmt"
	"syscall/js"

	. "github.com/AbdelrahmanWM/SyncVerse/document/crdt/crdt_conn/conn_types"
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
	peerEvents           PeerEvents
}

func NewWebRTCPeer(replicaID global.ReplicaID) *WebRTCPeer {
	peerConnections := make(map[string]*webrtcpeerconn.PeerConnection)
	replicaIDToServerIDs := make(map[global.ReplicaID]string)
	connectionMap := make(map[string]signalingserverconn.Connection)
	for key, pc := range peerConnections {
		connectionMap[key] = pc
	}
	peerEvents := make(PeerEvents, 1)

	signalingServerConn := signalingserverconn.NewSignalingServerConn(&peerEvents, connectionMap)

	p := &WebRTCPeer{replicaID, replicaIDToServerIDs, signalingServerConn, peerConnections, peerEvents}
	go p.peerEventListener() // starting the event listener
	return p
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
	Log("sending offer")
	return peer.peerConnections[peerID].SendOffer()
}
func (peer *WebRTCPeer) SendAnswer(peerID string) error {
	return peer.peerConnections[peerID].SendAnswer()
}

func (peer *WebRTCPeer) GetAllPeerIDs() error {
	if peer.signalingServerConn.Socket().IsUndefined() {
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

func (p *WebRTCPeer) SendToAll(message string) error { ////
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

//////////////////////////////////////////

func (p *WebRTCPeer) peerEventListener() {
	for event := range p.peerEvents {
		Log(fmt.Sprintf("Peer moved to state %s:", event.State))
		p.handlePeerEvent(event)
	}
}
func (p *WebRTCPeer) handlePeerEvent(event PeerEvent) {
	switch event.State {
	case SignalingInitiated:
		Log(fmt.Sprintf("Peer have initialized signaling"))
		p.GetAllPeerIDs()
	case GotAllPeerIDs:
		Log("Successfully got all peer ids")
		peerIDs, ok := event.Data.(GetAllPeersMetadata)
		if !ok {
			Log("Error parsing event data")
			break
		}
		for _, peerID := range peerIDs.PeerIDS {
			err := p.NewPeerConnection(peerID)
			if err != nil {
				Log("Error creating peerConnection with peer: " + peerID)
			}
			p.peerConnections[peerID].PushEvent(PeerEvent{PeerConnectionOpened, nil})
			// p.SendOffer(peerID)
		}
	case GotAnOffer:
		Peer, ok := event.Data.(GotAnOfferMetadata)
		if !ok {
			Log("Error setting remote description")
			break
		}
		pc, ok := p.peerConnections[Peer.SenderID]
		if !ok {
			p.NewPeerConnection(Peer.SenderID)
			pc = p.peerConnections[Peer.SenderID]
		}
		err := pc.SetRemoteDescription(Peer.Offer)
		if err != nil {
			Log(fmt.Sprintf("Error setting remote description: %v", err))
		}
		// pc.PushEvent(PeerEvent{OfferDescriptionSet,nil})
	// case SendPendingICECandidates:
	// 		err = pc.SendPendingICECandidates()
	// 		if err != nil {
	// 			Log("Error sending ICE candidates: " + err.Error())
	// 		}
	// 		pc.SendAnswer()
	}
}

// ////////////////////////////////////////////
// func (p *WebRTCPeer) JoinSession() {
// 	// connect to signaling server
// 	p.signalingServerConn.Connect()
// 	p.signalingServerConn.SindIdentifySelfMessage()

// }
func (p *WebRTCPeer) GetReplicaID() global.ReplicaID {
	return p.replicaID
}

func (p *WebRTCPeer) ConnectToMesh() {
	p.ConnectToSignalingServer()
	p.GetAllPeerIDs()
}

// func (pr *WebRTCPeer) GetAllPeers(v js.Value, p []js.Value) any {
// 	pr.GetAllPeerIDs()
// 	return nil
// }

func (pr *WebRTCPeer) JoinSession(v js.Value, p []js.Value) any {
	pr.ConnectToSignalingServer()
	pr.signalingServerConn.SindIdentifySelfMessage()
	return nil
}
