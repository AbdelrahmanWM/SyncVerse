package conntypes

import "encoding/json"

type PeerStatus string

const (
	START_CONNECTING PeerStatus = "START_CONNECTING"//peer 1
	SIGNALING_INITIATED PeerStatus = "SIGNALING_INITIATED"//peer 2
	GOT_ALL_PEER_IDS      PeerStatus = "GOT_ALL_PEER_IDS"//peer 3
	GOT_OFFER         PeerStatus = "GOT_OFFER" //peer 5 
    PEER_CONNECTION_AVAILABLE PeerStatus = "PEER_CONNECTION_AVAILABLE"   //peer 8
	START_DISCONNECTING PeerStatus = "START_DISCONNECTING" // peer -1
)

type PeerEvent struct {
	State PeerStatus
	Data  PEMetadata
}

type PeerEvents chan PeerEvent

func (pes PeerEvents) PushEvent(state PeerStatus, Data PEMetadata){
	pes<-PeerEvent{state,Data}
}

////////////////////////////////////////////////////////////////////

type PeerConnectionStatus string 
const (
	PEER_CONNECTION_OPENED PeerConnectionStatus = "PEER_CONNECTION_OPENED" //pc 5
	OFFER_DESCRIPTION_SET PeerConnectionStatus = "OFFER_DESCRIPTION_SET"//pc 6
	PEER_CONNECTION_ESTABLISHED PeerConnectionStatus = "PEER_CONNECTION_ESTABLISHED" //pc 7
)
type PeerConnectionEvent struct {
	State PeerConnectionStatus
	Data PEMetadata
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


type PeerMode string 

const (
	CONNECTING PeerMode = "CONNECTING"
	DISCONNECTING PeerMode = "DISCONNECTING"
	STABLE PeerMode = "STABLE"
)

type PeerConnectionMode  string
const (
	AVAILABLE PeerConnectionMode = "AVAILABLE"
	BUSY PeerConnectionMode = "BUSY"
)