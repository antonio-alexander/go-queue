package finite

import (
	"sync"

	goqueue "github.com/antonio-alexander/go-queue"
	internal "github.com/antonio-alexander/go-queue/internal"
)

type queueFinite struct {
	sync.RWMutex
	signalIn  chan struct{}
	signalOut chan struct{}
	data      []any
}

func New(size int) interface {
	goqueue.Owner
	goqueue.GarbageCollecter
	goqueue.Dequeuer
	goqueue.Enqueuer
	goqueue.EnqueueInFronter
	goqueue.Length
	goqueue.Event
	goqueue.Peeker
	EnqueueLossy
	Resizer
	Capacity
} {
	maxSize := max(size, 1)
	return &queueFinite{
		signalIn:  make(chan struct{}, maxSize),
		signalOut: make(chan struct{}, maxSize),
		data:      make([]any, 0, maxSize),
	}
}

func (q *queueFinite) Close() []any {
	q.Lock()
	defer q.Unlock()

	remainingElements, _, _ := internal.DequeueMultiple(cap(q.data), q.data)
	if q.signalIn != nil {
		select {
		default:
			close(q.signalIn)
		case <-q.signalIn:
		}
	}
	if q.signalOut != nil {
		select {
		default:
			close(q.signalOut)
		case <-q.signalOut:
		}
	}
	q.data, q.signalIn, q.signalOut = nil, nil, nil
	return remainingElements
}

func (q *queueFinite) GarbageCollect() {
	q.Lock()
	defer q.Unlock()

	//create a new slice to hold the data copy the data
	// from the old slice to the new slice and set the
	// internal data to be the new slice
	data := make([]any, 0, cap(q.data))
	copy(data, q.data)
	q.data = data
}

func (q *queueFinite) Resize(newSize int) []any {
	q.Lock()
	defer q.Unlock()

	var items []any

	//ensure that no operations occur if the size hasn't changed,
	// if there's a need to remove items, remove them, then copy the old
	// data to the newly created slice, create new signal channels
	if newSize == cap(q.data) {
		return nil
	}
	if newSize < 1 {
		newSize = 1
	}
	if len(q.data) > newSize {
		items, q.data, _ = internal.DequeueMultiple(len(q.data)-newSize, q.data)
	}
	data := make([]any, len(q.data), newSize)
	copy(data, q.data[:len(q.data)])
	if q.signalIn != nil {
		select {
		default:
			close(q.signalIn)
		case <-q.signalIn:
		}
	}
	if q.signalOut != nil {
		select {
		default:
			close(q.signalOut)
		case <-q.signalOut:
		}
	}
	q.data = data
	q.signalIn = make(chan struct{}, newSize)
	q.signalOut = make(chan struct{}, newSize)
	return items
}

func (q *queueFinite) GetSignalIn() <-chan struct{} {
	q.RLock()
	defer q.RUnlock()
	return q.signalIn
}

func (q *queueFinite) GetSignalOut() <-chan struct{} {
	q.RLock()
	defer q.RUnlock()
	return q.signalOut
}

func (q *queueFinite) Dequeue() (any, bool) {
	q.Lock()
	defer q.Unlock()

	var underflow bool
	var item any

	if item, q.data, underflow = internal.Dequeue(q.data); !underflow {
		internal.SendSignal(q.signalOut)
	}
	return item, underflow
}

func (q *queueFinite) DequeueMultiple(n int) []any {
	q.Lock()
	defer q.Unlock()

	var underflow bool
	var items []any

	if items, q.data, underflow = internal.DequeueMultiple(n, q.data); !underflow {
		internal.SendSignal(q.signalOut)
	}
	return items
}

func (q *queueFinite) Flush() []any {
	q.Lock()
	defer q.Unlock()

	var underflow bool
	var items []any

	if len(q.data) <= 0 {
		return nil
	}
	if items, q.data, underflow = internal.DequeueMultiple(cap(q.data), q.data); !underflow {
		internal.SendSignal(q.signalOut)
	}
	return items
}

func (q *queueFinite) Enqueue(item any) bool {
	q.Lock()
	defer q.Unlock()

	var overflow bool

	if overflow, q.data = internal.Enqueue(q.data, item); !overflow {
		internal.SendSignal(q.signalIn)
	}
	return overflow
}

func (q *queueFinite) EnqueueMultiple(items []any) ([]any, bool) {
	q.Lock()
	defer q.Unlock()

	var remainingElements []any
	var overflow bool

	for i, item := range items {
		if overflow, q.data = internal.Enqueue(q.data, item); overflow {
			remainingElements = items[i:]
			internal.SendSignal(q.signalIn)
			return remainingElements, overflow
		}
	}
	return remainingElements, overflow
}

func (q *queueFinite) EnqueueLossy(item any) (any, bool) {
	q.Lock()
	defer q.Unlock()

	var discardedElement any
	var discard bool

	if len(q.data) >= cap(q.data) {
		discard = true
		discardedElement, q.data, _ = internal.Dequeue(q.data)
	}
	_, q.data = internal.Enqueue(q.data, item)
	internal.SendSignal(q.signalIn)
	return discardedElement, discard
}

func (q *queueFinite) EnqueueInFront(item any) bool {
	q.Lock()
	defer q.Unlock()

	var overflow bool

	if overflow, q.data = internal.EnqueueInFront(q.data, item); !overflow {
		internal.SendSignal(q.signalIn)
	}
	return overflow
}

func (q *queueFinite) Length() int {
	q.RLock()
	defer q.RUnlock()
	return len(q.data)
}

func (q *queueFinite) Capacity() int {
	q.RLock()
	defer q.RUnlock()
	return cap(q.data)
}

func (q *queueFinite) Peek() []any {
	q.RLock()
	defer q.RUnlock()

	items := make([]any, 0, len(q.data))
	for i := range len(q.data) {
		items = append(items, q.data[i])
	}
	return items
}

func (q *queueFinite) PeekHead() (any, bool) {
	q.RLock()
	defer q.RUnlock()

	if len(q.data) <= 0 {
		return nil, true
	}
	return q.data[0], false
}

func (q *queueFinite) PeekFromHead(n int) []any {
	q.RLock()
	defer q.RUnlock()
	if len(q.data) == 0 {
		return nil
	}
	if n > len(q.data) {
		n = len(q.data)
	}
	items := make([]any, 0, n)
	for i := 0; i < n; i++ {
		items = append(items, q.data[i])
	}
	return items
}
