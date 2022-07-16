package goqueue_tests

import (
	"sync"
	"testing"
	"time"

	goqueue "github.com/antonio-alexander/go-queue"

	"github.com/stretchr/testify/assert"
)

const casef string = "case: %s"

//REVIEW: implement tests for sanity/security checks
// * When using dequeue methods that output slices, can we ensure we don't accidentally leak the
//   underlying slice? This should be possible using runtime garbage collection and total/allocated
//   heap memory.
// func getHeap() (allocated, totalAllocated uint64) {
// 	//via https://golangcode.com/print-the-current-memory-usage/
// 	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
// 	m := runtime.MemStats{}
// 	runtime.ReadMemStats(&m)
// 	return m.Alloc, m.TotalAlloc
// }

//KIM: this test doesn't directly do a good job of confirming garbage collection
// due to the difficulty of figuring out how much data in the heap was present
// before and after garbage collection
func TestGarbageCollect(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.GarbageCollecter
	goqueue.Enqueuer
	goqueue.Dequeuer
}) func(*testing.T) {
	return func(t *testing.T) {
		cases := map[string]struct {
			iSize    int
			iEnqueue int
		}{
			"normal": {
				iSize:    100,
				iEnqueue: 100,
			},
		}
		for cDesc, c := range cases {
			q := newQueue(c.iSize)
			for i := 0; i < c.iEnqueue; i++ {
				overflow := q.Enqueue(&goqueue.Example{Int: i})
				if !assert.False(t, overflow, casef, cDesc) {
					break
				}
			}
			q.GarbageCollect()
			for i := 0; i < c.iEnqueue; i++ {
				value, underflow := goqueue.ExampleDequeue(q)
				if assert.False(t, underflow, casef, cDesc) {
					assert.Equal(t, i, value, casef, cDesc)
				}
			}
			q.Close()
		}
	}
}

func TestDequeue(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
	goqueue.Length
}) func(*testing.T) {
	return func(t *testing.T) {
		cases := map[string]struct {
			iSize      int
			iExamples  []*goqueue.Example
			oUnderflow []bool
			oExamples  []*goqueue.Example
		}{
			"empty": {
				iSize:      5,
				iExamples:  []*goqueue.Example{},
				oUnderflow: []bool{true},
				oExamples:  []*goqueue.Example{{Int: 1}},
			},
			"single item": {
				iSize:      5,
				iExamples:  []*goqueue.Example{{Int: 1}},
				oUnderflow: []bool{false},
				oExamples:  []*goqueue.Example{{Int: 1}},
			},
			"full": {
				iSize:      5,
				iExamples:  []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}, {Int: 5}},
				oUnderflow: []bool{false, false, false, false, false, true},
				oExamples:  []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}, {Int: 5}},
			},
		}
		for cDesc, c := range cases {
			q := newQueue(c.iSize)
			for _, item := range c.iExamples {
				overflow := q.Enqueue(item)
				assert.False(t, overflow, casef, cDesc)
			}
			length := q.Length()
			for i, oUnderflow := range c.oUnderflow {
				item, underflow := goqueue.ExampleDequeue(q)
				assert.Equal(t, oUnderflow, underflow)
				if !underflow {
					assert.Equal(t, c.oExamples[i], item, casef, cDesc)
					length--
					assert.Equal(t, length, q.Length(), casef, cDesc)
				}
			}
			q.Close()
		}
	}
}

