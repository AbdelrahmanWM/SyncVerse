package vector_clock

type VectorClock map[string]int

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
