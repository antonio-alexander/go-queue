package infinite_test

import (
	"math/rand"
	"testing"
	"time"

	goqueue "github.com/antonio-alexander/go-queue"
	infinite "github.com/antonio-alexander/go-queue/infinite"
	internal "github.com/antonio-alexander/go-queue/internal"

	infinite_tests "github.com/antonio-alexander/go-queue/infinite/tests"
	goqueue_tests "github.com/antonio-alexander/go-queue/tests"
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

func TestNew(t *testing.T) {
	goqueue_tests.New(t, func(size int) interface {
		goqueue.Owner
		goqueue.Info
	} {
		return infinite.New(size)
	})
}

func TestGarbageCollect(t *testing.T) {
	newQueue := func(size int) interface {
		goqueue.Owner
		goqueue.GarbageCollecter
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Info
	} {
		return infinite.New(size)
	}
	goqueue_tests.GarbageCollect(t, newQueue)
	infinite_tests.GarbageCollect(t, newQueue)
}

func TestDequeue(t *testing.T) {
	goqueue_tests.Dequeue(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Info
	} {
		return infinite.New(size)
	})
}

func TestDequeueMultiple(t *testing.T) {
	goqueue_tests.DequeueMultiple(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return infinite.New(size)
	})
}

func TestFlush(t *testing.T) {
	goqueue_tests.Flush(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return infinite.New(size)
	})
}

func TestEnqueue(t *testing.T) {
	infinite_tests.Enqueue(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Info
	} {
		return infinite.New(size)
	})
}

func TestEnqueueMultiple(t *testing.T) {
	infinite_tests.EnqueueMultiple(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Info
	} {
		return infinite.New(size)
	})
}

func TestEnqueueInFront(t *testing.T) {
	infinite_tests.EnqueueInFront(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.EnqueueInFronter
		goqueue.Peeker
	} {
		return infinite.New(size)
	})
}

func TestPeek(t *testing.T) {
	goqueue_tests.Peek(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Info
		goqueue.Peeker
	} {
		return infinite.New(size)
	})
}

func TestPeekFromHead(t *testing.T) {
	goqueue_tests.PeekFromHead(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Info
		goqueue.Peeker
	} {
		return infinite.New(size)
	})
}

func TestEvent(t *testing.T) {
	//configure the timeout to something since the default is 0 and this
	// test would otherwise fail because the signal channels aren't buffered
	internal.ConfigSignalTimeout = 1 * time.Millisecond
	goqueue_tests.Event(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Event
	} {
		return infinite.New(size)
	})
}

func TestInfo(t *testing.T) {
	goqueue_tests.Info(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Info
	} {
		return infinite.New(size)
	})
}

func TestQueue(t *testing.T) {
	goqueue_tests.Queue(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Info
	} {
		return infinite.New(size)
	})
	infinite_tests.Queue(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Info
	} {
		return infinite.New(size)
	})
}

func TestAsync(t *testing.T) {
	goqueue_tests.Async(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Info
	} {
		return infinite.New(size)
	})
}
