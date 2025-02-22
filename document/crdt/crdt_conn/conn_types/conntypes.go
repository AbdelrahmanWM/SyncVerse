package conntypes

import "encoding/json"

type PeerStatus string

const (
	SignalingInitiated PeerStatus = "SIGNALING_INITIATED"//peer
	GotAllPeerIDs      PeerStatus = "GOT_ALL_PEER_IDS"//peer
	PeerConnectionOpened PeerStatus = "PEER_CONNECTION_OPENED" //pc
	GotAnOffer         PeerStatus = "GOT_AN_OFFER" //peer
	OfferDescriptionSet PeerStatus = "OFFER_DESCRIPTION_SET"//pc
	// SendPendingICECandidates PeerStatus = "SEND_PENDING_ICE_CANDIDATES"
	PeerConnectionEstablished PeerStatus = "PEER_CONNECTION_ESTABLISHED"
)

type PeerEvent struct {
	State PeerStatus
	Data  PEMetadata
}

type PeerEvents chan PeerEvent

type PeerConnectionEvents chan PeerEvent

type PEMetadata interface{}

type GetAllPeersMetadata struct {
	PeerIDS []string
}

type GotAnOfferMetadata struct {
	SenderID string
	Offer    json.RawMessage
}