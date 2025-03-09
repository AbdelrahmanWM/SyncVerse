package vector_clock

import (
	"slices"
	"strconv"
	"strings"

	"github.com/AbdelrahmanWM/SyncVerse/document/global"
)

type VectorClock map[global.ReplicaID]int

func NewVectorClock(replicaID global.ReplicaID) VectorClock {
	if replicaID == "" {
		return VectorClock{}
	}
	return VectorClock{replicaID: 1}
}
func (v VectorClock) NewVectorClock(replicaID global.ReplicaID) VectorClock { // for now, constant size increment as number of concurrent users is minimal (<10)
	newVectorClock := copyMap(v)
	_, ok := v[replicaID]
	if ok {
		newVectorClock[replicaID]++
	} else {
		newVectorClock[replicaID] = 1
	}
	return newVectorClock
}
func (v VectorClock) Copy() VectorClock {
	return copyMap(v)

}
func (vc VectorClock) CompareHashes(vc2 VectorClock) int {
	str1 := vc.String()
	str2 := vc2.String()
	if str1 == str2 {
		return 0
	} else if str1 < str2 {
		return -1
	} else {
		return 1
	}
}
func (vc1 VectorClock) Equals(vc2 VectorClock) bool {
	for k, v := range vc1 {
		if v != vc2[k] {
			return false
		}
	}
	return len(vc1) == len(vc2)
}

// return -1 (earlier), 0 (concurrent), 1 (later)
func (vc1 VectorClock) Compare(vc2 VectorClock) int {
	vc1Higher := true
	vc2Higher := true

	union := mergeMaps(vc1, vc2)
	for k, _ := range union {
		v1, ok1 := vc1[k]
		v2, ok2 := vc2[k]
		if ok1 && ok2 {
			if v1 > v2 {
				vc2Higher = false
			} else if v1 < v2 {
				vc1Higher = false
			}
		} else if ok1 {
			vc2Higher = false
		} else if ok2 {
			vc1Higher = false
		}
	}

	if vc1Higher && vc2Higher {
		return 0
	} else if vc1Higher {
		return 1
	} else if vc2Higher {
		return -1
	} else {
		return 0
	}
}
func (vc1 VectorClock) IsValidSuccessor(vc2 VectorClock) bool {
	union := mergeMaps(vc1, vc2)
	for k, _ := range union {
		v1, ok1 := vc1[k]
		v2, ok2 := vc2[k]
		if ok1 && ok2 {
			if v2-v1 > 1 {
				return false
			}
		} else if ok2 {
			if v2 > 1 {
				return false
			}
		}
	}
	return true
}

func (vc1 VectorClock) Merge(vc2 VectorClock) VectorClock {
	return mergeMaps(vc1, vc2)
}

func mergeMaps(m1 map[global.ReplicaID]int, m2 map[global.ReplicaID]int) map[global.ReplicaID]int {
	union := make(map[global.ReplicaID]int)
	for k, val := range m1 {
		union[k] = val
	}
	for k, val := range m2 {
		union[k] = max(val, union[k])
	}
	return union
}

func copyMap(m map[global.ReplicaID]int) map[global.ReplicaID]int {
	newM := make(map[global.ReplicaID]int)
	for k, v := range m {
		newM[k] = v
	}
	return newM
}
func toKeysArray(m map[global.ReplicaID]int) []global.ReplicaID {
	results := make([]global.ReplicaID, len(m))
	i := 0
	for k, _ := range m {
		results[i] = k
		i++
	}
	return results
}
func (vc VectorClock) String() string {
	var result strings.Builder
	keys := toKeysArray(vc)
	slices.Sort(keys)
	for i, v := range keys {
		result.WriteString(string(v))
		result.WriteString("-")
		result.WriteString(strconv.Itoa(vc[v]))
		if i+1 != len(keys) {
			result.WriteString("-")
		}
	}
	return result.String()

}
