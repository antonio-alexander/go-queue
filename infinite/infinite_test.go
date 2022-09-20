package infinite_test

import (
	"math/rand"
	"testing"
	"time"

	infinite "github.com/antonio-alexander/go-queue/infinite"
	infinite_tests "github.com/antonio-alexander/go-queue/infinite/tests"

	goqueue "github.com/antonio-alexander/go-queue"
	goqueue_tests "github.com/antonio-alexander/go-queue/tests"
)

const (
	queueGrowSize = 1024
	mustTimeout   = time.Second
	mustRate      = time.Millisecond
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

//TODO: test the growing of the internal slice

func TestInfiniteQueue(t *testing.T) {
	t.Run("Test Enqueue", infinite_tests.TestEnqueue(t, mustRate, mustTimeout, func() interface {
		goqueue.Dequeuer
		goqueue.Enqueuer
		goqueue.Owner
	} {
		return infinite.New(queueGrowSize)
	}))
	t.Run("Enqueue Multiple", infinite_tests.TestEnqueueMultiple(t, mustRate, mustTimeout, func() interface {
		goqueue.Dequeuer
		goqueue.Enqueuer
		goqueue.Owner
	} {
		return infinite.New(queueGrowSize)
	}))
	t.Run("Test Enqueue In Front", infinite_tests.TestEnqueueInFront(t, mustRate, mustTimeout, func() interface {
		goqueue.Dequeuer
		goqueue.EnqueueInFronter
		goqueue.Enqueuer
		goqueue.Owner
	} {
		return infinite.New(queueGrowSize)
	}))
}

func TestQueue(t *testing.T) {
	t.Run("Test Garbage Collect", goqueue_tests.TestGarbageCollect(t, mustRate, mustTimeout, func(size int) interface {
		goqueue.Owner
		goqueue.GarbageCollecter
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return infinite.New(size)
	}))
	t.Run("Test Dequeue", goqueue_tests.TestDequeue(t, mustRate, mustTimeout, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return infinite.New(size)
	}))
	t.Run("Test Dequeue Multiple", goqueue_tests.TestDequeueMultiple(t, mustRate, mustTimeout, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return infinite.New(size)
	}))
	t.Run("Test Flush", goqueue_tests.TestFlush(t, mustRate, mustTimeout, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
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
	t.Run("Test Queue", goqueue_tests.TestQueue(t, mustRate, mustTimeout, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return infinite.New(size)
	}))
	t.Run("Test Asynchronous", goqueue_tests.TestAsync(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return infinite.New(size)
	}))
}
