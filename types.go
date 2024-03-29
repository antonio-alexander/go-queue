package goqueue

import (
	"encoding"
)

//These types are specifically provided to attempt to communicate support
// for how queues would be able to store data in a persistent way no matter
// the data type (empty interface)
type (
	BinaryMarshaler = encoding.BinaryMarshaler

	BinaryUnmarshaler = encoding.BinaryUnmarshaler

	//Bytes is provided to make it easier to create jagged arrays; two
	// dimensional arrays are nice, but they work off the idea that
	// each row has the same number of elements which doesnt work for
	// the use case for a queue...
	//KIM: Bytes is a type that can be used to traverse package boundaries
	// unlike an anonymous struct or Bytes defined by some other package
	Bytes []byte
)

//Owner provides functions that directly affect the underlying pointers
// and data structures of a queue pointers. The Close() function should
// ready the underlying pointer for garbage collection and return a slice
// of any items that remain in the queue
type Owner interface {
	Close() (items []interface{})
}

//GarbageCollecter can be implemented to re-create the underlying pointers
// so that they can be garabge collected, you can think of this as creating
// an opportunity to defrag the memory
type GarbageCollecter interface {
	GarbageCollect()
}

//Dequeuer can be used to destructively remove one or more items from the
// queue, it can remove one item via Dequeue(), multiple items via
// DequeueMultiple() or all items using Flush() underflow will be true if
// the queue is empty
type Dequeuer interface {
	Dequeue() (item interface{}, underflow bool)
	DequeueMultiple(n int) (items []interface{})
	Flush() (items []interface{})
}

//Peeker can be used to non-destructively remove one or more items from
// the queue, it can remove all items via Peek(), remove an item from the
// front of the queue via PeekHead() or remove multiple items via
// PeekFromHead(). Underflow will be true, if the queue is empty
type Peeker interface {
	Peek() (items []interface{})
	PeekHead() (item interface{}, underflow bool)
	PeekFromHead(n int) (items []interface{})
}

//Enqueuer can be used to put one or more items into the queue
// Enqueue() can be used to place one item while EnqueueMultiple()
// can be used to place multiple items, in the event the queue is full
// the remaining items will be provided (if applicable) and overflow
// will be true
type Enqueuer interface {
	Enqueue(item interface{}) (overflow bool)
	EnqueueMultiple(items []interface{}) (itemsRemaining []interface{}, overflow bool)
}

//EnqueueInFronter describes an operation where you enqueue a single item at the
// front of the queue, if the queue is full overflow will be true
type EnqueueInFronter interface {
	EnqueueInFront(item interface{}) (overflow bool)
}

//Length can be used to determine how many items are inside a queue at
// any given time
type Length interface {
	Length() (size int)
}

//Event can be used to get a read-only signal that would indicate whether data was
// removed from the queue (out) or put into the queue (in). Keep in mind that whether
// the channel is buffered or un-buffered depends on the underlying implementation
type Event interface {
	GetSignalIn() (signal <-chan struct{})
	GetSignalOut() (signal <-chan struct{})
}
