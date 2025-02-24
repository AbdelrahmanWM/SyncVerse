//go:build js && wasm

// +build: js,wasm
package webrtcpeerconn

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	. "github.com/AbdelrahmanWM/SyncVerse/document/crdt/crdt_conn/conn_types"
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/crdt_conn/internal/signalingserverconn"
	. "github.com/AbdelrahmanWM/SyncVerse/document/crdt/crdt_conn/internal/utils"
	"github.com/AbdelrahmanWM/signalingserver/signalingserver/message"
	"github.com/pion/webrtc/v4"
)

type PeerConnection struct {
	peerIDs              [2]string
	signalingServerConn  *signalingserverconn.SignalingServerConn
	peerConnection       *webrtc.PeerConnection
	dataChannel          *webrtc.DataChannel
	pendingCandidates    []*webrtc.ICECandidate
	candidateMux         sync.Mutex
	peerConnectionEvents PeerConnectionEvents
	peerConnectionMode PeerConnectionMode
	notifyPeer PeerEvents
}

func NewPeerConnection(config *webrtc.Configuration, signalingServerConn *signalingserverconn.SignalingServerConn, connectedPeerID string,notifyPeer PeerEvents) (*PeerConnection, error) {
	peerIDs := [2]string{signalingServerConn.PeerID(), connectedPeerID}
	pendingCandidates := make([]*webrtc.ICECandidate, 0)
	var candidateMux sync.Mutex
	pc, err := webrtc.NewPeerConnection(*config)
	if err != nil {
		return nil, err
	}
	dataChannel, err := pc.CreateDataChannel("chan"+peerIDs[0]+peerIDs[1], nil)
	if err != nil {
		Log("Problem creating dataChannel")
		return nil, err
	}

	peerConnectionEvents := make(PeerConnectionEvents, 1)
	PeerConnectionMode:=BUSY
	peerConnection := PeerConnection{peerIDs, signalingServerConn, pc, dataChannel, pendingCandidates, candidateMux, peerConnectionEvents, PeerConnectionMode,notifyPeer}

	// handlers for datachannel events
	dataChannel.OnOpen(func() { peerConnection.dataChannelOnOpen(dataChannel) })
	dataChannel.OnClose(func() { peerConnection.dataChannelOnClose(dataChannel) })
	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) { peerConnection.dataChannelOnMessage(msg, dataChannel) })

	// handlers for peer connection events
	pc.OnConnectionStateChange(func(s webrtc.PeerConnectionState) { handleConnectionStateChange(s) })
	pc.OnDataChannel(func(d *webrtc.DataChannel) { handleDataChannel(d) })
	pc.OnICECandidate(func(candidate *webrtc.ICECandidate) { peerConnection.OnICECandidate(candidate) })

	// handlers for peer connection object events
	go peerConnection.peerConnectionEventListener()
	return &peerConnection, nil
}

func (p *PeerConnection) SendOffer() error {
	offer, err := p.peerConnection.CreateOffer(nil)
	if err != nil {
		return err
	}
	if err := p.peerConnection.SetLocalDescription(offer); err != nil {
		return err
	}
	offerObj, err := json.Marshal(message.OfferContent{Type: int(offer.Type), SDP: offer.SDP})
	if err != nil {
		Log("Error marshalling offer message")
		return err
	}
	msg := message.Message{
		Kind:    message.Offer,
		Sender:  p.peerIDs[0],
		PeerID:  p.peerIDs[1], // the opposite peer
		Reach:   message.OnePeer,
		Content: offerObj,
	}
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		Log("Error marshalling message: " + err.Error())
	}
	err = p.signalingServerConn.Send(string(msgJSON))
	if err != nil {
		Log("Error sending message: " + err.Error())
	}
	return nil
}
func (p *PeerConnection) SetLocalDescription(input json.RawMessage) error {
	var sdp message.OfferContent
	err := json.Unmarshal(input, &sdp)
	if err != nil {
		Log("Error marshalling message: " + err.Error())
	}
	var sdpType webrtc.SDPType
	switch sdp.Type {
	case 1:
		sdpType = webrtc.SDPTypeOffer
	case 3:
		sdpType = webrtc.SDPTypeAnswer
	default:
		return fmt.Errorf("sdp type in neither an offer nor an answer")
	}
	description := webrtc.SessionDescription{Type: sdpType, SDP: sdp.SDP}

	go func() {
		err := p.peerConnection.SetRemoteDescription(description)
		if err != nil {
			Log(err.Error())
		} else {
			Log("Successfully set remote description using sdp")
		}

	}()
	return nil //temp
}

