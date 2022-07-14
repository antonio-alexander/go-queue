package finite_test

import (
	"math/rand"
	"testing"
	"time"

	goqueue "github.com/antonio-alexander/go-queue"
	finite "github.com/antonio-alexander/go-queue/finite"
	finite_tests "github.com/antonio-alexander/go-queue/finite/tests"
	goqueue_tests "github.com/antonio-alexander/go-queue/tests"
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

func TestFiniteQueue(t *testing.T) {
	t.Run("Test New", finite_tests.TestNew(t, func(size int) interface {
		goqueue.Owner
		finite.Capacity
	} {
		return finite.New(size)
	}))
	t.Run("Test Garbage Collect", goqueue_tests.TestGarbageCollect(t, func(size int) interface {
		goqueue.Owner
		goqueue.GarbageCollecter
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return finite.New(size)
	}))
	t.Run("Test Dequeue", goqueue_tests.TestDequeue(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Length
	} {
		return finite.New(size)
	}))
	t.Run("Test Dequeue Multiple", goqueue_tests.TestDequeueMultiple(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return finite.New(size)
	}))
	t.Run("Test Flush", goqueue_tests.TestFlush(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return finite.New(size)
	}))
	t.Run("Test Enqueue", finite_tests.TestEnqueue(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
	} {
		return finite.New(size)
	}))
	t.Run("Test Enqueue Multiple", finite_tests.TestEnqueueMultiple(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
	} {
		return finite.New(size)
	}))
	t.Run("Test Resize", finite_tests.TestResize(t, func(size int) interface {
		finite.Capacity
		goqueue.Enqueuer
		goqueue.Owner
		finite.Resizer
	} {
		return finite.New(size)
	}))
	t.Run("Test Enqueue Lossy", finite_tests.TestEnqueueLossy(t, func(size int) interface {
		goqueue.Owner
		finite.EnqueueLossy
	} {
		return finite.New(size)
	}))
	t.Run("Test Enqueue In Front", finite_tests.TestEnqueueInFront(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.EnqueueInFronter
		goqueue.Peeker
	} {
		return finite.New(size)
	}))
	t.Run("Test Peek", goqueue_tests.TestPeek(t, func(size int) interface {
		goqueue.Dequeuer
		goqueue.Enqueuer
		goqueue.Owner
		goqueue.Peeker
	} {
		return finite.New(size)
	}))
	t.Run("Test Peek From Head", goqueue_tests.TestPeekFromHead(t, func(size int) interface {
		goqueue.Dequeuer
		goqueue.Enqueuer
		goqueue.Owner
		goqueue.Peeker
	} {
		return finite.New(size)
	}))
	t.Run("Test Event", goqueue_tests.TestEvent(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Event
	} {
		return finite.New(size)
	}))
	t.Run("Test Length", goqueue_tests.TestLength(t, func(size int) interface {
		goqueue.Dequeuer
		goqueue.Enqueuer
		goqueue.Length
		goqueue.Owner
	} {
		return finite.New(size)
	}))
	t.Run("Test Capacity", finite_tests.TestCapacity(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Length
		finite.Capacity
	} {
		return finite.New(size)
	}))
	t.Run("Test Queue", goqueue_tests.TestQueue(t, func(size int) interface {
		goqueue.Dequeuer
		goqueue.Enqueuer
		goqueue.Length
		goqueue.Owner
	} {
		return finite.New(size)
	}))
	t.Run("Test Asynchronous", goqueue_tests.TestAsync(t, func(size int) interface {
		goqueue.Dequeuer
		goqueue.Enqueuer
		goqueue.Length
		goqueue.Owner
	} {
		return finite.New(size)
	}))
}
