package goqueue_tests

import (
	"context"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"

	goqueue "github.com/antonio-alexander/go-queue"

	"github.com/stretchr/testify/assert"
)

const casef string = "case: %s"

func getHeap() (allocated, totalAllocated uint64) {
	//REFERENCE: https://golangcode.com/print-the-current-memory-usage/
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	return m.Alloc, m.TotalAlloc
}

func randomString(nLetters ...int) string {
	//REFERENCE: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	nLetter := 20
	if len(nLetters) > 0 {
		nLetter = nLetters[0]
	}
	b := make([]rune, nLetter)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

//TestGarbageCollect attempts to validate that memory is returned to the heap and can be properly
// garbage collected it executes and garbage collection is manually triggered. This makes some
// heavy underlying assumptions that a slice or map data structure. The idea behind the function
// is that it'll re-create those data structures with new maps/slices such that that the items
// can be de-allocated. This test attempts to put a significant amount of data on the heap to
// validate it's functionality
func TestGarbageCollect(t *testing.T, rate, timeout time.Duration, newQueue func(size int) interface {
	goqueue.Enqueuer
	goqueue.Dequeuer
	goqueue.Owner
	goqueue.GarbageCollecter
}) func(*testing.T) {
	return func(t *testing.T) {
		const kiloByte = 1024 - 13 //number of bytes when example is JSON
		const nElements = 1024

		//create queue
		q := newQueue(1024)
		defer q.Close()

		//generate examples and establish current heap
		allocatedBeforeEnqueue, _ := getHeap()
		for i := 0; i < nElements; i++ {
			example := goqueue.Example{String: randomString(kiloByte)}
			ctx, cancel := context.WithTimeout(context.TODO(), timeout)
			defer cancel()
			overflow := goqueue.MustEnqueue(q, example, ctx.Done(), rate)
			assert.False(t, overflow)
			cancel()
		}

		//validate that more memory has been allocated
		allocatedAfterEnqueue, _ := getHeap()
		assert.Greater(t, allocatedAfterEnqueue, allocatedBeforeEnqueue)

		//flush the queue
		ctx, cancel := context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		items := goqueue.MustFlush(q, ctx.Done(), rate)
		cancel()
		assert.Len(t, items, 1024)

		//garbage collect
		q.GarbageCollect()
		runtime.GC()

		//validate that the allocated memory has gone down
		allocatedAfterGarbageCollect, _ := getHeap()
		assert.Less(t, allocatedAfterGarbageCollect, allocatedAfterEnqueue)

		//close queue
		q.Close()
	}
}

//TestDequeue will confirm the functionality of the underflow output, with an infinite queue
// the expectation is that the enqueue will never overflow, and the dequeue will only underflow
// if the queue is empty. Although this use case is the "same" for infinite and finite queues
// the dequeue function is based on the behavior of the enqueue function. This test has some
// "configuration" items that can be "tweaked" for your specific queue implementation:
//  - rate: this is the rate at which the test will attempt to enqueue/dequeue
//  - timeout: this is when the test will "give up"
// Some assumptions this test does make:
//  - your queue can handle valid data, as a plus the example data type supports the
//    BinaryMarshaller
//  - your queue maintains order
//  - it's safe to use a single instance of your queue for each test case
// Some assumptions this test won't make:
//  - the "size" of the queue affects the behavior of enqueue
func TestDequeue(t *testing.T, rate, timeout time.Duration, newQueue func(size int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
}) func(*testing.T) {
	return func(t *testing.T) {
		//generate examples
		examples := goqueue.ExampleGenFloat64()
		items := make([]interface{}, 0, len(examples))
		for _, example := range examples {
			items = append(items, example)
		}

		//create the queue
		q := newQueue(len(examples))
		defer q.Close()

		//attempt to dequeue (confirm underflow)
		ctx, cancel := context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		//KIM: we don't confirm that the item is nil because it's
		// inconsequential to this test
		_, underflow := goqueue.MustDequeue(q, ctx.Done(), rate)
		cancel()
		assert.True(t, underflow)

		//enqueue items
		ctx, cancel = context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		_, overflow := goqueue.MustEnqueueMultiple(q, items, ctx.Done(), rate)
		assert.False(t, overflow)

		//dequeue set number of items (confirm underflow when empty)
		for i := 0; i < len(examples); i++ {
			ctx, cancel = context.WithTimeout(context.TODO(), timeout)
			defer cancel()
			item, underflow := goqueue.MustDequeue(q, ctx.Done(), rate)
			cancel()
			assert.False(t, underflow)
			example := goqueue.ExampleConvertSingle(item)
			assert.Equal(t, examples[i], example)
		}

		//attempt to dequeue again to confirm that the queue is empty
		ctx, cancel = context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		//KIM: we don't confirm that the item is nil because it's
		// inconsequential to this test
		_, underflow = goqueue.MustDequeue(q, ctx.Done(), rate)
		cancel()
		assert.True(t, underflow)

		//close the queue
		q.Close()
	}
}

func TestDequeueMultiple(t *testing.T, rate, timeout time.Duration, newQueue func(size int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
}) func(*testing.T) {
	return func(t *testing.T) {
		//generate examples
		examples := goqueue.ExampleGenFloat64()

		//create the queue
		q := newQueue(len(examples))
		defer q.Close()

		//attempt to dequeue (confirm underflow)
		ctx, cancel := context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		//KIM: we don't confirm that the item is nil because it's
		// inconsequential to this test
		items := goqueue.MustDequeueMultiple(q, ctx.Done(), len(examples), rate)
		cancel()
		assert.Empty(t, items)

		//enqueue items
		items = make([]interface{}, 0, len(examples))
		for _, example := range examples {
			items = append(items, example)
		}
		ctx, cancel = context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		_, overflow := goqueue.MustEnqueueMultiple(q, items, ctx.Done(), rate)
		assert.False(t, overflow)

		//dequeue set number of items (confirm underflow when empty)
		ctx, cancel = context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		items = goqueue.MustDequeueMultiple(q, ctx.Done(), len(examples), rate)
		assert.Equal(t, len(examples), len(items))
		for i, item := range items {
			example := goqueue.ExampleConvertSingle(item)
			assert.Equal(t, examples[i], example)
		}

		//attempt to dequeue again to confirm that the queue is empty
		ctx, cancel = context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		//KIM: we don't confirm that the item is nil because it's
		// inconsequential to this test
		_, underflow := goqueue.MustDequeue(q, ctx.Done(), rate)
		cancel()
		assert.True(t, underflow)

		//close the queue
		q.Close()
	}
}

func TestFlush(t *testing.T, rate, timeout time.Duration, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
}) func(*testing.T) {
	return func(t *testing.T) {
		//generate examples
		examples := goqueue.ExampleGenFloat64()

		//create the queue
		q := newQueue(len(examples))
		defer q.Close()

		//attempt to dequeue (confirm underflow)
		ctx, cancel := context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		//KIM: we don't confirm that the item is nil because it's
		// inconsequential to this test
		items := goqueue.MustFlush(q, ctx.Done(), rate)
		cancel()
		assert.Empty(t, items)

		//enqueue items
		items = make([]interface{}, 0, len(examples))
		for _, example := range examples {
			items = append(items, example)
		}
		ctx, cancel = context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		_, overflow := goqueue.MustEnqueueMultiple(q, items, ctx.Done(), rate)
		assert.False(t, overflow)

		//dequeue set number of items (confirm underflow when empty)
		ctx, cancel = context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		items = goqueue.MustFlush(q, ctx.Done(), rate)
		assert.Equal(t, len(examples), len(items))
		for i, item := range items {
			example := goqueue.ExampleConvertSingle(item)
			assert.Equal(t, examples[i], example)
		}

		//attempt to dequeue again to confirm that the queue is empty
		ctx, cancel = context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		//KIM: we don't confirm that the item is nil because it's
		// inconsequential to this test
		items = goqueue.MustFlush(q, ctx.Done(), rate)
		cancel()
		assert.Empty(t, items)

		//close the queue
		q.Close()
	}
}

