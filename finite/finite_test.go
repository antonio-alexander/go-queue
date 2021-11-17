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

func TestNew(t *testing.T) {
	goqueue_tests.New(t, func(size int) interface {
		goqueue.Owner
		goqueue.Info
	} {
		return finite.New(size)
	})
}

func TestGarbageCollect(t *testing.T) {
	goqueue_tests.GarbageCollect(t, func(size int) interface {
		goqueue.Owner
		goqueue.GarbageCollecter
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Info
	} {
		return finite.New(size)
	})
}

func TestDequeue(t *testing.T) {
	goqueue_tests.Dequeue(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Info
	} {
		return finite.New(size)
	})
}

func TestDequeueMultiple(t *testing.T) {
	goqueue_tests.DequeueMultiple(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return finite.New(size)
	})
}

func TestFlush(t *testing.T) {
	goqueue_tests.Flush(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
	} {
		return finite.New(size)
	})
}

func TestEnqueue(t *testing.T) {
	finite_tests.Enqueue(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
	} {
		return finite.New(size)
	})
}
func TestEnqueueMultiple(t *testing.T) {
	finite_tests.EnqueueMultiple(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
	} {
		return finite.New(size)
	})
}

func TestResize(t *testing.T) {
	finite_tests.Resize(t, func(size int) interface {
		finite.Resizer
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Info
	} {
		return finite.New(size)
	})
}

func TestEnqueueLossy(t *testing.T) {
	finite_tests.EnqueueLossy(t, func(size int) interface {
		goqueue.Owner
		finite.EnqueueLossy
	} {
		return finite.New(size)
	})
}

func TestEnqueueInFront(t *testing.T) {
	finite_tests.EnqueueInFront(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.EnqueueInFronter
		goqueue.Peeker
	} {
		return finite.New(size)
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
		return finite.New(size)
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
		return finite.New(size)
	})
}

func TestEvent(t *testing.T) {
	goqueue_tests.Event(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Event
	} {
		return finite.New(size)
	})
}

func TestInfo(t *testing.T) {
	goqueue_tests.Info(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Info
	} {
		return finite.New(size)
	})
}

func TestQueue(t *testing.T) {
	goqueue_tests.Queue(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Info
	} {
		return finite.New(size)
	})
}

func TestAsync(t *testing.T) {
	goqueue_tests.Async(t, func(size int) interface {
		goqueue.Owner
		goqueue.Enqueuer
		goqueue.Dequeuer
		goqueue.Info
	} {
		return finite.New(size)
	})
}
