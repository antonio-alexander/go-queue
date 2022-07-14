package finite_tests

import (
	"testing"

	goqueue "github.com/antonio-alexander/go-queue"
	finite "github.com/antonio-alexander/go-queue/finite"

	"github.com/stretchr/testify/assert"
)

const casef string = "case: %s"

func TestNew(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	finite.Capacity
}) func(*testing.T) {
	return func(t *testing.T) {
		cases := map[string]struct {
			iSize     int
			oCapacity int
		}{
			"normal; size > 1": {
				iSize:     10,
				oCapacity: 10,
			},
			"failure; size < 0": {
				iSize:     0,
				oCapacity: 1,
			},
		}
		for cDesc, c := range cases {
			//create the queue using the newQueue function, assert
			// that it is not nil and validate that the capacity is
			// as expected, then close the queue
			q := newQueue(c.iSize)
			if assert.NotNil(t, q) {
				if capacity := q.Capacity(); capacity != c.oCapacity {
					assert.Equal(t, c.oCapacity, casef, cDesc)
				}
				q.Close()
			}
		}
	}
}

func TestResize(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	finite.Capacity
	finite.Resizer
}) func(*testing.T) {
	return func(t *testing.T) {
		cases := map[string]struct {
			iSize            int
			iNewSize         int
			iExamples        []*goqueue.Example
			oSize            int
			oRemovedExamples []*goqueue.Example
		}{
			"same": {
				iSize:    1,
				iNewSize: 1,
				oSize:    1,
			},
			"greater": {
				iSize:     1,
				iNewSize:  5,
				iExamples: []*goqueue.Example{{Int: 1}},
				oSize:     5,
			},
			"less": {
				iSize:            5,
				iNewSize:         1,
				iExamples:        []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}, {Int: 5}},
				oSize:            1,
				oRemovedExamples: []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}},
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
			for _, value := range c.iExamples {
				overflow := q.Enqueue(value)
				assert.False(t, overflow, "unxpected overflow", cDesc)
			}
			removedExamples := finite.ExampleResize(q, c.iNewSize)
			size := q.Capacity()
			assert.Equal(t, c.oSize, size, "new queue size", cDesc)
			assert.Equal(t, len(c.oRemovedExamples), len(removedExamples))
			for i, removedExample := range removedExamples {
				assert.Equal(t, c.oRemovedExamples[i], removedExample, "elements removed value", cDesc)
			}
			q.Close()
		}
	}
}

func TestCapacity(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
	goqueue.Length
	finite.Capacity
}) func(*testing.T) {
	return func(t *testing.T) {
		cases := map[string]struct {
			iSize int
		}{
			"one": {
				iSize: 1,
			},
			"two": {
				iSize: 2,
			},
			"ten": {
				iSize: 10,
			},
			"hundred": {
				iSize: 100,
			},
		}
		for cDesc, c := range cases {
			//TODO: add documentation
			q := newQueue(c.iSize)
			capacity := q.Capacity()
			assert.Equal(t, c.iSize, capacity, casef, cDesc)
			overflow := q.Enqueue(struct{}{})
			assert.False(t, overflow, casef, cDesc)
			assert.Equal(t, 1, q.Length(), casef, cDesc)
			_, underflow := q.Dequeue()
			assert.False(t, underflow, casef, cDesc)
			assert.Equal(t, 0, q.Length(), casef, cDesc)
			for i := 0; i < c.iSize; i++ {
				overflow := q.Enqueue(struct{}{})
				assert.False(t, overflow, casef, cDesc)
			}
			assert.Equal(t, c.iSize, q.Length(), casef, cDesc)
			for i := 0; i < c.iSize; i++ {
				_, underflow := q.Dequeue()
				assert.False(t, underflow, casef, cDesc)
			}
			assert.Equal(t, 0, q.Length(), casef, cDesc)
			q.Close()
		}
	}
}

func TestEnqueue(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
}) func(*testing.T) {
	return func(t *testing.T) {
		cases := map[string]struct {
			iSize      int
			iExamples  []*goqueue.Example
			oOverflows []bool
		}{
			"min": {
				iSize:      1,
				iExamples:  []*goqueue.Example{{Int: 1}, {Int: 2}},
				oOverflows: []bool{false, true},
			},
			"max-1": {
				iSize:      5,
				iExamples:  []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}},
				oOverflows: []bool{false, false, false, false},
			},
			"min+2": {
				iSize:      5,
				iExamples:  []*goqueue.Example{{Int: 1}, {Int: 2}},
				oOverflows: []bool{false, false},
			},
			"max+1": {
				iSize:      5,
				iExamples:  []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}, {Int: 5}, {Int: 6}},
				oOverflows: []bool{false, false, false, false, false, true},
			},
			"max+2": {
				iSize:      5,
				iExamples:  []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}, {Int: 5}, {Int: 6}, {Int: 7}},
				oOverflows: []bool{false, false, false, false, false, true, true},
			},
		}
		for cDesc, c := range cases {
			q := newQueue(c.iSize)
			for i, element := range c.iExamples {
				overflow := q.Enqueue(element)
				assert.Equal(t, c.oOverflows[i], overflow, casef, cDesc)
			}
			q.Close()
		}
	}
}

