package vector_clock_test

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/AbdelrahmanWM/SyncVerse/document/crdt/internal/vector_clock"
)

func TestNewVectorClock(t *testing.T) {
	testCases := []struct {
		v         VectorClock
		replicaID string
		newV      VectorClock
	}{
		{
			VectorClock{"A": 1, "B": 1},
			"A",
			VectorClock{"A": 2, "B": 1},
		},
		{
			VectorClock{"A": 1, "B": 1},
			"B",
			VectorClock{"A": 1, "B": 2},
		},
		{
			VectorClock{"A": 2, "B": 1},
			"C",
			VectorClock{"A": 2, "B": 1, "C": 1},
		},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("NewVectorClock %v->%v", testCase.v, testCase.replicaID), func(t *testing.T) {
			got := testCase.v.NewVectorClock(testCase.replicaID)
			if !reflect.DeepEqual(got, testCase.newV) {
				t.Errorf("expected %v, got %v", testCase.newV, got)
			}
		})
	}
}
func TestCompare(t *testing.T) {
	testCases := []struct {
		v1       VectorClock
		v2       VectorClock
		expected int
	}{
		{
			v1:       VectorClock{"A": 1, "B": 2},
			v2:       VectorClock{"A": 1, "B": 1},
			expected: 1, // v1 is ahead on B
		},
		{
			v1:       VectorClock{"A": 1, "B": 2},
			v2:       VectorClock{"A": 2, "B": 2},
			expected: -1, // v2 is ahead on A
		},
		{
			v1:       VectorClock{"A": 1, "B": 2, "C": 1},
			v2:       VectorClock{"A": 1, "B": 2, "C": 2},
			expected: -1, // v2 is ahead on C
		},
		{
			v1:       VectorClock{"A": 1, "B": 2},
			v2:       VectorClock{"A": 1, "B": 2},
			expected: 0, // v1 and v2 are equal
		},
		{
			v1:       VectorClock{"A": 2, "B": 3},
			v2:       VectorClock{"A": 1, "B": 3},
			expected: 1, // v1 is ahead on A
		},
		{
			v1:       VectorClock{"A": 1, "B": 1, "C": 3},
			v2:       VectorClock{"A": 1, "B": 1, "C": 2},
			expected: 1, // v1 is ahead on C
		},
		{
			v1:       VectorClock{"A": 3, "B": 3, "C": 3},
			v2:       VectorClock{"A": 2, "B": 2, "C": 2},
			expected: 1, // v1 is ahead on all sequence numbers
		},
		{
			v1:       VectorClock{"A": 2, "B": 4},
			v2:       VectorClock{"A": 4, "B": 2},
			expected: 0, // concurrent
		},
		{
			v1:       VectorClock{"A": 1, "B": 2, "C": 1},
			v2:       VectorClock{"A": 1, "B": 2},
			expected: 1, // v1 has an extra replica C
		},
		{
			v1:       VectorClock{"A": 1, "B": 1},
			v2:       VectorClock{"A": 1, "B": 1, "C": 1},
			expected: -1, // v2 has an extra replica C
		},
		{
			v1:       VectorClock{"A": 1, "B": 2, "C": 1},
			v2:       VectorClock{"A": 1, "B": 2, "D": 1},
			expected: 0, // different branches (concurrent)
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Compare(%v,%v)", testCase.v1, testCase.v2), func(t *testing.T) {
			got := testCase.v1.Compare(testCase.v2)
			if got != testCase.expected {
				t.Errorf("expected %d, got %d", testCase.expected, got)
			}
		})
	}
}

