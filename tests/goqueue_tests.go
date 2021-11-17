package goqueue_tests

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	goqueue "github.com/antonio-alexander/go-queue"

	"github.com/stretchr/testify/assert"
)

const casef string = "case: %s"

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

//genFloat64 will generate a random number of random float values if n is equal to 0
// not to exceed the constant TestMaxValues, if n is provided, it will generate that many items
func genFloat64(n int) (floats []float64) {
	if n <= 0 {
		n = int(rand.Float64() * 1000)
	}
	for i := 0; i < n; i++ {
		floats = append(floats, rand.Float64())
	}

	return
}

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

func New(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Info
}) {

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

func GarbageCollect(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.GarbageCollecter
	goqueue.Enqueuer
	goqueue.Dequeuer
	goqueue.Info
}) {

	cases := map[string]struct {
		iSize    int
		iEnqueue int
	}{
		"normal": {
			iSize: 100,
		},
		//TODO: full
		//TODO: empty
		//TODO: full-1
		//TODO: empty+1
	}

	for cDesc, c := range cases {
		q := newQueue(c.iSize)
		for i := 0; i < c.iEnqueue; i++ {
			overflow := q.Enqueue(i)
			if !assert.False(t, overflow, casef, cDesc) {
				break
			}
		}
		q.GarbageCollect()
		for i := 0; i < c.iEnqueue; i++ {
			item, underflow := q.Dequeue()
			if assert.False(t, underflow, casef, cDesc) {
				assert.Equal(t, i, item, casef, cDesc)
			}
		}
		q.Close()
	}
}

func Dequeue(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
	goqueue.Info
}) {

	cases := map[string]struct {
		iSize       int
		iInts       []int
		oUnderflow  []bool
		oOutputInts []int
	}{
		"empty": {
			iSize:       5,
			iInts:       []int{},
			oUnderflow:  []bool{true},
			oOutputInts: []int{1},
		},
		"single item": {
			iSize:       5,
			iInts:       []int{1},
			oUnderflow:  []bool{false},
			oOutputInts: []int{1},
		},
		"full": {
			iSize:       5,
			iInts:       []int{1, 2, 3, 4, 5},
			oUnderflow:  []bool{false, false, false, false, false, true},
			oOutputInts: []int{1, 2, 3, 4, 5},
		},
	}
	for cDesc, c := range cases {
		//TODO: add documentation
		q := newQueue(c.iSize)
		capacity := q.Capacity()
		for _, item := range c.iInts {
			overflow := q.Enqueue(item)
			assert.False(t, overflow, casef, cDesc)
		}
		for i, oUnderflow := range c.oUnderflow {
			item, underflow := q.Dequeue()
			if assert.Equal(t, oUnderflow, underflow) && !underflow {
				if assert.NotNil(t, item, casef, cDesc) {
					value, _ := item.(int)
					assert.Equal(t, c.oOutputInts[i], value, casef, cDesc)
				}
			}
			assert.Equal(t, capacity, q.Capacity(), casef, cDesc)
		}
		q.Close()
	}
}

func DequeueMultiple(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
}) {

	cases := map[string]struct {
		iSize            int
		iInts            []int
		iDequeueMultiple int
		oLength          int
	}{
		"zero": {
			iSize:            5,
			iInts:            []int{1},
			iDequeueMultiple: 0,
			oLength:          0,
		},
		"empty queue": {
			iSize:   5,
			oLength: 0,
		},
		"less than in": {
			iSize:            5,
			iInts:            []int{1, 2, 3, 4, 5},
			iDequeueMultiple: 3,
			oLength:          3,
		},
		"equal": {
			iSize:            5,
			iInts:            []int{1, 2, 3, 4, 5},
			iDequeueMultiple: 5,
			oLength:          5,
		},
		"greater than in": {
			iSize:            5,
			iInts:            []int{1, 2, 3, 4, 5},
			iDequeueMultiple: 10,
			oLength:          5,
		},
	}
	for cDesc, c := range cases {
		//TODO: add documentation
		q := newQueue(c.iSize)
		for _, item := range c.iInts {
			overflow := q.Enqueue(item)
			assert.False(t, overflow, casef, cDesc)
		}
		items := q.DequeueMultiple(c.iDequeueMultiple)
		if assert.Len(t, items, c.oLength, casef, cDesc) {
			for i, item := range items {
				assert.Equal(t, c.iInts[i], item, casef, cDesc)
			}
		}
		q.Close()
	}
}

