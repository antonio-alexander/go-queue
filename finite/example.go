package finite

import (
	goqueue "github.com/antonio-alexander/go-queue"
)

func ExampleResize(queue Resizer, size int) []*goqueue.Example {
	items := queue.Resize(size)
	return goqueue.ExampleConvertMultiple(items)
}

func ExampleEnqueueLossy(queue EnqueueLossy, value *goqueue.Example) (*goqueue.Example, bool) {
	item, discarded := queue.EnqueueLossy(value)
	if !discarded {
		return nil, false
	}
	return goqueue.ExampleConvertSingle(item), true
}
