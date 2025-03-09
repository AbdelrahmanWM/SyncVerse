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
	"github.com/AbdelrahmanWM/SyncVerse/document/global"

	"github.com/AbdelrahmanWM/signalingserver/signalingserver/message"
	"github.com/pion/webrtc/v4"
)

type WebRTCPeer struct {
	replicaID            global.ReplicaID
	replicaIDToServerIDs map[global.ReplicaID]string
	signalingServerConn  *signalingserverconn.SignalingServerConn
	peerConnections      map[string]*webrtcpeerconn.PeerConnection
	peerEvents           PeerEvents
	peerMode             PeerMode
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
	peerMode := DISCONNECTED
	p := &WebRTCPeer{replicaID, replicaIDToServerIDs, signalingServerConn, peerConnections, peerEvents, peerMode}
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
	conn, err := webrtcpeerconn.NewPeerConnection(&config, p.signalingServerConn, peerConnectionID, p.peerEvents)
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
func (peer *WebRTCPeer) SendDisconnectionMessage(notifyAll bool) error {
	if peer.signalingServerConn.Socket().IsUndefined() {
		Log("Socket connection not found.")
		return errors.New("[ERROR] socket connection not found") // temp
	}
	disconnectContent := message.DisconnectContent{notifyAll}
	disconnectContentJSON, err := json.Marshal(disconnectContent)
	if err != nil {
		Log("Failed to send disconnection message")
		return nil
	}
	disconnectionMsg := message.Message{
		Kind:    message.Disconnect,
		PeerID:  "",
		Content: disconnectContentJSON,
		Reach:   message.AllPeers,
		Sender:  "",
	}
	msgJSON, err := json.Marshal(disconnectionMsg)
	if err != nil {
		Log("Error marshalling message:" + err.Error())
		return nil
	}
	Log("Sending message: " + string(msgJSON))
	peer.signalingServerConn.Send(string(msgJSON))
	// disconnecting the signaling server client
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

// ////////////////////////////////////////
func (pc *WebRTCPeer) PushEvent(state PeerStatus, metaData PEMetadata) {
	pc.peerEvents.PushEvent(state, metaData)
}
func (p *WebRTCPeer) peerEventListener() {
	for event := range p.peerEvents {
		Log(fmt.Sprintf("Peer moved to state %s:", event.State))
		p.handlePeerEvent(event)
	}
}
func (p *WebRTCPeer) handlePeerEvent(event PeerEvent) {

	switch p.peerMode {
	case CONNECTED: // may need mux protection in the future
		switch event.State {

		case START_DISCONNECTING:
			p.peerMode = DISCONNECTING
			for _, pc := range p.peerConnections {
				pc.PushEvent(DISCONNECTION_INITIATED, nil)
			}
			p.SendDisconnectionMessage(true)
			p.DisconnectFromSignalingServer()
			p.peerMode = DISCONNECTED
			// disconnection logic
		case GOT_OFFER:
			// p.peerMode = CONNECTING ///////////////////
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
		case GOT_ALL_PEER_IDS:
			Log("Successfully got all peer ids")
			peerIDs, ok := event.Data.(GetAllPeersMetadata)
			if !ok {
				Log("Error parsing event data")
				break
			}
			if len(peerIDs.PeerIDS) == 0 {

				p.peerMode = CONNECTED
			} else {

				for _, peerID := range peerIDs.PeerIDS {
					err := p.NewPeerConnection(peerID)
					if err != nil {
						Log("Error creating peerConnection with peer: " + peerID)
					}
					p.peerConnections[peerID].PushEvent(PEER_CONNECTION_OPENED, nil)
					// p.SendOffer(peerID)
				}
			}

		}
	case DISCONNECTED: // may need mux protection in the future
		switch event.State {
		case START_CONNECTING:
			p.peerMode = CONNECTING
			p.ConnectToSignalingServer()
			p.signalingServerConn.SindIdentifySelfMessage()

		case PEER_CONNECTION_CLOSED:
			disconnectionMetadata, ok := event.Data.(PeerConnectionDisconnectedMetadata)
			if !ok {
				Log("Error extracting disconnected peer id")
				break
			}
			Log("removing peerconnection: " + disconnectionMetadata.PeerID)
			p.RemovePeerConnection(disconnectionMetadata.PeerID)
		}
	case CONNECTING:
		switch event.State {
		case SIGNALING_INITIATED:
			Log(fmt.Sprintf("Peer have initialized signaling"))
			p.peerMode = CONNECTED
			p.GetAllPeerIDs()

		case PEER_CONNECTION_AVAILABLE:
			if p.AllPeerConnectionsStable() {
				p.peerMode = CONNECTED
			}
		}
	case DISCONNECTING:
		switch event.State {
		case PEER_CONNECTION_CLOSED:
			disconnectionMetadata, ok := event.Data.(PeerConnectionDisconnectedMetadata)
			if !ok {
				Log("Error extracting disconnected peer id")
				break
			}
			Log("removing peerconnection: " + disconnectionMetadata.PeerID)
			p.RemovePeerConnection(disconnectionMetadata.PeerID)
		}
	}
}

func (p *WebRTCPeer) AllPeerConnectionsStable() bool {
	allPeerConnectionsStable := true
	for _, pc := range p.peerConnections {
		if !pc.Stable() {
			allPeerConnectionsStable = false
		}
	}
	return allPeerConnectionsStable
}

// ////////////////////////////////////////////

func (p *WebRTCPeer) GetReplicaID() global.ReplicaID {
	return p.replicaID
}

func (pr *WebRTCPeer) JoinSession(v js.Value, p []js.Value) any {
	pr.PushEvent(START_CONNECTING, nil)
	return nil
}
func (pr *WebRTCPeer) LeaveSession(v js.Value, p []js.Value) any {
	pr.PushEvent(START_DISCONNECTING, nil)
	return nil
}
func (pr *WebRTCPeer) GetPeerModeJS(v js.Value, p []js.Value) any {
	Log(string(pr.peerMode))
	return nil

}
func (pr *WebRTCPeer) GetAllPeersJS(v js.Value, p []js.Value) any {
	pr.GetAllPeerIDs()
	return nil
}