func Flush(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
}) {

	cases := map[string]struct {
		iSize     int
		iElements []interface{}
	}{
		"empty": {
			iSize: 1,
		},
		"single": {
			iSize:     5,
			iElements: []interface{}{1},
		},
		"max 5": {
			iSize:     5,
			iElements: []interface{}{1, 2, 3, 4, 5},
		},
		"max 1": {
			iSize:     1,
			iElements: []interface{}{1},
		},
	}
	for cDesc, c := range cases {
		//TODO: add documentation
		q := newQueue(c.iSize)
		remainingElements, overflow := q.EnqueueMultiple(c.iElements)
		//KIM: this is testing dequeue, so this should always succeed
		assert.False(t, overflow, casef, cDesc)
		assert.Len(t, remainingElements, 0, casef, cDesc)
		items := q.Flush()
		if assert.Equal(t, len(items), len(c.iElements), casef, cDesc) {
			for i, item := range items {
				assert.Equal(t, c.iElements[i], item, casef, cDesc)
			}
		}
		q.Close()
	}
}

func Peek(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
	goqueue.Info
	goqueue.Peeker
}) {

	cases := map[string]struct {
		iSize      int
		iInts      []int
		oUnderflow []bool
	}{
		"No Elements": {
			iSize:      5,
			iInts:      []int{},
			oUnderflow: []bool{true},
		},
		"Single Element": {
			iSize:      5,
			iInts:      []int{1},
			oUnderflow: []bool{false},
		},
		"Multiple Elements": {
			iSize:      5,
			iInts:      []int{1, 2, 3, 4, 5},
			oUnderflow: []bool{false},
		},
		"Peek Greater than Size": {
			iSize:      5,
			iInts:      []int{1, 2, 3, 4, 5},
			oUnderflow: []bool{false},
		},
	}
	for cDesc, c := range cases {
		//TODO: add documentation
		q := newQueue(c.iSize)
		for _, input := range c.iInts {
			overflow := q.Enqueue(input)
			assert.False(t, overflow, casef, cDesc)
		}
		peeked := q.Peek()
		if assert.Equal(t, len(peeked), len(c.iInts), casef, cDesc) {
			for i, peek := range peeked {
				assert.Equal(t, c.iInts[i], peek, casef, cDesc)
			}
		}
		item, underflow := q.PeekHead()
		if assert.Equal(t, c.oUnderflow[0], underflow) {
			if !underflow {
				assert.Equal(t, c.iInts[0], item, casef, cDesc)
			}
		}
		q.Close()
	}
}

func PeekFromHead(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
	goqueue.Info
	goqueue.Peeker
}) {

	cases := map[string]struct {
		iSize         int
		iElements     []interface{}
		iPeekFromHead int
		oPeeked       []interface{}
	}{
		"min": {
			iSize:         5,
			iElements:     []interface{}{1, 2, 3, 4, 5},
			iPeekFromHead: 1,
			oPeeked:       []interface{}{1},
		},
		"max": {
			iSize:         5,
			iElements:     []interface{}{1, 2, 3, 4, 5},
			iPeekFromHead: 5,
			oPeeked:       []interface{}{1, 2, 3, 4, 5},
		},
		"zero": {
			iSize:         5,
			iElements:     []interface{}{1, 2, 3, 4, 5},
			iPeekFromHead: 0,
		},
		"max+1": {
			iSize:         5,
			iElements:     []interface{}{1, 2, 3, 4, 5},
			iPeekFromHead: 6,
			oPeeked:       []interface{}{1, 2, 3, 4, 5},
		},
		"empty queue": {
			iSize:         5,
			iPeekFromHead: 1,
		},
	}
	for cDesc, c := range cases {
		//TODO: add documentation
		q := newQueue(c.iSize)
		for _, input := range c.iElements {
			overflow := q.Enqueue(input)
			assert.False(t, overflow, casef, cDesc)
		}
		items := q.PeekFromHead(c.iPeekFromHead)
		if assert.Equal(t, len(items), len(c.oPeeked), casef, cDesc) {
			for i, item := range items {
				assert.Equal(t, c.oPeeked[i], item, casef, cDesc)
			}
		}
		q.Close()
	}
}

