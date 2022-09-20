package infinite_tests

import (
	"context"
	"testing"
	"time"

	goqueue "github.com/antonio-alexander/go-queue"

	"github.com/stretchr/testify/assert"
)

const casef string = "case: %s"

//TestEnqueue attempt to unit test the enqueue function, in general it confirms the
// behavior, that for an infinite queue, no matter how much data you put into the
// queue, the queue will never overflow and there will be no data loss. This test
// also assumes that the "size" of the queue won't affect the behavior of enqueue
// This test has some "configuration" items that can be "tweaked" for your specific
// queue implementation:
//  - rate: this is the rate at which the test will attempt to enqueue/dequeue
//  - timeout: this is when the test will "give up"
// Some assumptions this test does make:
//  - your queue can handle valid data, as a plus the example data type supports the
//    BinaryMarshaller
// Some assumptions this test won't make:
//  - order is maintained
//  - the "size" of the queue affects the behavior of enqueue
func TestEnqueue(t *testing.T, rate, timeout time.Duration, newQueue func() interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
}) func(*testing.T) {
	return func(t *testing.T) {
		//create the queue
		q := newQueue()
		defer q.Close()

		//generate and enqueue the examples
		examples := goqueue.ExampleGenFloat64()
		for _, example := range examples {
			ctx, cancel := context.WithTimeout(context.TODO(), timeout)
			defer cancel()
			overflow := goqueue.MustEnqueue(q, example, ctx.Done(), rate)
			assert.False(t, overflow)
			cancel()
		}

		//attempt to dequeue the examples and ensure no data loss
		ctx, cancel := context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		items := goqueue.MustFlush(q, ctx.Done(), rate)
		cancel()
		for _, item := range items {
			example := goqueue.ExampleConvertSingle(item)
			if !assert.NotNil(t, example) {
				continue
			}
			assert.Condition(t, goqueue.AssertExamples(example, examples))
		}

		//KIM: we do the same thing a second time, this would catch situations
		// where slices or data structures aren't properly allocated or something
		// explodes
		//generate and enqueue the examples
		examples = goqueue.ExampleGenFloat64()
		for _, example := range examples {
			ctx, cancel := context.WithTimeout(context.TODO(), timeout)
			defer cancel()
			overflow := goqueue.MustEnqueue(q, example, ctx.Done(), rate)
			assert.False(t, overflow)
			cancel()
		}

		//attempt to dequeue the examples and ensure no data loss
		ctx, cancel = context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		items = goqueue.MustFlush(q, ctx.Done(), rate)
		cancel()
		for _, example := range goqueue.ExampleConvertMultiple(items) {
			if !assert.NotNil(t, example) {
				continue
			}
			assert.Condition(t, goqueue.AssertExamples(example, examples))
		}
	}
}

//TestEnqueueMultiple will attempt to unit test the EnqueueMultiple function;
// for an infinite queue, this function will never overflow nor will it return
// items that weren't able to be enqueued.
// This test has some "configuration" items that can be "tweaked" for your specific
// queue implementation:
//  - rate: this is the rate at which the test will attempt to enqueue/dequeue
//  - timeout: this is when the test will "give up"
// Some assumptions this test does make:
//  - your queue can handle valid data, as a plus the example data type supports the
//    BinaryMarshaller
// Some assumptions this test won't make:
//  - order is maintained
//  - the "size" of the queue affects the behavior of enqueue
func TestEnqueueMultiple(t *testing.T, rate, timeout time.Duration, newQueue func() interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Dequeuer
}) func(*testing.T) {
	return func(t *testing.T) {
		//create the queue
		q := newQueue()
		defer q.Close()

		//generate and enqueue the examples
		examples := goqueue.ExampleGenFloat64()
		items := make([]interface{}, 0, len(examples))
		for _, example := range examples {
			items = append(items, example)
		}
		ctx, cancel := context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		items, overflow := goqueue.MustEnqueueMultiple(q, items, ctx.Done(), rate)
		assert.False(t, overflow)
		assert.Empty(t, items)
		cancel()

		//attempt to dequeue the examples and ensure no data loss
		ctx, cancel = context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		items = goqueue.MustFlush(q, ctx.Done(), rate)
		cancel()
		for _, example := range goqueue.ExampleConvertMultiple(items) {
			if !assert.NotNil(t, example) {
				continue
			}
			assert.Condition(t, goqueue.AssertExamples(example, examples))
		}
		//KIM: we do the same thing a second time, this would catch situations
		// where slices or data structures aren't properly allocated or something
		// explodes
		examples = goqueue.ExampleGenFloat64()
		items = make([]interface{}, 0, len(examples))
		for _, example := range examples {
			items = append(items, example)
		}
		ctx, cancel = context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		items, overflow = goqueue.MustEnqueueMultiple(q, items, ctx.Done(), rate)
		assert.False(t, overflow)
		assert.Empty(t, items)
		cancel()

		//attempt to dequeue the examples and ensure no data loss
		ctx, cancel = context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		items = goqueue.MustFlush(q, ctx.Done(), rate)
		cancel()
		for _, example := range goqueue.ExampleConvertMultiple(items) {
			if !assert.NotNil(t, example) {
				continue
			}
			assert.Condition(t, goqueue.AssertExamples(example, examples))
		}
	}
}