//TestDequeueMultiple is used to validate the functionality of
// attempting to dequeue multiple items, it does this by pushing
// a number of items into the queue, and then popping them out
// with the idea that the function will always the number of elements
// requested (in the order pushed) or whatever is in the queue
func TestDequeueMultiple(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
}) func(*testing.T) {
	return func(t *testing.T) {
		cases := map[string]struct {
			iSize     int
			iExamples []*goqueue.Example
			iN        int
			oLength   int
		}{
			"zero": {
				iSize:     5,
				iExamples: []*goqueue.Example{{Int: 1}},
				iN:        0,
				oLength:   0,
			},
			"empty queue": {
				iSize:   5,
				oLength: 0,
			},
			"less than in": {
				iSize:     5,
				iExamples: []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}, {Int: 5}},
				iN:        3,
				oLength:   3,
			},
			"equal": {
				iSize:     5,
				iExamples: []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}, {Int: 5}},
				iN:        5,
				oLength:   5,
			},
			"greater than in": {
				iSize:     5,
				iExamples: []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}, {Int: 5}},
				iN:        10,
				oLength:   5,
			},
		}
		for cDesc, c := range cases {
			q := newQueue(c.iSize)
			for _, item := range c.iExamples {
				overflow := q.Enqueue(item)
				assert.False(t, overflow, casef, cDesc)
			}
			values := goqueue.ExampleDequeueMultiple(q, c.iN)
			if assert.Len(t, values, c.oLength, casef, cDesc) {
				for i, value := range values {
					assert.Equal(t, c.iExamples[i], value, casef, cDesc)
				}
			}
			q.Close()
		}
	}
}

func TestFlush(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
}) func(*testing.T) {
	return func(t *testing.T) {
		cases := map[string]struct {
			iSize     int
			iExamples []*goqueue.Example
		}{
			"empty": {
				iSize: 1,
			},
			"single": {
				iSize:     5,
				iExamples: []*goqueue.Example{{Int: 1}},
			},
			"max 5": {
				iSize:     5,
				iExamples: []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}, {Int: 5}},
			},
			"max 1": {
				iSize:     1,
				iExamples: []*goqueue.Example{{Int: 1}},
			},
		}
		for cDesc, c := range cases {
			q := newQueue(c.iSize)
			remainingElements, overflow := goqueue.ExampleEnqueueMultiple(q, c.iExamples)
			//KIM: this is testing dequeue, so this should always succeed
			assert.False(t, overflow, casef, cDesc)
			assert.Len(t, remainingElements, 0, casef, cDesc)
			values := goqueue.ExampleFlush(q)
			if assert.Equal(t, len(values), len(c.iExamples), casef, cDesc) {
				for i, value := range values {
					assert.Equal(t, c.iExamples[i], value, casef, cDesc)
				}
			}
			q.Close()
		}
	}
}

func TestPeek(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
	goqueue.Peeker
}) func(*testing.T) {
	return func(t *testing.T) {
		cases := map[string]struct {
			iSize      int
			iExamples  []*goqueue.Example
			oUnderflow []bool
		}{
			"No Elements": {
				iSize:      5,
				iExamples:  []*goqueue.Example{},
				oUnderflow: []bool{true},
			},
			"Single Element": {
				iSize:      5,
				iExamples:  []*goqueue.Example{{Int: 1}},
				oUnderflow: []bool{false},
			},
			"Multiple Elements": {
				iSize:      5,
				iExamples:  []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}, {Int: 5}},
				oUnderflow: []bool{false},
			},
			"Peek Greater than Size": {
				iSize:      5,
				iExamples:  []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}, {Int: 5}},
				oUnderflow: []bool{false},
			},
		}
		for cDesc, c := range cases {
			q := newQueue(c.iSize)
			for _, input := range c.iExamples {
				overflow := q.Enqueue(input)
				assert.False(t, overflow, casef, cDesc)
			}
			peeked := goqueue.ExamplePeek(q)
			if assert.Equal(t, len(peeked), len(c.iExamples), casef, cDesc) {
				for i, peek := range peeked {
					assert.Equal(t, c.iExamples[i], peek, casef, cDesc)
				}
			}
			value, underflow := goqueue.ExamplePeekHead(q)
			if assert.Equal(t, c.oUnderflow[0], underflow) {
				if !underflow {
					assert.Equal(t, c.iExamples[0], value, casef, cDesc)
				}
			}
			q.Close()
		}
	}
}