func Event(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
	goqueue.Event
}) {

	cases := map[string]struct {
		iSize int
	}{
		"normal": {
			iSize: 1,
		},
	}
	for cDesc, c := range cases {
		var wg sync.WaitGroup

		//TODO: add documentation
		q := newQueue(c.iSize)
		signalIn, signalOut := q.GetSignalIn(), q.GetSignalOut()
		wg.Add(2)
		go func() {
			defer wg.Done()
			overflow := q.Enqueue(struct{}{})
			assert.False(t, overflow, casef, cDesc)
		}()
		go func() {
			defer wg.Done()

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
		wg.Wait()
		wg.Add(2)
		go func() {
			defer wg.Done()
			_, underflow := q.Dequeue()
			assert.False(t, underflow, casef, cDesc)
		}()
		go func() {
			defer wg.Done()
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
		wg.Wait()
		q.Close()
	}
}

func Info(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
	goqueue.Info
}) {

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

func Queue(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
	goqueue.Info
}) {

	randFloats := genFloat64(0)
	singleFloat := genFloat64(1)
	cases := map[string]struct {
		Size   int
		Floats []float64
	}{
		"random length": {
			Size:   len(randFloats),
			Floats: randFloats,
		},
		"negative length": {
			Size:   -1,
			Floats: singleFloat,
		},
		"zero length": {
			Size:   0,
			Floats: singleFloat,
		},
	}
	for cDesc, c := range cases {
		// 1. Use the New() function to create/populate a queue of the size for the case
		// 2. Use the Length() function to verify that the queue is empty (size of 0)
		// 3. Use the Enqueue() function for the number of itemsIn to place data in the queue and verify that
		//  the length increases by one each time.
		// 4. Use the Length() to check to see if the queue is the expected size
		// 5. Use the Dequeue() Function again to verify an underflow as the queue should now be empty (length of 0)
		// 6. Use the Close() function to clean up all internal pointers for the queue
		q := newQueue(c.Size)
		assert.Equal(t, 0, q.Length(), casef, cDesc)
		for _, item := range c.Floats {
			overflow := q.Enqueue(item)
			assert.False(t, overflow, casef, cDesc)
		}
		assert.Equal(t, len(c.Floats), q.Length(), casef, cDesc)
		for i := 0; i < len(c.Floats); i++ {
			item, underflow := q.Dequeue()
			if assert.False(t, underflow, casef, cDesc) {
				value, _ := item.(float64)
				assert.Equal(t, c.Floats[i], value, casef, cDesc)
			}
		}
		_, underflow := q.Dequeue()
		assert.True(t, underflow, casef, cDesc)
		q.Close()
	}
}

func Async(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
	goqueue.Info
}) {

	randFloats := genFloat64(0)
	cases := map[string]struct {
		Size   int
		Floats []float64
		Ints   []int
	}{
		"basic": {
			Size:   len(randFloats),
			Floats: randFloats,
		},
	}
	for cDesc, c := range cases {
		var valuesEnqueued, valuesDequeued []float64
		var wg sync.WaitGroup

		//1. Populate async interface using New()
		//2. Create two goRoutines:
		//a. goRoutine (dequeue):
		//(2) Stop when signal received after enqueue function is finished enqueing data
		//(1) Constantly attempt to dequeue, when underflow is false, add item to slice of float64
		//b. goRoutine (enqueue):
		//(1) Enqueue all the data from randFloats, store data in queue
		//(2) Send signal when finished enqueuing data
		//5. Compare the items dequeued to the items enqueued, they should be equal although their quantity may not be the same (see verification)

		q := newQueue(c.Size)
		stopDequeue := make(chan (struct{}))
		wg.Add(1)
		go func() { //dequeue
			defer wg.Done()

			for {
				select {
				case <-stopDequeue:
					return
				default:
					if item, underflow := q.Dequeue(); !underflow {
						value, _ := item.(float64)
						valuesDequeued = append(valuesDequeued, value)
					}
				}
			}
		}()
		wg.Add(1)
		go func() { //enqueue
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