func TestDequeueEvent(t *testing.T, rate, timeout time.Duration, newQueue func(size int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
	goqueue.Event
}) func(*testing.T) {
	return func(t *testing.T) {
		//generate examples
		examples := goqueue.ExampleGenFloat64()

		//create queue
		q := newQueue(len(examples))
		defer q.Close()
		signalOut := q.GetSignalOut()

		//enqueue multiple items (to be dequeued)
		items := make([]interface{}, 0, len(examples))
		for _, example := range examples {
			items = append(items, example)
		}
		ctx, cancel := context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		items, overflow := goqueue.MustEnqueueMultiple(q, items, ctx.Done(), rate)
		cancel()
		assert.False(t, overflow)
		assert.Empty(t, items)

		//attempt to dequeue each item and check if signals received
		for range examples {
			ctx, cancel := context.WithTimeout(context.TODO(), timeout)
			defer cancel()
			item, underflow := goqueue.MustDequeue(q, ctx.Done(), rate)
			assert.False(t, underflow)
			assert.NotNil(t, item)
			select {
			default:
				assert.Fail(t, "expected signal out not received")
			case <-signalOut:
				//KIM: this assumes that at least one signal is received
			}
		}

		//enqueue multiple items (to be dequeued)
		items = make([]interface{}, 0, len(examples))
		for _, example := range examples {
			items = append(items, example)
		}
		ctx, cancel = context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		items, overflow = goqueue.MustEnqueueMultiple(q, items, ctx.Done(), rate)
		cancel()
		assert.False(t, overflow)
		assert.Empty(t, items)

		//dequeue multiple items
		ctx, cancel = context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		n := len(examples) - 3
		items = goqueue.MustDequeueMultiple(q, ctx.Done(), n, rate)
		cancel()
		assert.Equal(t, n, len(items))
		select {
		default:
			assert.Fail(t, "expected signal out not received")
		case <-signalOut:
			//KIM: this assumes that at least one signal is received
		}

		//flush items
		ctx, cancel = context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		items = goqueue.MustFlush(q, ctx.Done(), rate)
		cancel()
		assert.Equal(t, 3, len(items))
		select {
		default:
			assert.Fail(t, "expected signal out not received")
		case <-signalOut:
			//KIM: this assumes that at least one signal is received
		}

		//close the queue
		q.Close()
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
func TestQueue(t *testing.T, rate, timeout time.Duration, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
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
			for _, item := range c.iFloats {
				overflow := q.Enqueue(item)
				assert.False(t, overflow, casef, cDesc)
			}
			for i := 0; i < len(c.iFloats); i++ {
				ctx, cancel := context.WithTimeout(context.TODO(), timeout)
				defer cancel()
				value, underflow := goqueue.MustDequeue(q, ctx.Done(), rate)
				cancel()
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

//REVIEW: implement tests for sanity/security checks
// * When using dequeue methods that output slices, can we ensure we don't accidentally leak the
//   underlying slice?
