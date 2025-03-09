package types

import(
	"github.com/AbdelrahmanWM/SyncVerse/document/crdt/vector_clock"
	"strings"
	"strconv"
)
type ModifyMetadata struct {
	ClockOffset *vector_clock.ClockOffset
	Rng         [2]int
}

func (mm *ModifyMetadata) String() string {
	var result strings.Builder
	if mm.ClockOffset != nil {
		result.WriteString(mm.ClockOffset.String())
	} else {
		result.WriteString("nil ClockOffset")
	}

	result.WriteString("\n")
	result.WriteString("Range: ")
	result.WriteString(strconv.Itoa(mm.Rng[0]))
	result.WriteString(" ")
	result.WriteString(strconv.Itoa(mm.Rng[1]))
	return result.String()
}

type ModifyMetadataArray []*ModifyMetadata

func (mma ModifyMetadataArray) String() string {
	var result strings.Builder
	for _, mm := range mma {
		result.WriteString(mm.String())
		result.WriteString("\n---\n")
	}
	return result.String()
}