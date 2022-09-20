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

const (
	mustTimeout = time.Second
	mustRate    = time.Millisecond
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

func TestFiniteQueue(t *testing.T) {
	t.Run("Test Enqueue", finite_tests.TestEnqueue(t, mustRate, mustTimeout, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return finite.New(size)
	}))
	t.Run("Test Enqueue Multiple", finite_tests.TestEnqueueMultiple(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
	} {
		return finite.New(size)
	}))
	t.Run("Test Enqueue Event", finite_tests.TestEnqueueEvent(t, mustRate, mustTimeout, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Event
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
	t.Run("Test Capacity", finite_tests.TestCapacity(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		finite.Capacity
	} {
		return finite.New(size)
	}))
}

func TestQueue(t *testing.T) {
	t.Run("Test Dequeue", goqueue_tests.TestDequeue(t, mustRate, mustTimeout, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return finite.New(size)
	}))
	t.Run("Test Dequeue Event", goqueue_tests.TestDequeueEvent(t, mustRate, mustTimeout, func(size int) interface {
		goqueue.Dequeuer
		goqueue.Enqueuer
		goqueue.Event
		goqueue.Owner
	} {
		return finite.New(size)
	}))
	t.Run("Test Dequeue Multiple", goqueue_tests.TestDequeueMultiple(t, mustRate, mustTimeout, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return finite.New(size)
	}))
	t.Run("Test Flush", goqueue_tests.TestFlush(t, mustRate, mustTimeout, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return finite.New(size)
	}))
	t.Run("Test Peek", goqueue_tests.TestPeek(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Peeker
	} {
		return finite.New(size)
	}))
	t.Run("Test Peek From Head", goqueue_tests.TestPeekFromHead(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Peeker
	} {
		return finite.New(size)
	}))
	t.Run("Test Length", goqueue_tests.TestLength(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Length
	} {
		return finite.New(size)
	}))
	t.Run("Test Garbage Collect", goqueue_tests.TestGarbageCollect(t, mustRate, mustTimeout, func(size int) interface {
		goqueue.Owner
		goqueue.GarbageCollecter
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return finite.New(size)
	}))
	//
	t.Run("Test Queue", goqueue_tests.TestQueue(t, mustRate, mustTimeout, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return finite.New(size)
	}))
	t.Run("Test Asynchronous", goqueue_tests.TestAsync(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return finite.New(size)
	}))
}