func (p *PeerConnection) SetRemoteDescription(input json.RawMessage) error {

	var sdp message.OfferContent
	err := json.Unmarshal(input, &sdp)
	if err != nil {
		Log("Error marshalling message: " + err.Error())
	}
	var sdpType webrtc.SDPType
	switch sdp.Type {
	case 1:
		sdpType = webrtc.SDPTypeOffer
	case 3:
		sdpType = webrtc.SDPTypeAnswer
	default:
		return fmt.Errorf("sdp type in neither an offer nor an answer")
	}
	description := webrtc.SessionDescription{Type: sdpType, SDP: sdp.SDP}

	go func() {
		err := p.peerConnection.SetRemoteDescription(description)
		if err != nil {
			Log(err.Error())
		} else {
			Log("Successfully set remote description using sdp")
		}
		//temp
		p.PushEvent(OFFER_DESCRIPTION_SET, nil)

		// p.signalingServerConn.

	}()
	return nil //temp
}
func (p *PeerConnection) AddICECandidate(input json.RawMessage) error {
	var iceCandidateContent message.ICECandidateContent
	err := json.Unmarshal(input, &iceCandidateContent)
	if err != nil {
		Log("Error marshalling message: " + err.Error())
	}
	go func() {
		// Perform the addICECandidate operation inside the goroutine
		err := p.peerConnection.AddICECandidate(webrtc.ICECandidateInit{Candidate: iceCandidateContent.Candidate, SDPMid: iceCandidateContent.SdpMid, SDPMLineIndex: iceCandidateContent.SdpMLineIndex, UsernameFragment: iceCandidateContent.UsernameFragment})
		if err != nil {
			Log("Error adding ICE candidate: " + err.Error())
		} else {
			Log("Successfully added ICE Candidate")
		}

	}()
	return nil
}
func (p *PeerConnection) SendPendingICECandidates() error {
	var err error
	for _, candidate := range p.pendingCandidates {
		if err = p.peerConnection.AddICECandidate(candidate.ToJSON()); err != nil {
			Log("Error sending ICE candidate: " + err.Error())

		}
	}
	p.pendingCandidates = nil
	return err
}
func (p *PeerConnection) SendAnswer() error {
	answer, err := p.peerConnection.CreateAnswer(nil)
	if err != nil {
		return err
	}
	if err := p.peerConnection.SetLocalDescription(answer); err != nil {
		return err
	}
	offerObj, err := json.Marshal(message.OfferContent{Type: int(answer.Type), SDP: answer.SDP})
	if err != nil {
		Log("Error marshalling offer message")
		return err
	}
	msg := message.Message{
		Kind:    message.Offer,
		Sender:  p.peerIDs[0], // self
		PeerID:  p.peerIDs[1], // the opposite peer
		Reach:   message.OnePeer,
		Content: offerObj,
	}
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		Log("Error marshalling message: " + err.Error())
	}
	err = p.signalingServerConn.Send(string(msgJSON))
	if err != nil {
		Log("Error sending message: " + err.Error())
	}
	return nil
}
func (p *PeerConnection) SendMessage(message []byte) error {

	err := p.dataChannel.Send(message)
	if err != nil {
		Log("Error sending message on dataChannel " + p.dataChannel.Label() + err.Error())
		return err
	}
	return nil
}

