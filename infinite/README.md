# infinite (github.com/antonio-alexander/go-queue/infinite)

The infinite "queue" is an implementation of go-queue where the underlying data structure is un-bounded. The queue will automatically re-size itself at runtime in the event it's capacity is reached. Although semantically the same as a finite queue, the "size" provided at the start is not the size of the queue, but the amount to grow the slice by each time the queue is filled.

This is mostly a proof of concept, I think there are some valid use cases for un-bounded queues, but generally it's a code smell and this should never be used in production code outside of testing.

The infinite queue implementation is unique in that it should NEVER overflow since the queue should grow as a result of it being full. Overflow is always false

## Usage

Usage of the infinite queue is straight forward:

1. Create a queue via the New() constructor (supply a size)
2. Use the Enqueue/Dequeue functions to get data in and out of the queue
3. Use the Close() function to clean up the queue

```go
import "github.com/antonio-alexander/go-queue/infinite"

func main() {
    q := infinite.New(1)
    q.Enqueue(1.234)
    item, underflow := q.Dequeue()
    if underflow {
        fmt.Println("underflow occured")
    } else {
        val, _ := item.(float64)
        fmt.Println("value: %f\n", val)
    }
    q.Close()
}
```

All of the patterns introduced in the [github.com/antonio-alexander/go-queue](github.com/antonio-alexander/go-queue) documentation apply for the finite queue.

## Caveats

Keep in mind that internally, the event queue is provided for usability, but because the channel isn't buffered, it can fail if no-one is listening to the other side of the channel. By default this timeout is set to 0, to prevent loss of performance, a polling implementation with flush or DequeueMultiple() is preferred when using infinite queues.

The "grow" size provided on instantiation is a source of tuning, in the event the size is too small, you may allocate data more often and burn cpu usage, if the size is too big, you may allocate more data than you actually need.
