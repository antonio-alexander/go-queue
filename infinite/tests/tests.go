package infinite_tests

import (
	"testing"

	goqueue "github.com/antonio-alexander/go-queue"

	"github.com/stretchr/testify/assert"
)

const casef string = "case: %s"

func TestNew(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Length
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
			q := newQueue(c.iSize)
			if !assert.NotNil(t, q, cDesc) {
				continue
			}
			q.Close()
		}
	}
}

func TestGarbageCollect(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.GarbageCollecter
	goqueue.Enqueuer
	goqueue.Dequeuer
}) func(*testing.T) {
	return func(t *testing.T) {
		cases := map[string]struct {
			iGrowSize   int
			iNtoEnqueue int
			iNtoDequeue int
			oCapacity   int
		}{
			"grow + 1": {
				iGrowSize:   5,
				iNtoEnqueue: 9,
				iNtoDequeue: 3,
				oCapacity:   10,
			},
			"grow - 1": {
				iGrowSize:   5,
				iNtoEnqueue: 10,
				iNtoDequeue: 10,
				oCapacity:   5,
			},
			"grow + 0": {
				iGrowSize:   5,
				iNtoEnqueue: 5,
				iNtoDequeue: 5,
				oCapacity:   5,
			},
		}

		for cDesc, c := range cases {
			q := newQueue(c.iGrowSize)
			for i := 0; i < c.iNtoEnqueue; i++ {
				overflow := q.Enqueue(i)
				assert.False(t, overflow, casef, cDesc)
			}
			for i := 0; i < c.iNtoDequeue; i++ {
				element, underflow := q.Dequeue()
				if assert.False(t, underflow, casef, cDesc) {
					assert.Equal(t, i, element)
				}
			}
			q.GarbageCollect()
			for i := 0 + c.iNtoDequeue; i < c.iNtoEnqueue; i++ {
				element, underflow := q.Dequeue()
				if assert.False(t, underflow, casef, cDesc) {
					assert.Equal(t, i, element)
				}
			}
			q.Close()
		}
	}
}

func TestEnqueue(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Length
}) func(*testing.T) {
	return func(t *testing.T) {
		cases := map[string]struct {
			iGrowSize int
			iValues   []*goqueue.Example
			oLengths  []int
		}{
			"No Grow": {
				iGrowSize: 10,
				iValues:   []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}},
				oLengths:  []int{1, 2, 3},
			},
			"grow once": {
				iGrowSize: 2,
				iValues:   []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}},
				oLengths:  []int{1, 2, 3},
			},
			"grow twice": {
				iGrowSize: 2,
				iValues:   []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}, {Int: 5}, {Int: 6}},
				oLengths:  []int{1, 2, 3, 4, 5, 6},
			},
		}
		for cDesc, c := range cases {
			q := newQueue(c.iGrowSize)
			for i, value := range c.iValues {
				overflow := q.Enqueue(value)
				assert.False(t, overflow, casef, cDesc)
				assert.Equal(t, c.oLengths[i], q.Length())
			}
			q.Close()
		}
	}
}

