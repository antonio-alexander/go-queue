package finite_tests

import (
	"testing"

	goqueue "github.com/antonio-alexander/go-queue"
	finite "github.com/antonio-alexander/go-queue/finite"

	"github.com/stretchr/testify/assert"
)

const casef string = "case: %s"

func Resize(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Info
	finite.Resizer
}) {

	cases := map[string]struct {
		iSize            int
		iNewSize         int
		iInts            []int
		oSize            int
		oRemovedElements []interface{}
	}{
		"same": {
			iSize:    1,
			iNewSize: 1,
			oSize:    1,
		},
		"greater": {
			iSize:    1,
			iNewSize: 5,
			iInts:    []int{1},
			oSize:    5,
		},
		"less": {
			iSize:            5,
			iNewSize:         1,
			iInts:            []int{1, 2, 3, 4, 5},
			oSize:            1,
			oRemovedElements: []interface{}{1, 2, 3, 4},
		},
		"invalid": {
			iSize:    1,
			iNewSize: -1,
			oSize:    1,
		},
	}
	for cDesc, c := range cases {
		//TODO: add documentation
		q := newQueue(c.iSize)
		for _, element := range c.iInts {
			overflow := q.Enqueue(element)
			assert.False(t, overflow, "unxpected overflow", cDesc)
		}
		removedElements := q.Resize(c.iNewSize)
		size := q.Capacity()
		assert.Equal(t, c.oSize, size, "new queue size", cDesc)
		assert.Equal(t, len(c.oRemovedElements), len(removedElements))
		for i, removedElement := range removedElements {
			assert.Equal(t, c.oRemovedElements[i], removedElement.(int), "elements removed value", cDesc)
		}
		q.Close()
	}
}

func Enqueue(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
}) {

	cases := map[string]struct {
		iSize      int
		iElements  []interface{}
		oOverflows []bool
	}{
		"min": {
			iSize:      1,
			iElements:  []interface{}{1, 2},
			oOverflows: []bool{false, true},
		},
		"max-1": {
			iSize:      5,
			iElements:  []interface{}{1, 2, 3, 4},
			oOverflows: []bool{false, false, false, false},
		},
		"min+2": {
			iSize:      5,
			iElements:  []interface{}{1, 2},
			oOverflows: []bool{false, false},
		},
		"max+1": {
			iSize:      5,
			iElements:  []interface{}{1, 2, 3, 4, 5, 6},
			oOverflows: []bool{false, false, false, false, false, true},
		},
		"max+2": {
			iSize:      5,
			iElements:  []interface{}{1, 2, 3, 4, 5, 6, 7},
			oOverflows: []bool{false, false, false, false, false, true, true},
		},
	}
	for cDesc, c := range cases {
		q := newQueue(c.iSize)
		for i, element := range c.iElements {
			overflow := q.Enqueue(element)
			assert.Equal(t, c.oOverflows[i], overflow, casef, cDesc)
		}
		q.Close()
	}
}

func EnqueueMultiple(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
}) {

	cases := map[string]struct {
		iSize             int
		iInterfaces       []interface{}
		oOverflow         bool
		oOverflowElements []interface{}
	}{
		"max-1": {
			iSize:       5,
			iInterfaces: []interface{}{1, 2, 3, 4},
			oOverflow:   false,
		},
		"min+2": {
			iSize:       5,
			iInterfaces: []interface{}{1, 2},
			oOverflow:   false,
		},
		"max+1": {
			iSize:             5,
			iInterfaces:       []interface{}{1, 2, 3, 4, 5, 6},
			oOverflow:         true,
			oOverflowElements: []interface{}{6},
		},
		"max+2": {
			iSize:             5,
			iInterfaces:       []interface{}{1, 2, 3, 4, 5, 6, 7},
			oOverflow:         true,
			oOverflowElements: []interface{}{6, 7},
		},
	}
	for cDesc, c := range cases {
		q := newQueue(c.iSize)
		elements, overflow := q.EnqueueMultiple(c.iInterfaces)
		if assert.Equal(t, c.oOverflow, overflow, casef, cDesc) {
			if overflow {
				for i, element := range elements {
					value, _ := element.(int)
					assert.Equal(t, c.oOverflowElements[i], value, casef, cDesc)
				}
			}
		}
		q.Close()
	}
}

func EnqueueInFront(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.EnqueueInFronter
	goqueue.Peeker
}) {

	cases := map[string]struct {
		iSize     int
		iInFront  int
		iInts     []int
		oOverflow bool
	}{
		"empty_queue": {
			iSize:     5,
			iInFront:  10,
			oOverflow: false,
		},
		"queue_not_full": {
			iSize:     5,
			iInts:     []int{1, 2, 3, 4},
			iInFront:  10,
			oOverflow: false,
		},
		"queue_full": {
			iSize:     5,
			iInts:     []int{1, 2, 3, 4, 5},
			iInFront:  88,
			oOverflow: true,
		},
		"almost_full": {
			iSize:     5,
			iInts:     []int{1},
			iInFront:  66,
			oOverflow: false,
		},
		"max-1": {
			iSize:     5,
			iInts:     []int{1, 2, 3, 4},
			iInFront:  44,
			oOverflow: false,
		},
	}
	for cDesc, c := range cases {
		//TODO: add documentation
		q := newQueue(c.iSize)
		for _, element := range c.iInts {
			overflow := q.Enqueue(element)
			assert.False(t, overflow, casef, cDesc)
		}
		overflow := q.EnqueueInFront(c.iInFront)
		if assert.Equal(t, c.oOverflow, overflow, casef, cDesc) {
			if !overflow {
				element, underflow := q.PeekHead()
				if assert.False(t, underflow, casef, cDesc) {
					assert.Equal(t, c.iInFront, element)
				}
			}
		}
		q.Close()
	}
}

func EnqueueLossy(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	finite.EnqueueLossy
}) {

	cases := map[string]struct {
		iSize            int
		iInts            []int
		oDiscard         []bool
		oDiscardElements []interface{}
	}{
		"single_element": {
			iInts: []int{
				1, 2, 3},
			oDiscard:         []bool{false, true, true},
			oDiscardElements: []interface{}{0, 1, 2},
		},
		"normal_enqueue": {
			iSize:    5,
			iInts:    []int{1, 2, 3, 4, 5},
			oDiscard: []bool{false, false, false, false, false},
		},
		"lossy_enqueue": {
			iSize:            5,
			iInts:            []int{1, 2, 3, 4, 5, 6, 7},
			oDiscard:         []bool{false, false, false, false, false, true, true},
			oDiscardElements: []interface{}{0, 0, 0, 0, 0, 1, 2},
		},
	}
	for cDesc, c := range cases {
		q := newQueue(c.iSize)
		for i, element := range c.iInts {
			element, discard := q.EnqueueLossy(element)
			if assert.Equal(t, c.oDiscard[i], discard, casef, cDesc) {
				if discard {
					assert.Equal(t, element, c.oDiscardElements[i], casef, cDesc)
				}
			}
		}
		q.Close()
	}
}
