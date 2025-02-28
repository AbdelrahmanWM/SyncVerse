package conntypes

import "encoding/json"

type PeerStatus string

const (
	START_CONNECTING          PeerStatus = "START_CONNECTING"          //1
	SIGNALING_INITIATED       PeerStatus = "SIGNALING_INITIATED"       //2
	GOT_ALL_PEER_IDS          PeerStatus = "GOT_ALL_PEER_IDS"          //3
	GOT_OFFER                 PeerStatus = "GOT_OFFER"                 //5
	PEER_CONNECTION_AVAILABLE PeerStatus = "PEER_CONNECTION_AVAILABLE" //8
	/////////////////
	START_DISCONNECTING    PeerStatus           = "START_DISCONNECTING"    //-1
	PEER_CONNECTION_CLOSED PeerStatus = "PEER_CONNECTION_CLOSED" //-3
)

type PeerEvent struct {
	State PeerStatus
	Data  PEMetadata
}

type PeerEvents chan PeerEvent

func (pes PeerEvents) PushEvent(state PeerStatus, Data PEMetadata) {
	pes <- PeerEvent{state, Data}
}

////////////////////////////////////////////////////////////////////

type PeerConnectionStatus string

const (
	PEER_CONNECTION_OPENED      PeerConnectionStatus = "PEER_CONNECTION_OPENED"      //5
	OFFER_DESCRIPTION_SET       PeerConnectionStatus = "OFFER_DESCRIPTION_SET"       //6
	PEER_CONNECTION_ESTABLISHED PeerConnectionStatus = "PEER_CONNECTION_ESTABLISHED" //7
	////////////////////////////////
	DISCONNECTION_INITIATED PeerConnectionStatus = "DISCONNECTION_INITIATED" //-2
)

type PeerConnectionEvent struct {
	State PeerConnectionStatus
	Data  PEMetadata
}
type PeerConnectionEvents chan PeerConnectionEvent

//////////////////////////////////////////////////////
type PEMetadata interface{}

type GetAllPeersMetadata struct {
	PeerIDS []string
}

type GotAnOfferMetadata struct {
	SenderID string
	Offer    json.RawMessage
}
type PeerConnectionDisconnectedMetadata struct{
	PeerID string
}
type PeerMode string

const (
	CONNECTING    PeerMode = "CONNECTING"
	DISCONNECTING PeerMode = "DISCONNECTING"
	CONNECTED        PeerMode = "CONNECTED"
	DISCONNECTED PeerMode = "DISCONNECTED"
)

type PeerConnectionMode string

const (
	AVAILABLE PeerConnectionMode = "AVAILABLE"
	BUSY      PeerConnectionMode = "BUSY"
)