func handleConnectionStateChange(s webrtc.PeerConnectionState) {
	log.Printf("Peer connection state has changed: %s\n", s.String())
	if s == webrtc.PeerConnectionStateFailed {
		Log("Peer connection has gone to failed exiting")
		os.Exit(0)
	}
	if s == webrtc.PeerConnectionStateClosed {
		Log("Peer connection has gone to closed exiting")
		os.Exit(0)
	}
}
func handleDataChannel(d *webrtc.DataChannel) {
	{
		Log(fmt.Sprintf("New data channel '%s'-'%d' open.", d.Label(), *d.ID()))
		d.OnOpen(func() {
			Log(fmt.Sprintf("Data channel '%s'-'%d' open.", d.Label(), *d.ID()))
		})
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			Log(fmt.Sprintf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data)))
		})
		d.OnClose(func() {
			Log(fmt.Sprintf("Data channel '%s'-'%d' closed.", d.Label(), *d.ID()))
		})
	}
}
func (pc *PeerConnection) dataChannelOnOpen(dataChannel *webrtc.DataChannel) {
	Log(fmt.Sprintf("Data channel '%s'-'%d' open.", dataChannel.Label(), *dataChannel.ID()))
	pc.PushEvent(PEER_CONNECTION_ESTABLISHED, nil)
}
func (pc *PeerConnection) dataChannelOnClose(dataChannel *webrtc.DataChannel) {
	Log(fmt.Sprintf("Data channel '%s'-'%d' closed.", dataChannel.Label(), *dataChannel.ID()))
}
func (pc *PeerConnection) dataChannelOnMessage(msg webrtc.DataChannelMessage, dataChannel *webrtc.DataChannel) {
	Log(fmt.Sprintf("Message from DataChannel '%s': '%s'\n", dataChannel.Label(), string(msg.Data)))
}
func (pc *PeerConnection) OnICECandidate(c *webrtc.ICECandidate) {
	if c == nil {
		return
	}
	pc.candidateMux.Lock()
	defer pc.candidateMux.Unlock()

	desc := pc.peerConnection.RemoteDescription()
	if desc == nil {
		pc.pendingCandidates = append(pc.pendingCandidates, c)
	} else {
		candidateContent := message.ICECandidateContent{Candidate: c.ToJSON().Candidate, SdpMid: c.ToJSON().SDPMid, SdpMLineIndex: c.ToJSON().SDPMLineIndex, UsernameFragment: c.ToJSON().UsernameFragment}
		candidateContentJson, err := json.Marshal(candidateContent)
		if err != nil {
			Log("Error marshalling message content :" + err.Error())
		}

		candidateMsg := message.Message{
			Kind:    message.ICECandidate,
			PeerID:  pc.peerIDs[1],
			Reach:   message.OnePeer,
			Sender:  pc.peerIDs[0],
			Content: candidateContentJson,
		}
		candidateMsgJson, err := json.Marshal(candidateMsg)
		if err != nil {
			Log("Error marshaling message: " + err.Error())
		}
		err = pc.signalingServerConn.Send(string(candidateMsgJson))
		if err != nil {
			Log("Error sending message: " + err.Error())
		}
	}
}

func (pc *PeerConnection) peerConnectionEventListener() {
	for event := range pc.peerConnectionEvents {
		Log(fmt.Sprintf("Peer connection moved to state %s:", event.State))
		pc.handlePeerConnectionEvent(&event)
	}
}
func (pc *PeerConnection) handlePeerConnectionEvent(event *PeerConnectionEvent) {
	switch event.State {
	case PEER_CONNECTION_OPENED:
		pc.SendOffer()
	case OFFER_DESCRIPTION_SET:
		err := pc.SendPendingICECandidates()
		if err != nil {
			Log("Error sending ICE candidates: " + err.Error())
		}
		Log("Successfully sent ICE candidates")
		pc.SendAnswer()
	case PEER_CONNECTION_ESTABLISHED:
		pc.peerConnectionMode=AVAILABLE
		pc.notifyPeer<-PeerEvent{}
		pc.SendMessage([]byte("Hello Peer!"))
	}

}
func (pc *PeerConnection) PushEvent(state PeerConnectionStatus, metadata PEMetadata) {
	pc.peerConnectionEvents <- PeerConnectionEvent{state, metadata}
}