func TestMerge(t *testing.T) {
	testCases := []struct {
		v1     VectorClock
		v2     VectorClock
		VMerge VectorClock
	}{
		// Simple merge, both have the same replica IDs but different sequence numbers
		{
			v1:     VectorClock{"A": 1, "B": 2},
			v2:     VectorClock{"A": 2, "B": 1},
			VMerge: VectorClock{"A": 2, "B": 2},
		},
		// v2 has an additional replica
		{
			v1:     VectorClock{"A": 1, "B": 1},
			v2:     VectorClock{"A": 2, "B": 1, "C": 1},
			VMerge: VectorClock{"A": 2, "B": 1, "C": 1},
		},
		// v1 has an additional replica
		{
			v1:     VectorClock{"A": 1, "B": 2, "C": 1},
			v2:     VectorClock{"A": 2, "B": 2},
			VMerge: VectorClock{"A": 2, "B": 2, "C": 1},
		},
		// v1 and v2 are identical
		{
			v1:     VectorClock{"A": 1, "B": 2},
			v2:     VectorClock{"A": 1, "B": 2},
			VMerge: VectorClock{"A": 1, "B": 2},
		},
		// v2 has a higher sequence number for one replica
		{
			v1:     VectorClock{"A": 1, "B": 1},
			v2:     VectorClock{"A": 1, "B": 3},
			VMerge: VectorClock{"A": 1, "B": 3},
		},
		// v2 has a lower sequence number for one replica (v1 is ahead)
		{
			v1:     VectorClock{"A": 3, "B": 2},
			v2:     VectorClock{"A": 1, "B": 2},
			VMerge: VectorClock{"A": 3, "B": 2},
		},
		// v1 and v2 are divergent
		{
			v1:     VectorClock{"A": 1, "B": 2, "C": 3},
			v2:     VectorClock{"A": 1, "B": 3, "D": 1},
			VMerge: VectorClock{"A": 1, "B": 3, "C": 3, "D": 1},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Merge(%v, %v)", tc.v1, tc.v2), func(t *testing.T) {
			got := tc.v1.Merge(tc.v2)
			if !reflect.DeepEqual(got, tc.VMerge) {
				t.Errorf("expected %+v, got %+v", tc.VMerge, tc.v1)
			}
		})
	}
}
func TestEquals(t *testing.T) {
	testCases := []struct {
		v1   VectorClock
		v2   VectorClock
		want bool
	}{
		{
			VectorClock{"A": 0},
			VectorClock{"A": 0},
			true,
		},
		{
			VectorClock{"A": 1},
			VectorClock{"A": 0},
			false,
		},
		{
			VectorClock{"A": 1, "B": 2},
			VectorClock{"A": 1, "B": 2},
			true,
		},
		{
			VectorClock{"A": 1, "B": 2},
			VectorClock{"A": 1, "B": 3},
			false,
		},
		{
			VectorClock{"A": 1},
			VectorClock{"A": 1, "B": 0},
			false,
		},
		{
			VectorClock{},
			VectorClock{},
			true,
		},
		{
			VectorClock{"A": 1},
			VectorClock{},
			false,
		},
	}
	for _, ts := range testCases {
		got := ts.v1.Equals(ts.v2)
		if got != ts.want {
			t.Errorf("got %v, want %v, v1: %v, v2: %v", got, ts.want, ts.v1, ts.v2)
		}
	}
}
func TestCompareHashes(t *testing.T) {
	testCases := []struct {
		v1   VectorClock
		v2   VectorClock
		want int
	}{
		{
			VectorClock{"A": 0},
			VectorClock{"A": 0},
			0,
		},
		{
			VectorClock{"A": 1},
			VectorClock{"A": 0},
			1,
		},
		{
			VectorClock{"A": 1, "B": 2},
			VectorClock{"A": 1, "B": 2},
			0,
		},
		{
			VectorClock{"A": 1, "B": 2},
			VectorClock{"A": 1, "B": 3},
			-1,
		},
		{
			VectorClock{"A": 1},
			VectorClock{"A": 1, "B": 0},
			-1,
		},
		{
			VectorClock{},
			VectorClock{},
			0,
		},
		{
			VectorClock{"A": 1, "B": 1, "C": 2},
			VectorClock{"B": 1, "C": 2, "D": 2},
			-1,
		},
	}
	for _, ts := range testCases {
		got := ts.v1.CompareHashes(ts.v2)
		if got != ts.want {
			t.Errorf("expected %v, got %v", ts.want, got)
		}
	}
}