func TestEnqueueMultiple(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
}) func(*testing.T) {
	return func(t *testing.T) {
		cases := map[string]struct {
			iSize             int
			iExamples         []*goqueue.Example
			oOverflow         bool
			oOverflowExamples []*goqueue.Example
		}{
			"max-1": {
				iSize:     5,
				iExamples: []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}},
				oOverflow: false,
			},
			"min+2": {
				iSize:     5,
				iExamples: []*goqueue.Example{{Int: 1}, {Int: 2}},
				oOverflow: false,
			},
			"max+1": {
				iSize:             5,
				iExamples:         []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}, {Int: 5}, {Int: 6}},
				oOverflow:         true,
				oOverflowExamples: []*goqueue.Example{{Int: 6}},
			},
			"max+2": {
				iSize:             5,
				iExamples:         []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}, {Int: 5}, {Int: 6}, {Int: 7}},
				oOverflow:         true,
				oOverflowExamples: []*goqueue.Example{{Int: 6}, {Int: 7}},
			},
		}
		for cDesc, c := range cases {
			q := newQueue(c.iSize)
			values, overflow := goqueue.ExampleEnqueueMultiple(q, c.iExamples)
			if assert.Equal(t, c.oOverflow, overflow, casef, cDesc) {
				if overflow {
					for i, value := range values {
						assert.Equal(t, c.oOverflowExamples[i], value, casef, cDesc)
					}
				}
			}
			q.Close()
		}
	}
}

func TestEnqueueInFront(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.EnqueueInFronter
	goqueue.Peeker
}) func(*testing.T) {
	return func(t *testing.T) {
		cases := map[string]struct {
			iSize     int
			iInFront  *goqueue.Example
			iExamples []*goqueue.Example
			oOverflow bool
		}{
			"empty_queue": {
				iSize:     5,
				iInFront:  &goqueue.Example{Int: 10},
				oOverflow: false,
			},
			"queue_not_full": {
				iSize:     5,
				iExamples: []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}},
				iInFront:  &goqueue.Example{Int: 10},
				oOverflow: false,
			},
			"queue_full": {
				iSize:     5,
				iExamples: []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}, {Int: 5}},
				iInFront:  &goqueue.Example{Int: 88},
				oOverflow: true,
			},
			"almost_full": {
				iSize:     5,
				iExamples: []*goqueue.Example{{Int: 1}},
				iInFront:  &goqueue.Example{Int: 66},
				oOverflow: false,
			},
			"max-1": {
				iSize:     5,
				iExamples: []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}},
				iInFront:  &goqueue.Example{Int: 44},
				oOverflow: false,
			},
		}
		for cDesc, c := range cases {
			//TODO: add documentation
			q := newQueue(c.iSize)
			for _, element := range c.iExamples {
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
}

func TestEnqueueLossy(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	finite.EnqueueLossy
}) func(*testing.T) {
	return func(t *testing.T) {
		cases := map[string]struct {
			iSize            int
			iExamples        []*goqueue.Example
			oDiscard         []bool
			oDiscardExamples []*goqueue.Example
		}{
			"single_element": {
				iExamples:        []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}},
				oDiscard:         []bool{false, true, true},
				oDiscardExamples: []*goqueue.Example{{Int: 0}, {Int: 1}, {Int: 2}},
			},
			"normal_enqueue": {
				iSize:     5,
				iExamples: []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}, {Int: 5}},
				oDiscard:  []bool{false, false, false, false, false},
			},
			"lossy_enqueue": {
				iSize:            5,
				iExamples:        []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}, {Int: 5}, {Int: 6}, {Int: 7}},
				oDiscard:         []bool{false, false, false, false, false, true, true},
				oDiscardExamples: []*goqueue.Example{{Int: 0}, {Int: 0}, {Int: 0}, {Int: 0}, {Int: 0}, {Int: 1}, {Int: 2}},
			},
		}
		for cDesc, c := range cases {
			q := newQueue(c.iSize)
			for i, value := range c.iExamples {
				value, discard := finite.ExampleEnqueueLossy(q, value)
				if assert.Equal(t, c.oDiscard[i], discard, casef, cDesc) {
					if discard {
						assert.Equal(t, value, c.oDiscardExamples[i], casef, cDesc)
					}
				}
			}
			q.Close()
		}
	}
}
