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
	data      []interface{}
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
	if growSize < 1 {
		growSize = 1
	}
	return &queueInfinite{
		growSize:  growSize,
		data:      make([]interface{}, 0, growSize),
		signalIn:  make(chan struct{}),
		signalOut: make(chan struct{}),
	}
}

func (q *queueInfinite) Close() (remainingElements []interface{}) {
	q.Lock()
	defer q.Unlock()

	remainingElements, q.data, _ = internal.DequeueMultiple(cap(q.data), q.data)
	close(q.signalIn)
	close(q.signalOut)
	q.data, q.signalIn, q.signalOut = nil, nil, nil
	q.growSize = 0

	return
}

func (q *queueInfinite) GarbageCollect() {
	q.Lock()
	defer q.Unlock()

	var length, newSize, r int
	var data []interface{}

	//this collection will attempt to create a new underlying data structure and
	// down-size it if it's grown more than necessary
	length = len(q.data)
	newSize = int(math.Trunc(float64(length)/float64(q.growSize)) * float64(q.growSize))
	if r = length % q.growSize; r > 0 || newSize == 0 {
		newSize += q.growSize
	}
	data = make([]interface{}, len(q.data), newSize)
	copy(data, q.data[:len(q.data)])
	q.data = data
}

func (q *queueInfinite) Dequeue() (item interface{}, underflow bool) {
	q.Lock()
	defer q.Unlock()

	item, q.data, underflow = internal.Dequeue(q.data)
	internal.SendSignal(q.signalOut, ConfigSignalTimeout)

	return
}

func (q *queueInfinite) DequeueMultiple(n int) (items []interface{}) {
	q.Lock()
	defer q.Unlock()

	var underflow bool

	if items, q.data, underflow = internal.DequeueMultiple(n, q.data); !underflow {
		internal.SendSignal(q.signalOut, ConfigSignalTimeout)
	}

	return
}

func (q *queueInfinite) Flush() (items []interface{}) {
	q.Lock()
	defer q.Unlock()

	var underflow bool

	if items, q.data, underflow = internal.DequeueMultiple(cap(q.data), q.data); !underflow {
		internal.SendSignal(q.signalOut, ConfigSignalTimeout)
	}

	return
}

func (q *queueInfinite) Enqueue(item interface{}) (overflow bool) {
	q.Lock()
	defer q.Unlock()

	q.data = enqueue(q.data, item, q.growSize)
	internal.SendSignal(q.signalIn, ConfigSignalTimeout)

	return
}

func (q *queueInfinite) EnqueueMultiple(items []interface{}) (remainingElements []interface{}, overflow bool) {
	q.Lock()
	defer q.Unlock()

	for _, item := range items {
		q.data = enqueue(q.data, item, q.growSize)
		internal.SendSignal(q.signalIn, ConfigSignalTimeout)
	}

	return
}

func (q *queueInfinite) EnqueueInFront(item interface{}) (overflow bool) {
	q.Lock()
	defer q.Unlock()

	q.data = enqueueInFront(q.data, item, q.growSize)
	internal.SendSignal(q.signalIn, ConfigSignalTimeout)

	return
}

func (q *queueInfinite) Length() (size int) {
	q.RLock()
	defer q.RUnlock()
	return len(q.data)
}

func (q *queueInfinite) GetSignalIn() (signal <-chan struct{}) {
	q.RLock()
	defer q.RUnlock()
	return q.signalIn
}

func (q *queueInfinite) GetSignalOut() (signal <-chan struct{}) {
	q.RLock()
	defer q.RUnlock()
	return q.signalOut
}

func (q *queueInfinite) Peek() (items []interface{}) {
	q.RLock()
	defer q.RUnlock()

	for i := 0; i < len(q.data); i++ {
		items = append(items, q.data[i])
	}

	return
}

func (q *queueInfinite) PeekHead() (item interface{}, underflow bool) {
	q.RLock()
	defer q.RUnlock()

	if len(q.data) <= 0 {
		return nil, true
	}
	return q.data[0], false
}

func (q *queueInfinite) PeekFromHead(n int) (items []interface{}) {
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
