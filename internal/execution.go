package internal

import (
	"time"
)

//RotateLeft can be used to perform an in-place rotation
// left of a slice of empty interface
func RotateLeft(dataIn []interface{}) []interface{} {
	if len(dataIn) > 1 {
		copy(dataIn, append(dataIn[1:], dataIn[:1]...))
		return dataIn
	}
	return dataIn
}

//RotateRight can be used to perform an in-place rotation
// right of a slice of empty interface
func RotateRight(dataIn []interface{}) []interface{} {
	if len(dataIn) > 1 {
		copy(dataIn, append(dataIn[len(dataIn)-1:], dataIn[:len(dataIn)-1]...))
		return dataIn
	}
	return dataIn
}

//Enqueue can be used  to add an item to the back of a queue while maintaining
// it's capacity (e.g. in-place) it will return true if the queue is full
func Enqueue(dataIn []interface{}, item interface{}) (bool, []interface{}) {
	if len(dataIn) >= cap(dataIn) {
		return true, dataIn
	}
	return false, append(dataIn, item)
}

//EnqueueInFront can be used to add an item to the front of the queue while maintaining
// its capacity, it will return true if the queue is full
func EnqueueInFront(data []interface{}, item interface{}) (bool, []interface{}) {
	if len(data) >= cap(data) {
		return true, data
	}
	data = append(data, item)
	return false, RotateRight(data)
}

//Dequeue can be used to remove an item from the queue and reduce its
// capacity by one
func Dequeue(data []interface{}) (interface{}, []interface{}, bool) {
	if len(data) <= 0 {
		return nil, data, true
	}
	item := data[0]
	data[0] = nil
	data = RotateLeft(data)
	if len(data) > 0 {
		data = data[:len(data)-1] //truncate the slice
	}
	return item, data, false
}

//DequeueMultiple will return a number of items less than or equal to the value of
//n while maintaining the input data on the second slice of interface, it will return
// true if there are no items to dequeue
func DequeueMultiple(n int, data []interface{}) ([]interface{}, []interface{}, bool) {
	var l int

	//get the length of the data, underflow if no data, then
	// check to see if n is negative or greater than -1
	if l = len(data); l <= 0 {
		return nil, data, true
	}
	if n > l {
		n = l
	}
	items := make([]interface{}, 0, n)
	for i := 0; i < n; i++ {
		var item interface{}

		item, data, _ = Dequeue(data)
		items = append(items, item)
	}
	return items, data, false
}

//SendSignal will perform a non-blocking send with or without
// a timeout depending on whether ConfigSignalTimeout is greater
// than 0
func SendSignal(signal chan struct{}) bool {
	if timeout := ConfigSignalTimeout; timeout > 0 {
		select {
		case <-time.After(timeout):
		case signal <- struct{}{}:
			return true
		}
		return false
	}
	select {
	default:
	case signal <- struct{}{}:
		return true
	}
	return false
}
