package goqueue

import "time"

//MustEnqueue will attempt to use the Enqueue() function until the enqueue is successful
// (no overflow); this function will block until success occurs or the done channel receives
// a signal. An enqueue will attempt to occur at the rate configured
func MustEnqueue(queue Enqueuer, item interface{}, done <-chan struct{}, rate time.Duration) bool {
	if overflow := queue.Enqueue(item); !overflow {
		return overflow
	}
	tEnqueue := time.NewTicker(rate)
	defer tEnqueue.Stop()
	if done != nil {
		for {
			select {
			case <-done:
				return queue.Enqueue(item)
			case <-tEnqueue.C:
				if overflow := queue.Enqueue(item); !overflow {
					return overflow
				}
			}
		}
	}
	for {
		<-tEnqueue.C
		if overflow := queue.Enqueue(item); !overflow {
			return overflow
		}
	}
}

//MustEnqueue will attempt to use the Enqueue() function until the enqueue is successful
// (no overflow); this function will block until success occurs or the done channel receives
// a signal. An enqueue will be attempted for every signal received
func MustEnqueueEvent(queue interface {
	Enqueuer
	Event
}, item interface{}, done <-chan struct{}) bool {
	if overflow := queue.Enqueue(item); !overflow {
		return false
	}
	signalOut := queue.GetSignalOut()
	if done != nil {
		for {
			select {
			case <-done:
				return queue.Enqueue(item)
			case <-signalOut:
				if overflow := queue.Enqueue(item); !overflow {
					return overflow
				}
			}
		}
	}
	for {
		<-signalOut
		if overflow := queue.Enqueue(item); !overflow {
			return overflow
		}
	}
}

//MustEnqueueMultiple will attempt to enqueue until the done channel completes,
//  at the configured rate or the number of elements are successfully enqueued
//  into the provided queue
//KIM: this function doesn't preserve the unit of work and may not be consistent
// with concurent usage (although it is safe)
func MustEnqueueMultiple(queue Enqueuer, items []interface{}, done <-chan struct{}, rate time.Duration) ([]interface{}, bool) {
	itemsRemaining, overflow := queue.EnqueueMultiple(items)
	if !overflow {
		return nil, false
	}
	items = itemsRemaining
	tEnqueueMultiple := time.NewTicker(rate)
	defer tEnqueueMultiple.Stop()
	if done != nil {
		for {
			select {
			case <-done:
				return queue.EnqueueMultiple(items)
			case <-tEnqueueMultiple.C:
				itemsRemaining, overflow := queue.EnqueueMultiple(items)
				if !overflow {
					return nil, false
				}
				items = itemsRemaining
			}
		}
	}
	for {
		<-tEnqueueMultiple.C
		itemsRemaining, overflow := queue.EnqueueMultiple(items)
		if !overflow {
			return nil, false
		}
		items = itemsRemaining
	}
}

//MustEnqueueMultipleEvent will attempt to enqueue one or more items, upon initial
// failure, it'll use the event channels/signals to attempt to enqueue items
//KIM: this function doesn't preserve the unit of work and may not be consistent
// with concurent usage (although it is safe)
func MustEnqueueMultipleEvent(queue interface {
	Enqueuer
	Event
}, items []interface{}, done <-chan struct{}) ([]interface{}, bool) {
	itemsRemaining, overflow := queue.EnqueueMultiple(items)
	if !overflow {
		return nil, false
	}
	items = itemsRemaining
	signalOut := queue.GetSignalOut()
	if done != nil {
		for {
			select {
			case <-done:
				return queue.EnqueueMultiple(items)
			case <-signalOut:
				itemsRemaining, overflow := queue.EnqueueMultiple(items)
				if !overflow {
					return nil, false
				}
				items = itemsRemaining
			}
		}
	}
	for {
		<-signalOut
		itemsRemaining, overflow := queue.EnqueueMultiple(items)
		if !overflow {
			return nil, false
		}
		items = itemsRemaining
	}
}

//MustDequeue will attempt to dequeue at least one item at the rate configured until
// the done channel signals.
//KIM: It's possible to provide a nil channel and this function will block (forever)
// until a dequeue is successful
func MustDequeue(queue Dequeuer, done <-chan struct{}, rate time.Duration) (interface{}, bool) {
	if item, underflow := queue.Dequeue(); !underflow {
		return item, false
	}
	tDequeue := time.NewTicker(rate)
	defer tDequeue.Stop()
	if done != nil {
		for {
			select {
			case <-done:
				return queue.Dequeue()
			case <-tDequeue.C:
				if item, underflow := queue.Dequeue(); !underflow {
					return item, false
				}
			}
		}
	}
	for {
		<-tDequeue.C
		if item, underflow := queue.Dequeue(); !underflow {
			return item, false
		}
	}
}

func MustDequeueEvent(queue interface {
	Dequeuer
	Event
}, done <-chan struct{}) (interface{}, bool) {
	signalIn := queue.GetSignalIn()
	if item, underflow := queue.Dequeue(); !underflow {
		return item, false
	}
	if done != nil {
		for {
			select {
			case <-done:
				return queue.Dequeue()
			case <-signalIn:
				if item, underflow := queue.Dequeue(); !underflow {
					return item, false
				}
			}
		}
	}
	for {
		<-signalIn
		if item, underflow := queue.Dequeue(); !underflow {
			return item, false
		}
	}
}

