package infinite_test

import (
	"math/rand"
	"testing"
	"time"

	goqueue "github.com/antonio-alexander/go-queue"
	infinite "github.com/antonio-alexander/go-queue/infinite"
	infinite_tests "github.com/antonio-alexander/go-queue/infinite/tests"
	goqueue_tests "github.com/antonio-alexander/go-queue/tests"
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

func TestInfiniteQueue(t *testing.T) {
	t.Run("Test New", infinite_tests.TestNew(t, func(size int) interface {
		goqueue.Owner
		goqueue.Length
	} {
		return infinite.New(size)
	}))
	t.Run("Test Garbage Collect", goqueue_tests.TestGarbageCollect(t, func(size int) interface {
		goqueue.Owner
		goqueue.GarbageCollecter
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return infinite.New(size)
	}))
	t.Run("Test Garbage Collect Infinite", infinite_tests.TestGarbageCollect(t, func(size int) interface {
		goqueue.Owner
		goqueue.GarbageCollecter
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return infinite.New(size)
	}))
	t.Run("Test Dequeue", goqueue_tests.TestDequeue(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Length
	} {
		return infinite.New(size)
	}))
	t.Run("Test Dequeue Multiple", goqueue_tests.TestDequeueMultiple(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return infinite.New(size)
	}))
	t.Run("Test Flush", goqueue_tests.TestFlush(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return infinite.New(size)
	}))
	t.Run("Test Enqueue", infinite_tests.TestEnqueue(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Length
	} {
		return infinite.New(size)
	}))
	t.Run("Enqueue Multiple", infinite_tests.TestEnqueueMultiple(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Length
	} {
		return infinite.New(size)
	}))
	t.Run("Test Enqueue In Front", infinite_tests.TestEnqueueInFront(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.EnqueueInFronter
		goqueue.Peeker
	} {
		return infinite.New(size)
	}))
	t.Run("Test Peek", goqueue_tests.TestPeek(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Peeker
	} {
		return infinite.New(size)
	}))
	t.Run("Test Peek From Head", goqueue_tests.TestPeekFromHead(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Peeker
	} {
		return infinite.New(size)
	}))
	//configure the timeout to something since the default is 0 and this
	// test would otherwise fail because the signal channels aren't buffered
	infinite.ConfigSignalTimeout = 1 * time.Millisecond
	t.Run("Test Event", goqueue_tests.TestEvent(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Event
	} {
		return infinite.New(size)
	}))
	t.Run("Test Length", goqueue_tests.TestLength(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Length
	} {
		return infinite.New(size)
	}))
	t.Run("Test Queue", goqueue_tests.TestQueue(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Length
	} {
		return infinite.New(size)
	}))
	t.Run("Test Queue Infinite", infinite_tests.TestQueue(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Length
	} {
		return infinite.New(size)
	}))
	t.Run("Test Asynchronous", goqueue_tests.TestAsync(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Length
	} {
		return infinite.New(size)
	}))
}