//TestEnqueueInFront will validate that if there is data in the queue and you attempt to
// enqueue in front, that special "data" will go to the front, while if the queue is empty
// that data will just be "in" the queue (a regular queue if the queue is empty).
// This test has some "configuration" items that can be "tweaked" for your specific
// queue implementation:
//  - rate: this is the rate at which the test will attempt to enqueue/dequeue
//  - timeout: this is when the test will "give up"
// Some assumptions this test does make:
//  - your queue can handle valid data, as a plus the example data type supports the
//    BinaryMarshaller
//  - your queue maintains order
//  - it's safe to use a single instance of your queue for each test case
// Some assumptions this test won't make:
//  - the "size" of the queue affects the behavior of enqueue
func TestEnqueueInFront(t *testing.T, rate, timeout time.Duration, newQueue func() interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.EnqueueInFronter
	goqueue.Dequeuer
}) func(*testing.T) {
	return func(t *testing.T) {
		//verify that eif doesn't overflow if queue is not full
		//verify that eif doesn't overflow if queue is empty

		cases := map[string]struct {
			iValues       []interface{}
			iInFrontValue *goqueue.Example
		}{
			"empty_queue": {
				iInFrontValue: &goqueue.Example{Int: 10},
			},
			"queue_not_full": {
				iValues: []interface{}{
					&goqueue.Example{Int: 1},
					&goqueue.Example{Int: 2},
					&goqueue.Example{Int: 3},
					&goqueue.Example{Int: 4},
				},
				iInFrontValue: &goqueue.Example{Int: 10},
			},
		}
		q := newQueue()
		defer q.Close()
		for cDesc, c := range cases {
			//enqueue values
			ctx, cancel := context.WithTimeout(context.TODO(), timeout)
			defer cancel()
			items, overflow := goqueue.MustEnqueueMultiple(q, c.iValues, ctx.Done(), rate)
			cancel()
			assert.False(t, overflow, casef, cDesc)
			assert.Empty(t, items, casef, cDesc)

			//attempt to enqueue value in front
			overflow = q.EnqueueInFront(c.iInFrontValue)
			assert.False(t, overflow, casef, cDesc)

			//validate that next dequeued value is value enqueued in front
			ctx, cancel = context.WithTimeout(context.TODO(), timeout)
			defer cancel()
			item, underflow := goqueue.MustDequeue(q, ctx.Done(), rate)
			cancel()
			assert.False(t, underflow)
			assert.IsType(t, &goqueue.Example{}, item, casef, cDesc)
			example, _ := item.(*goqueue.Example)
			assert.Equal(t, c.iInFrontValue, example, casef, cDesc)

			//flush the queue to empty it
			ctx, cancel = context.WithTimeout(context.TODO(), timeout)
			defer cancel()
			goqueue.MustDequeueMultiple(q, ctx.Done(), len(c.iValues), rate)
			cancel()
		}
		q.Close()
	}
}

//TestEnqueueEvent will confirm that the signal channels function correctly when data is enqueued,
// this function for an infinite queue is slightly different because it can't be lossless, there's
// no way to properly implement a buffered channel with an infinite queue
// Some assumptions this test does make:
func TestEnqueueEvent(t *testing.T, rate, timeout time.Duration, newQueue func() interface {
	goqueue.Owner
	goqueue.Enqueuer
	goqueue.Event
}) func(*testing.T) {
	return func(t *testing.T) {
		//create queue
		q := newQueue()
		defer q.Close()
		signalIn := q.GetSignalIn()

		//generate examples, enqueue and verify signals, we're going to test
		// a singal per enqueue
		examples := goqueue.ExampleGenFloat64()
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

		//attempt to enqueue multiple items and see that you get at least one
		// signal
		items := make([]interface{}, 0, len(examples))
		for _, example := range examples {
			items = append(items, example)
		}
		ctx, cancel := context.WithTimeout(context.TODO(), timeout)
		defer cancel()
		items, overflow := goqueue.MustEnqueueMultiple(q, items, ctx.Done(), rate)
		assert.False(t, overflow)
		assert.Empty(t, items)
		cancel()
		select {
		default:
			assert.Fail(t, "expected signal in not received")
		case <-signalIn:
			//KIM: this assumes that at least one signal is received
		}

		//close the queue
		q.Close()
	}
}
