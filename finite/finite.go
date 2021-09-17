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
	data      []interface{}
}

func New(size int) interface {
	goqueue.Owner
	goqueue.GarbageCollecter
	goqueue.Dequeuer
	goqueue.Enqueuer
	goqueue.EnqueueInFronter
	goqueue.Info
	goqueue.Event
	goqueue.Peeker
	EnqueueLossy
	Resizer
} {

	maxSize := size
	if maxSize < 1 {
		maxSize = 1
	}
	return &queueFinite{
		signalIn:  make(chan struct{}, maxSize),
		signalOut: make(chan struct{}, maxSize),
		data:      make([]interface{}, 0, maxSize),
	}
}

func (q *queueFinite) Close() (remainingElements []interface{}) {
	q.Lock()
	defer q.Unlock()

	remainingElements, q.data, _ = internal.DequeueMultiple(cap(q.data), q.data)
	close(q.signalIn)
	close(q.signalOut)
	q.data, q.signalIn, q.signalOut = nil, nil, nil

	return
}

func (q *queueFinite) GarbageCollect() {
	q.Lock()
	defer q.Unlock()

	//create a new slice to hold the data copy the data
	// from the old slice to the new slice and set the
	// internal data to be the new slice
	data := make([]interface{}, 0, cap(q.data))
	copy(data, q.data)
	q.data = data
}

func (q *queueFinite) Resize(newSize int) (items []interface{}) {
	q.Lock()
	defer q.Unlock()

	//ensure that no operations occur if the size hasn't changed,
	// if there's a need to remove items, remove them, then copy the old
	// data to the newly created slice, create new signal channels
	if newSize == cap(q.data) {
		return
	}
	if newSize < 1 {
		newSize = 1
	}
	if len(q.data) > newSize {
		items, q.data, _ = internal.DequeueMultiple(len(q.data)-newSize, q.data)
	}
	data := make([]interface{}, len(q.data), newSize)
	copy(data, q.data[:len(q.data)])
	close(q.signalIn)
	close(q.signalOut)
	q.data = data
	q.signalIn = make(chan struct{}, newSize)
	q.signalOut = make(chan struct{}, newSize)

	return
}

func (q *queueFinite) GetSignalIn() (signal <-chan struct{}) {
	q.RLock()
	defer q.RUnlock()
	return q.signalIn
}

func (q *queueFinite) GetSignalOut() (signal <-chan struct{}) {
	q.RLock()
	defer q.RUnlock()
	return q.signalOut
}

func (q *queueFinite) Dequeue() (item interface{}, underflow bool) {
	q.Lock()
	defer q.Unlock()

	if item, q.data, underflow = internal.Dequeue(q.data); !underflow {
		internal.SendSignal(q.signalOut)
	}

	return
}

func (q *queueFinite) DequeueMultiple(n int) (items []interface{}) {
	q.Lock()
	defer q.Unlock()

	var underflow bool

	if items, q.data, underflow = internal.DequeueMultiple(n, q.data); !underflow {
		internal.SendSignal(q.signalOut)
	}

	return
}

func (q *queueFinite) Flush() (items []interface{}) {
	q.Lock()
	defer q.Unlock()

	var underflow bool

	if len(q.data) <= 0 {
		return
	}
	if items, q.data, underflow = internal.DequeueMultiple(cap(q.data), q.data); !underflow {
		internal.SendSignal(q.signalOut)
	}

	return
}

func (q *queueFinite) Enqueue(item interface{}) (overflow bool) {
	q.Lock()
	defer q.Unlock()

	if overflow, q.data = internal.Enqueue(q.data, item); !overflow {
		internal.SendSignal(q.signalIn)
	}

	return
}

func (q *queueFinite) EnqueueMultiple(items []interface{}) (remainingElements []interface{}, overflow bool) {
	q.Lock()
	defer q.Unlock()

	for i, item := range items {
		if overflow, q.data = internal.Enqueue(q.data, item); overflow {
			remainingElements = items[i:]
			internal.SendSignal(q.signalIn)

			return
		}
	}

	return
}

func (q *queueFinite) EnqueueLossy(item interface{}) (discardedElement interface{}, discard bool) {
	q.Lock()
	defer q.Unlock()

	if len(q.data) >= cap(q.data) {
		discard = true
		discardedElement, q.data, _ = internal.Dequeue(q.data)
	}
	_, q.data = internal.Enqueue(q.data, item)
	internal.SendSignal(q.signalIn)

	return
}

func (q *queueFinite) EnqueueInFront(item interface{}) (overflow bool) {
	q.Lock()
	defer q.Unlock()

	if overflow, q.data = internal.EnqueueInFront(q.data, item); !overflow {
		internal.SendSignal(q.signalIn)
	}

	return
}

func (q *queueFinite) Length() (size int) {
	q.RLock()
	defer q.RUnlock()
	return len(q.data)
}

func (q *queueFinite) Capacity() (capacity int) {
	q.RLock()
	defer q.RUnlock()
	return cap(q.data)
}

func (q *queueFinite) Peek() (items []interface{}) {
	q.RLock()
	defer q.RUnlock()

	for i := 0; i < len(q.data); i++ {
		items = append(items, q.data[i])
	}

	return
}

func (q *queueFinite) PeekHead() (item interface{}, underflow bool) {
	q.RLock()
	defer q.RUnlock()
	if len(q.data) <= 0 {
		return nil, true
	}
	return q.data[0], false
}

func (q *queueFinite) PeekFromHead(n int) (items []interface{}) {
	q.RLock()
	defer q.RUnlock()

	if len(q.data) == 0 {
		return
	}
	if n > len(q.data) {
		n = len(q.data)
	}
	for i := 0; i < n; i++ {
		items = append(items, q.data[i])
	}

	return
}
