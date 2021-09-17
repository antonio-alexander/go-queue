package infinite

import (
	internal "github.com/antonio-alexander/go-queue/internal"
)

//growIfFull can be used to increase the size of the slice by the
// growSize if the length of the slice is greater than or equal
// to the capacity of the slice
func growIfFull(growSize int, data []interface{}) []interface{} {
	if len(data) < cap(data) {
		return data
	}
	return append(make([]interface{}, 0, cap(data)+growSize), data...)
}

//enqueue can be used to add an item to the back of the slice, if the slice is
// full, it's grown by the growSize, and then the item is appended to the slice
func enqueue(data []interface{}, item interface{}, growSize int) []interface{} {
	data = growIfFull(growSize, data)
	return append(data, item)
}

//enqueueInFront can be used to add an item to the front of the slice, if the slice is
// full, it's grown by the growSize, and then the item is added.
func enqueueInFront(data []interface{}, item interface{}, growSize int) []interface{} {
	data = growIfFull(growSize, data)
	data = append(data, item)
	return internal.RotateRight(data)
}