func TestPeekFromHead(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
	goqueue.Peeker
}) func(*testing.T) {
	return func(t *testing.T) {
		cases := map[string]struct {
			iSize           int
			iExamples       []*goqueue.Example
			iPeekFromHead   int
			oPeekedExamples []*goqueue.Example
		}{
			"min": {
				iSize:           5,
				iExamples:       []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}, {Int: 5}},
				iPeekFromHead:   1,
				oPeekedExamples: []*goqueue.Example{{Int: 1}},
			},
			"max": {
				iSize:           5,
				iExamples:       []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}, {Int: 5}},
				iPeekFromHead:   5,
				oPeekedExamples: []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}, {Int: 5}},
			},
			"zero": {
				iSize:         5,
				iExamples:     []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}, {Int: 5}},
				iPeekFromHead: 0,
			},
			"max+1": {
				iSize:           5,
				iExamples:       []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}, {Int: 5}},
				iPeekFromHead:   6,
				oPeekedExamples: []*goqueue.Example{{Int: 1}, {Int: 2}, {Int: 3}, {Int: 4}, {Int: 5}},
			},
			"empty queue": {
				iSize:         5,
				iPeekFromHead: 1,
			},
		}
		for cDesc, c := range cases {
			q := newQueue(c.iSize)
			for _, input := range c.iExamples {
				overflow := q.Enqueue(input)
				assert.False(t, overflow, casef, cDesc)
			}
			values := goqueue.ExamplePeekFromHead(q, c.iPeekFromHead)
			if assert.Equal(t, len(values), len(c.oPeekedExamples), casef, cDesc) {
				for i, item := range values {
					assert.Equal(t, c.oPeekedExamples[i], item, casef, cDesc)
				}
			}
			q.Close()
		}
	}
}

func TestEvent(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
	goqueue.Event
}) func(*testing.T) {
	return func(t *testing.T) {
		cases := map[string]struct {
			iSize int
		}{
			"normal": {
				iSize: 1,
			},
		}
		for cDesc, c := range cases {
			var wg sync.WaitGroup

			q := newQueue(c.iSize)
			signalIn, signalOut := q.GetSignalIn(), q.GetSignalOut()
			start := make(chan struct{})
			wg.Add(2)
			go func() {
				defer wg.Done()
				<-start
				overflow := q.Enqueue(&goqueue.Example{})
				assert.False(t, overflow, casef, cDesc)
			}()
			go func() {
				defer wg.Done()
				<-start
				select {
				case <-time.After(time.Second):
					assert.Fail(t, "no signal received when expected", casef, cDesc)
				case <-signalIn:
				}
				select {
				default:
				case <-signalIn:
					assert.Fail(t, "signal received when unexpected", casef, cDesc)
				}
			}()
			close(start)
			wg.Wait()
			start = make(chan struct{})
			wg.Add(2)
			go func() {
				defer wg.Done()
				<-start
				_, underflow := q.Dequeue()
				assert.False(t, underflow, casef, cDesc)
			}()
			go func() {
				defer wg.Done()
				<-start
				select {
				case <-time.After(time.Second):
					assert.Fail(t, "no signal received when expected", casef, cDesc)
				case <-signalOut:
				}
				select {
				default:
				case <-signalOut:
					assert.Fail(t, "signal received when unexpected", casef, cDesc)
				}
			}()
			close(start)
			wg.Wait()
			q.Close()
		}
	}
}