func TestEnqueueMultiple(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Length
}) func(*testing.T) {
	return func(t *testing.T) {
		cases := map[string]struct {
			iGrowSize int
			iValues   []*goqueue.Example
			oLength   int
		}{
			"No Grow": {
				iGrowSize: 10,
				iValues:   []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}},
				oLength:   3,
			},
			"grow once": {
				iGrowSize: 2,
				iValues:   []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}},
				oLength:   3,
			},
			"grow twice": {
				iGrowSize: 2,
				iValues:   []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}, {Int: 5}, {Int: 6}},
				oLength:   6,
			},
		}
		for cDesc, c := range cases {
			q := newQueue(c.iGrowSize)
			remainingElements, overflow := goqueue.ExampleEnqueueMultiple(q, c.iValues)
			if assert.False(t, overflow, casef, cDesc) &&
				assert.GreaterOrEqual(t, len(remainingElements), 0) {
				length := q.Length()
				assert.Equal(t, c.oLength, length, casef, cDesc)
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
			iSize         int
			iInFrontValue *goqueue.Example
			iInts         []*goqueue.Example
		}{
			"empty_queue": { //verify that eif doesn't overflow if queue is empty
				iSize:         5,
				iInFrontValue: &goqueue.Example{Int: 10},
			},
			"queue_not_full": { //verify that eif doesn't overflow if queue is not full
				iSize:         5,
				iInts:         []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}},
				iInFrontValue: &goqueue.Example{Int: 10},
			},
			"queue_full": { //verify that eif overflows if queue is full
				iSize:         5,
				iInts:         []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}, {Int: 5}},
				iInFrontValue: &goqueue.Example{Int: 88},
			},
			"almost_full": { //verify that eif places data in front if you enqueue one element
				iSize:         5,
				iInts:         []*goqueue.Example{{Int: 1}},
				iInFrontValue: &goqueue.Example{Int: 66},
			},
			"max-1": { //verify that eif places data in front if you enqueue max -1 number of elements
				iSize:         5,
				iInts:         []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}},
				iInFrontValue: &goqueue.Example{Int: 44},
			},
		}
		for cDesc, c := range cases {
			q := newQueue(c.iSize)
			for _, element := range c.iInts {
				overflow := q.Enqueue(element)
				assert.False(t, overflow, casef, cDesc)
			}
			overflow := q.EnqueueInFront(c.iInFrontValue)
			if assert.False(t, overflow, casef, cDesc) {
				element, underflow := q.PeekHead()
				if assert.False(t, underflow, casef, cDesc) {
					assert.Equal(t, c.iInFrontValue, element, casef, cDesc)
				}
			}
			q.Close()
		}
	}
}

func TestQueue(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
	goqueue.Length
}) func(*testing.T) {
	return func(t *testing.T) {
		randFloats := goqueue.ExampleGenFloat64(0)
		cases := map[string]struct {
			iSize   int
			iValues []*goqueue.Example
		}{
			"Ticker": { //test with ticker
				iSize:   len(randFloats),
				iValues: randFloats,
			},
		}

		for cDesc, c := range cases {
			//1. Use the New() function to create/populate a queue of the size for the case
			q := newQueue(c.iSize)
			for i := 0; i < len(c.iValues)-1; i++ {
				//2. Use the Enqueue() function for the number of elementsIn to place data in the queue and verify that the length increases by one each time.
				overflow := q.Enqueue(c.iValues[i])
				assert.False(t, overflow, casef, cDesc)
			}
			assert.Equal(t, q.Length(), len(c.iValues)-1)
			//3. Use the Enqueue() function again to verify that the queue doesn't overflow and that the capacity increases (capacity after enqueue is greater than the capacity before)
			overflow := q.Enqueue(c.iValues[len(c.iValues)-1])
			assert.False(t, overflow, casef, cDesc)
			assert.Equal(t, len(c.iValues), q.Length())
			for i := 0; i < c.iSize; i++ {
				//4. Use the Dequeue() function for the number of elementsIn to remove data from the queue adn verify that the length decreases by one each time.
				value, underflow := goqueue.ExampleDequeue(q)
				if assert.False(t, underflow, casef, cDesc) {
					assert.Equal(t, c.iValues[i], value)
				}
			}
			//5. Use the Dequeue() Function again to verify an underflow as the queue should now be empty (length of 0)
			_, underflow := q.Dequeue()
			assert.True(t, underflow, casef, cDesc)
			//6. Use the Close() function to clean up all internal pointers for the queue
			q.Close()
		}
	}
}

//REVIEW: implement tests for sanity/security checks
// * When using dequeue methods that output slices, can we ensure we don't accidentally leak the
//   underlying slice? This should be possible using runtime garbage collection and total/allocated
//   heap memory.
