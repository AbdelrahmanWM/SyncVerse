package sequence_number

import (
	"strconv"

	"github.com/AbdelrahmanWM/SyncVerse/document/global"
)

// Event sequence number
type SeqNum int

func (sn *SeqNum)Increment(){
	*sn+=1
}
func (sn *SeqNum)String()string{
	return strconv.Itoa(int(*sn))
}

////////////////////////////////////
type SeqNumMap map[global.ReplicaID]SeqNum

func (snm SeqNumMap) Update(replicaID global.ReplicaID, seqNum SeqNum) {
	if v, exists := snm[replicaID]; exists {
		snm[replicaID] = max(v, seqNum)
	} else {
		snm[replicaID] = seqNum
	}
}
func NewSeqNumMap() SeqNumMap {
	return make(map[global.ReplicaID]SeqNum)
}
