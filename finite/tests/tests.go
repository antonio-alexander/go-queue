package finite_tests

import (
	"context"
	"math/rand"
	"testing"
	"time"

	goqueue "github.com/antonio-alexander/go-queue"
	finite "github.com/antonio-alexander/go-queue/finite"

	"github.com/stretchr/testify/assert"
)

const casef string = "case: %s"

func TestEnqueue(t *testing.T, rate, timeout time.Duration, newQueue func(int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
}) func(*testing.T) {
	return func(t *testing.T) {
		//generate size and examples
		size := int(100 * rand.Float64())
		examples := goqueue.ExampleGenFloat64(size)

		//create queue
		q := newQueue(size)
		defer q.Close()

		//enqueue all examples
		for _, example := range examples {
			ctx, cancel := context.WithTimeout(context.TODO(), timeout)
			defer cancel()
			overflow := goqueue.MustEnqueue(q, example, ctx.Done(), rate)
			assert.False(t, overflow)
			cancel()
		}

		//attempt to enqueue (validate overflow)
		example := &goqueue.Example{Int: rand.Int()}
		ctx, cancel := context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		overflow := goqueue.MustEnqueue(q, example, ctx.Done(), rate)
		assert.True(t, overflow)

		//attempt to enqueue again (validate overflow)
		example = &goqueue.Example{Int: rand.Int()}
		ctx, cancel = context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		overflow = goqueue.MustEnqueue(q, example, ctx.Done(), rate)
		assert.True(t, overflow)

		//dequeue once to ensure the queue isn't full
		ctx, cancel = context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		item, underflow := goqueue.MustDequeue(q, ctx.Done(), rate)
		assert.False(t, underflow)
		assert.NotNil(t, item)
		assert.IsType(t, &goqueue.Example{}, item)

		//enqueue an item to validate that the queue wasn't full
		example = &goqueue.Example{Int: rand.Int()}
		ctx, cancel = context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		overflow = goqueue.MustEnqueue(q, example, ctx.Done(), rate)
		assert.False(t, overflow)

		//close the queue
		q.Close()
	}
}

//TestEnqueueEvent will confirm the behavior of the signal channel for enqueues, it will
// operate off of the idea that a fixed queue can utilize a buffered channel such that
// a signal is received for every item enqueued
func TestEnqueueEvent(t *testing.T, rate, timeout time.Duration, newQueue func(size int) interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
	goqueue.Event
}) func(*testing.T) {
	return func(t *testing.T) {
		//generate examples and size
		size := int(100 * rand.Float64())
		examples := goqueue.ExampleGenFloat64(size)

		//create queue
		q := newQueue(size)
		defer q.Close()
		signalIn := q.GetSignalIn()

		//enqueue and verify signals, we're going to test
		// a singal per enqueue
		for _, example := range examples {
			ctx, cancel := context.WithTimeout(context.TODO(), timeout)
			defer cancel()
			overflow := goqueue.MustEnqueue(q, example, ctx.Done(), rate)
			assert.False(t, overflow)
			cancel()
			select {
			default:
				assert.Fail(t, "expected signal not received")
				continue
			case <-signalIn:
				//KIM: this tests makes an assumption that the signal is
				// synchronously sent after a successful enqueue
			}
		}

		//flush queue
		ctx, cancel := context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		goqueue.MustFlush(q, ctx.Done(), rate)
		cancel()

		//attempt to enqueue multiple items and ensure you get a signal for
		// every item enqueued
		items := make([]interface{}, 0, len(examples))
		for _, example := range examples {
			items = append(items, example)
		}
		ctx, cancel = context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		items, overflow := goqueue.MustEnqueueMultiple(q, items, ctx.Done(), rate)
		assert.False(t, overflow)
		assert.Empty(t, items)
		cancel()
		for i := 0; i < len(items); i++ {
			select {
			default:
				assert.Fail(t, "expected signal in not received")
			case <-signalIn:
				//KIM: this assumes that at least one signal is received
			}
		}

		//flush queue
		ctx, cancel = context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		goqueue.MustFlush(q, ctx.Done(), rate)
		cancel()

		//enqueue and validate signal received for every enqueue
		for _, example := range examples {
			ctx, cancel := context.WithTimeout(context.TODO(), timeout)
			defer cancel()
			overflow := goqueue.MustEnqueue(q, example, ctx.Done(), rate)
			assert.False(t, overflow)
			cancel()
		}
		for i := 0; i < len(items); i++ {
			select {
			default:
				assert.Fail(t, "expected signal in not received")
			case <-signalIn:
				//KIM: this assumes that at least one signal is received
			}
		}

		//close the queue
		q.Close()
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

//TODO: test enqueue in front with event

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

//TODO: test enqueue lossy with event

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
	// goqueue.Length
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
			q := newQueue(c.iSize)
			capacity := q.Capacity()
			assert.Equal(t, c.iSize, capacity, casef, cDesc)
			overflow := q.Enqueue(&goqueue.Example{})
			assert.False(t, overflow, casef, cDesc)
			// assert.Equal(t, 1, q.Length(), casef, cDesc)
			_, underflow := q.Dequeue()
			assert.False(t, underflow, casef, cDesc)
			// assert.Equal(t, 0, q.Length(), casef, cDesc)
			for i := 0; i < c.iSize; i++ {
				overflow := q.Enqueue(&goqueue.Example{})
				assert.False(t, overflow, casef, cDesc)
			}
			// assert.Equal(t, c.iSize, q.Length(), casef, cDesc)
			for i := 0; i < c.iSize; i++ {
				_, underflow := q.Dequeue()
				assert.False(t, underflow, casef, cDesc)
			}
			// assert.Equal(t, 0, q.Length(), casef, cDesc)
			q.Close()
		}
	}
}