func TestLength(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
	goqueue.Length
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
			q := newQueue(c.iSize)
			overflow := q.Enqueue(&goqueue.Example{})
			assert.False(t, overflow, casef, cDesc)
			assert.Equal(t, 1, q.Length(), casef, cDesc)
			_, underflow := q.Dequeue()
			assert.False(t, underflow, casef, cDesc)
			assert.Equal(t, 0, q.Length(), casef, cDesc)
			for i := 0; i < c.iSize; i++ {
				overflow := q.Enqueue(&goqueue.Example{})
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

//TestQueue
// 1. Use the New() function to create/populate a queue of the size for the case
// 2. Use the Length() function to verify that the queue is empty (size of 0)
// 3. Use the Enqueue() function for the number of itemsIn to place data in the queue and verify that
//  the length increases by one each time.
// 4. Use the Length() to check to see if the queue is the expected size
// 5. Use the Dequeue() Function again to verify an underflow as the queue should now be empty (length of 0)
// 6. Use the Close() function to clean up all internal pointers for the queue
func TestQueue(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
	goqueue.Length
}) func(*testing.T) {
	return func(t *testing.T) {
		randFloats := goqueue.ExampleGenFloat64(0)
		singleFloat := goqueue.ExampleGenFloat64(1)
		cases := map[string]struct {
			iSize   int
			iFloats []*goqueue.Example
		}{
			"random length": {
				iSize:   len(randFloats),
				iFloats: randFloats,
			},
			"negative length": {
				iSize:   -1,
				iFloats: singleFloat,
			},
			"zero length": {
				iSize:   0,
				iFloats: singleFloat,
			},
		}
		for cDesc, c := range cases {
			q := newQueue(c.iSize)
			assert.Equal(t, 0, q.Length(), casef, cDesc)
			for _, item := range c.iFloats {
				overflow := q.Enqueue(item)
				assert.False(t, overflow, casef, cDesc)
			}
			assert.Equal(t, len(c.iFloats), q.Length(), casef, cDesc)
			for i := 0; i < len(c.iFloats); i++ {
				value, underflow := goqueue.ExampleDequeue(q)
				if assert.False(t, underflow, casef, cDesc) {
					assert.Equal(t, c.iFloats[i], value, casef, cDesc)
				}
			}
			_, underflow := q.Dequeue()
			assert.True(t, underflow, casef, cDesc)
			q.Close()
		}
	}
}

//TestAsync
// 1. Populate async interface using New()
// 2. Create two goRoutines:
//   a. goRoutine (dequeue):
//    (2) Stop when signal received after enqueue function is finished enqueing data
//    (1) Constantly attempt to dequeue, when underflow is false, add item to slice of float64
//   b. goRoutine (enqueue):
//    (1) Enqueue all the data from randFloats, store data in queue
//    (2) Send signal when finished enqueuing data
// 3. Compare the items dequeued to the items enqueued, they should be equal although their quantity may not be the same (see verification)
func TestAsync(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
	goqueue.Length
}) func(*testing.T) {
	return func(t *testing.T) {
		randFloats := goqueue.ExampleGenFloat64(0)
		cases := map[string]struct {
			Size   int
			Floats []*goqueue.Example
			Ints   []*goqueue.Example
		}{
			"basic": {
				Size:   len(randFloats),
				Floats: randFloats,
			},
		}
		for cDesc, c := range cases {
			var valuesEnqueued, valuesDequeued []*goqueue.Example
			var wg sync.WaitGroup

			q := newQueue(c.Size)
			stopDequeue := make(chan (struct{}))
			wg.Add(1)
			go func() {
				defer wg.Done()

				for {
					select {
					case <-stopDequeue:
						return
					default:
						if value, underflow := goqueue.ExampleDequeue(q); !underflow {
							valuesDequeued = append(valuesDequeued, value)
						}
					}
				}
			}()
			wg.Add(1)
			go func() {
				defer wg.Done()

				for _, item := range c.Floats {
					if overflow := q.Enqueue(item); !overflow {
						valuesEnqueued = append(valuesEnqueued, item)
					}
				}
				close(stopDequeue)
			}()
			wg.Wait()
			for i, value := range valuesDequeued {
				assert.Equal(t, valuesEnqueued[i], value, casef, cDesc)
			}
			q.Close()
		}
	}
}
