package infinite

import (
	"math"
	"sync"

	goqueue "github.com/antonio-alexander/go-queue"
	internal "github.com/antonio-alexander/go-queue/internal"
)

type queueInfinite struct {
	sync.RWMutex
	growSize  int
	signalIn  chan struct{}
	signalOut chan struct{}
	data      []any
}

func New(growSize int) interface {
	goqueue.Owner
	goqueue.GarbageCollecter
	goqueue.Dequeuer
	goqueue.Enqueuer
	goqueue.EnqueueInFronter
	goqueue.Length
	goqueue.Event
	goqueue.Peeker
} {
	growSize = max(growSize, 1)
	return &queueInfinite{
		growSize:  growSize,
		data:      make([]any, 0, growSize),
		signalIn:  make(chan struct{}),
		signalOut: make(chan struct{}),
	}
}

func (q *queueInfinite) Close() []any {
	q.Lock()
	defer q.Unlock()

	var remainingElements []any

	remainingElements, q.data, _ = internal.DequeueMultiple(cap(q.data), q.data)
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
	q.growSize = 0
	return remainingElements
}

func (q *queueInfinite) GarbageCollect() {
	q.Lock()
	defer q.Unlock()

	var length, newSize, r int
	var data []any

	//this collection will attempt to create a new underlying data structure and
	// down-size it if it's grown more than necessary
	length = len(q.data)
	newSize = int(math.Trunc(float64(length)/float64(q.growSize)) * float64(q.growSize))
	if r = length % q.growSize; r > 0 || newSize == 0 {
		newSize += q.growSize
	}
	data = make([]any, len(q.data), newSize)
	copy(data, q.data[:len(q.data)])
	q.data = data
}

func (q *queueInfinite) Dequeue() (any, bool) {
	q.Lock()
	defer q.Unlock()

	var underflow bool
	var item any

	item, q.data, underflow = internal.Dequeue(q.data)
	internal.SendSignal(q.signalOut, ConfigSignalTimeout)
	return item, underflow
}

func (q *queueInfinite) DequeueMultiple(n int) []any {
	q.Lock()
	defer q.Unlock()

	var underflow bool
	var items []any

	if items, q.data, underflow = internal.DequeueMultiple(n, q.data); !underflow {
		internal.SendSignal(q.signalOut, ConfigSignalTimeout)
	}
	return items
}

func (q *queueInfinite) Flush() []any {
	q.Lock()
	defer q.Unlock()

	var underflow bool
	var items []any

	if items, q.data, underflow = internal.DequeueMultiple(cap(q.data), q.data); !underflow {
		internal.SendSignal(q.signalOut, ConfigSignalTimeout)
	}
	return items
}

func (q *queueInfinite) Enqueue(item any) bool {
	q.Lock()
	defer q.Unlock()

	q.data = enqueue(q.data, item, q.growSize)
	internal.SendSignal(q.signalIn, ConfigSignalTimeout)
	return false
}

func (q *queueInfinite) EnqueueMultiple(items []any) ([]any, bool) {
	q.Lock()
	defer q.Unlock()

	for _, item := range items {
		q.data = enqueue(q.data, item, q.growSize)
		internal.SendSignal(q.signalIn, ConfigSignalTimeout)
	}
	return nil, false
}

func (q *queueInfinite) EnqueueInFront(item any) bool {
	q.Lock()
	defer q.Unlock()

	q.data = enqueueInFront(q.data, item, q.growSize)
	internal.SendSignal(q.signalIn, ConfigSignalTimeout)
	return false
}

func (q *queueInfinite) Length() int {
	q.RLock()
	defer q.RUnlock()
	return len(q.data)
}

func (q *queueInfinite) GetSignalIn() <-chan struct{} {
	q.RLock()
	defer q.RUnlock()
	return q.signalIn
}

func (q *queueInfinite) GetSignalOut() <-chan struct{} {
	q.RLock()
	defer q.RUnlock()
	return q.signalOut
}

func (q *queueInfinite) Peek() []any {
	q.RLock()
	defer q.RUnlock()

	items := make([]any, 0, len(q.data))
	for i := range len(q.data) {
		items = append(items, q.data[i])
	}
	return items
}

func (q *queueInfinite) PeekHead() (any, bool) {
	q.RLock()
	defer q.RUnlock()

	if len(q.data) <= 0 {
		return nil, true
	}
	return q.data[0], false
}

func (q *queueInfinite) PeekFromHead(n int) []any {
	q.RLock()
	defer q.RUnlock()

	if len(q.data) == 0 {
		return nil
	}
	if n > len(q.data) {
		n = len(q.data)
	}
	items := make([]any, 0, n)
	for i := range n {
		items = append(items, q.data[i])
	}
	return items
}
