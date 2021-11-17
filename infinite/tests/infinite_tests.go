package infinite_tests

import (
	"math/rand"
	"testing"
	"time"

	goqueue "github.com/antonio-alexander/go-queue"

	"github.com/stretchr/testify/assert"
)

const (
	casef          string = "case: %s"
	testCaseMap    string = ", case: \"%s\""
	testUnexpected string = "%s, unexpected %s"
)

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
		n := rand.Float64()

		floats = append(floats, n)
	}

	return
}

func GarbageCollect(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.GarbageCollecter
	goqueue.Enqueuer
	goqueue.Dequeuer
	goqueue.Info
}) {

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
		capacity := q.Capacity()
		assert.Equal(t, c.oCapacity, capacity, casef, cDesc)
		for i := 0 + c.iNtoDequeue; i < c.iNtoEnqueue; i++ {
			element, underflow := q.Dequeue()
			if assert.False(t, underflow, casef, cDesc) {
				assert.Equal(t, i, element)
			}
		}
		q.Close()
	}
}

func Enqueue(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Info
}) {

	cases := map[string]struct {
		iGrowSize   int
		iElements   []interface{}
		oCapacities []int
		oLengths    []int
	}{
		"No Grow": {
			iGrowSize:   10,
			iElements:   []interface{}{1, 2, 3},
			oCapacities: []int{10, 10, 10},
			oLengths:    []int{1, 2, 3},
		},
		"grow once": {
			iGrowSize:   2,
			iElements:   []interface{}{1, 2, 3},
			oCapacities: []int{2, 2, 4},
			oLengths:    []int{1, 2, 3},
		},
		"grow twice": {
			iGrowSize:   2,
			iElements:   []interface{}{1, 2, 3, 4, 5, 6},
			oCapacities: []int{2, 2, 4, 4, 6, 6},
			oLengths:    []int{1, 2, 3, 4, 5, 6},
		},
	}
	for cDesc, c := range cases {
		q := newQueue(c.iGrowSize)
		for i, element := range c.iElements {
			overflow := q.Enqueue(element)
			assert.False(t, overflow, casef, cDesc)
			capacity := q.Capacity()
			assert.Equal(t, c.oCapacities[i], capacity, casef, cDesc)
			length := q.Length()
			assert.Equal(t, c.oLengths[i], length)
		}
		q.Close()
	}
}

func EnqueueMultiple(t *testing.T, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Info
}) {

	cases := map[string]struct {
		iGrowSize int
		iElements []interface{}
		oCapacity int
		oLength   int
	}{
		"No Grow": {
			iGrowSize: 10,
			iElements: []interface{}{1, 2, 3},
			oLength:   3,
			oCapacity: 10,
		},
		"grow once": {
			iGrowSize: 2,
			iElements: []interface{}{1, 2, 3},
			oLength:   3,
			oCapacity: 4,
		},
		"grow twice": {
			iGrowSize: 2,
			iElements: []interface{}{1, 2, 3, 4, 5, 6},
			oLength:   6,
			oCapacity: 6,
		},
	}
	for cDesc, c := range cases {
		q := newQueue(c.iGrowSize)
		remainingElements, overflow := q.EnqueueMultiple(c.iElements)
		if assert.False(t, overflow, casef, cDesc) &&
			assert.GreaterOrEqual(t, len(remainingElements), 0) {
			capacity := q.Capacity()
			assert.GreaterOrEqual(t, capacity, c.oCapacity, casef, cDesc)
			length := q.Length()
			assert.Equal(t, c.oLength, length, casef, cDesc)
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
		iSize    int
		iInFront int
		iInts    []int
	}{
		"empty_queue": { //verify that eif doesn't overflow if queue is empty
			iSize:    5,
			iInFront: 10,
		},
		"queue_not_full": { //verify that eif doesn't overflow if queue is not full
			iSize:    5,
			iInts:    []int{1, 2, 3, 4},
			iInFront: 10,
		},
		"queue_full": { //verify that eif overflows if queue is full
			iSize:    5,
			iInts:    []int{1, 2, 3, 4, 5},
			iInFront: 88,
		},
		"almost_full": { //verify that eif places data in front if you enqueue one element
			iSize:    5,
			iInts:    []int{1},
			iInFront: 66,
		},
		"max-1": { //verify that eif places data in front if you enqueue max -1 number of elements
			iSize:    5,
			iInts:    []int{1, 2, 3, 4},
			iInFront: 44,
		},
	}
	for cDesc, c := range cases {
		q := newQueue(c.iSize)
		for _, element := range c.iInts {
			overflow := q.Enqueue(element)
			assert.False(t, overflow, casef, cDesc)
		}
		overflow := q.EnqueueInFront(c.iInFront)
		if assert.False(t, overflow, casef, cDesc) {
			element, underflow := q.PeekHead()
			if assert.False(t, underflow, casef, cDesc) {
				assert.Equal(t, c.iInFront, element, casef, cDesc)
			}
		}
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
	//randInts := genInts(0)

	cases := map[string]struct {
		Size   int
		Floats []float64
		Ints   []int
	}{
		"Ticker": { //test with ticker
			Size:   len(randFloats),
			Floats: randFloats,
		},
	}

	for cDesc, c := range cases {
		//1. Use the New() function to create/populate a queue of the size for the case
		q := newQueue(c.Size)
		for i := 0; i < len(c.Floats)-1; i++ {
			//2. Use the Enqueue() function for the number of elementsIn to place data in the queue and verify that the length increases by one each time.
			overflow := q.Enqueue(c.Floats[i])
			assert.False(t, overflow, casef, cDesc)
		}
		assert.Equal(t, q.Length(), len(c.Floats)-1)
		//3. Use the Enqueue() function again to verify that the queue doesn't overflow and that the capacity increases (capacity after enqueue is greater than the capacity before)
		overflow := q.Enqueue(c.Floats[len(c.Floats)-1])
		assert.False(t, overflow, casef, cDesc)
		assert.Equal(t, len(c.Floats), q.Length())
		for i := 0; i < c.Size; i++ {
			//4. Use the Dequeue() function for the number of elementsIn to remove data from the queue adn verify that the length decreases by one each time.
			element, underflow := q.Dequeue()
			if assert.False(t, underflow, casef, cDesc) {
				switch value := element.(type) {
				default:
					assert.Fail(t, testUnexpected+testCaseMap, "type doesn't match", cDesc)
				case int:
					assert.Equal(t, c.Ints[i], value)
				case float64:
					assert.Equal(t, c.Floats[i], value)
				}
			}
		}
		//5. Use the Dequeue() Function again to verify an underflow as the queue should now be empty (length of 0)
		_, underflow := q.Dequeue()
		assert.True(t, underflow, casef, cDesc)
		//6. Use the Close() function to clean up all internal pointers for the queue
		q.Close()
	}
}

//REVIEW: implement tests for sanity/security checks
// * When using dequeue methods that output slices, can we ensure we don't accidentally leak the
//   underlying slice? This should be possible using runtime garbage collection and total/allocated
//   heap memory.
