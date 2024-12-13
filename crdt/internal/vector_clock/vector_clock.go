package vector_clock

import (
	"slices"
	"strconv"
)

type VectorClock map[string]int

func NewVectorClock(replicaID string) VectorClock {
	if (replicaID==""){
		return VectorClock{}
	}
	return VectorClock{replicaID: 1}
}
func (v VectorClock) NewVectorClock(replicaID string) VectorClock { // for now, constant size increment as number of concurrent users is minimal (<10)
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
	str1 := vc.toString()
	str2 := vc2.toString()
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

func (vc1 VectorClock) Merge(vc2 VectorClock) VectorClock {
	return mergeMaps(vc1, vc2)
}

func mergeMaps(m1 map[string]int, m2 map[string]int) map[string]int {
	union := make(map[string]int)
	for k, val := range m1 {
		union[k] = val
	}
	for k, val := range m2 {
		union[k] = max(val, union[k])
	}
	return union
}

func copyMap(m map[string]int) map[string]int {
	newM := make(map[string]int)
	for k, v := range m {
		newM[k] = v
	}
	return newM
}
func toKeysArray(m map[string]int) []string {
	results := make([]string, len(m))
	i := 0
	for k, _ := range m {
		results[i] = k
		i++
	}
	return results
}
func (vc VectorClock) toString() string {
	var result string
	keys := toKeysArray(vc)
	slices.Sort(keys)
	for _, v := range keys {
		result += (v + "-" + strconv.Itoa(vc[v]) + "-")
	}
	return result

}