func MustDequeueMultiple(queue Dequeuer, done <-chan struct{}, n int, rate time.Duration) []interface{} {
	items := queue.DequeueMultiple(n)
	if len(items) == n {
		return items
	}
	n = n - len(items)
	tDequeueMultiple := time.NewTicker(rate)
	defer tDequeueMultiple.Stop()
	if done != nil {
		for {
			select {
			case <-done:
				return append(items, queue.DequeueMultiple(n)...)
			case <-tDequeueMultiple.C:
				items = append(items, queue.DequeueMultiple(n)...)
				if len(items) == n {
					return items
				}
				n = n - len(items)
			}
		}
	}
	for {
		<-tDequeueMultiple.C
		items := queue.DequeueMultiple(n)
		if len(items) == n {
			return items
		}
		n = n - len(items)
	}
}

func MustDequeueMultipleEvent(queue interface {
	Dequeuer
	Event
}, done <-chan struct{}, n int) []interface{} {
	items := queue.DequeueMultiple(n)
	if len(items) == n {
		return items
	}
	n = n - len(items)
	signalIn := queue.GetSignalIn()
	if done != nil {
		for {
			select {
			case <-done:
				return append(items, queue.DequeueMultiple(n))
			case <-signalIn:
				items := queue.DequeueMultiple(n)
				if len(items) == n {
					return items
				}
				n = n - len(items)
			}
		}
	}
	for {
		<-signalIn
		items := queue.DequeueMultiple(n)
		if len(items) == n {
			return items
		}
		n = n - len(items)
	}
}

func MustFlush(queue Dequeuer, done <-chan struct{}, rate time.Duration) []interface{} {
	items := queue.Flush()
	tFlush := time.NewTicker(rate)
	defer tFlush.Stop()
	if done != nil {
		for {
			select {
			case <-done:
				return append(items, queue.Flush()...)
			case <-tFlush.C:
				items = append(items, queue.Flush()...)
			}
		}
	}
	for {
		<-tFlush.C
		items = append(items, queue.Flush()...)
	}
}

func MustFlushEvent(queue interface {
	Dequeuer
	Event
}, done <-chan struct{}) []interface{} {
	items := queue.Flush()
	signalIn := queue.GetSignalIn()
	if done != nil {
		for {
			select {
			case <-done:
				return append(items, queue.Flush()...)
			case <-signalIn:
				items = append(items, queue.Flush()...)
			}
		}
	}
	for {
		<-signalIn
		items = append(items, queue.Flush()...)
	}
}

func MustPeekHead(queue Peeker, done <-chan struct{}, rate time.Duration) (interface{}, bool) {
	if item, underflow := queue.PeekHead(); !underflow {
		return item, false
	}
	tPeekHead := time.NewTicker(rate)
	defer tPeekHead.Stop()
	if done != nil {
		for {
			select {
			case <-done:
				return queue.PeekHead()
			case <-tPeekHead.C:
				if item, underflow := queue.PeekHead(); !underflow {
					return item, false
				}
			}
		}
	}
	for {
		<-tPeekHead.C
		if item, underflow := queue.PeekHead(); !underflow {
			return item, false
		}
	}
}

func MustPeekHeadEvent(queue interface {
	Peeker
	Event
}, done <-chan struct{}) (interface{}, bool) {
	signalIn := queue.GetSignalIn()
	if item, underflow := queue.PeekHead(); !underflow {
		return item, false
	}
	if done != nil {
		for {
			select {
			case <-done:
				return queue.PeekHead()
			case <-signalIn:
				if item, underflow := queue.PeekHead(); !underflow {
					return item, false
				}
			}
		}
	}
	for {
		<-signalIn
		if item, underflow := queue.PeekHead(); !underflow {
			return item, false
		}
	}
}

func MustPeekFromHead(queue Peeker, done <-chan struct{}, n int, rate time.Duration) []interface{} {
	items := queue.PeekFromHead(n)
	if len(items) == n {
		return items
	}
	n = n - len(items)
	tPeekFromHead := time.NewTicker(rate)
	defer tPeekFromHead.Stop()
	if done != nil {
		for {
			select {
			case <-done:
				return append(items, queue.PeekFromHead(n))
			case <-tPeekFromHead.C:
				items := queue.PeekFromHead(n)
				if len(items) == n {
					return items
				}
				n = n - len(items)
			}
		}
	}
	for {
		<-tPeekFromHead.C
		items := queue.PeekFromHead(n)
		if len(items) == n {
			return items
		}
		n = n - len(items)
	}
}

func MustPeekFromHeadEvent(queue interface {
	Peeker
	Event
}, done <-chan struct{}, n int) []interface{} {
	items := queue.PeekFromHead(n)
	if len(items) == n {
		return items
	}
	n = n - len(items)
	signalIn := queue.GetSignalIn()
	if done != nil {
		for {
			select {
			case <-done:
				return append(items, queue.PeekFromHead(n))
			case <-signalIn:
				items := queue.PeekFromHead(n)
				if len(items) == n {
					return items
				}
				n = n - len(items)
			}
		}
	}
	for {
		<-signalIn
		items := queue.PeekFromHead(n)
		if len(items) == n {
			return items
		}
		n = n - len(items)
	}
}

func MustPeek(queue Peeker, done <-chan struct{}, rate time.Duration) []interface{} {
	items := queue.Peek()
	tPeek := time.NewTicker(rate)
	defer tPeek.Stop()
	if done != nil {
		for {
			select {
			case <-done:
				return append(items, queue.Peek()...)
			case <-tPeek.C:
				items = append(items, queue.Peek()...)
			}
		}
	}
	for {
		<-tPeek.C
		items = append(items, queue.Peek()...)
	}
}

func MustPeekEvent(queue interface {
	Peeker
	Event
}, done <-chan struct{}) []interface{} {
	items := queue.Peek()
	signalIn := queue.GetSignalIn()
	if done != nil {
		for {
			select {
			case <-done:
				return append(items, queue.Peek()...)
			case <-signalIn:
				items = append(items, queue.Peek()...)
			}
		}
	}
	for {
		<-signalIn
		items = append(items, queue.Peek()...)
	}
}
