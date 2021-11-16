package finite

//Resizer can be used to modify the size of the queue, it will return any elements
// that can't fit in the new queue. Keep in mind that this is destructive and will
// invalidate and signal channels that have been created
type Resizer interface {
	Resize(size int) (items []interface{})
}

//EnqueueLossy can be used to add an element to the back of the queue, if
// the queue is full, the oldest element will be discarded and returned
type EnqueueLossy interface {
	EnqueueLossy(item interface{}) (discardedElement interface{}, discard bool)
}
