# finite (github.com/antonio-alexander/go-queue/finite)

The finite "queue" is an implementation of go-queue where the underlying data structure is finite (or bounded). Although the queue can be re-sized at runtime, it's expected that while "in-use" it's a fixed size.

The backing data structure is a slice of empty interface, be careful that even though the slice is "fixed", the actual memory it allocates depends on the data placed inside of it; this can generally be mitigated by use of pointers rather than structs.

The finite queue implementation is unique in that it is...wait for it...finite. This means that if you attempt to enqueue items when the queue is full, overflow will be true.

## Usage

Usage of the finite queue is straight forward:

1. Create a queue via the New() constructor (supply a size)
2. Use the Enqueue/Dequeue functions to get data in and out of the queue
3. Use the Close() function to clean up the queue

```go
import "github.com/antonio-alexander/go-queue/finite"

func main() {
    q := finite.New(1)
    overflow := q.Enqueue(1.234);
    if overflow {
        fmt.Println("overflow occured")
    }
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

## Finite Interfaces

The finite queue implementation is specific in that it will overflow if the queue is full. With that in mind, it has two additional interfaces not provided by the go-queue package. The Resize() function and the EnqueueLossy() function.

The implementation was designed such that the "size" is provided on instantiation to simplify the implementation, but in the event you are unable to know the true size at instantiation, the Reize() function can resolve that for you. It will re-create the underlying pointers (the buffered signal channels AND the slice of empty interface) at the given size and in the event it's less than the current size, it will provide any items that have been removed.

Keep in mind that Resize() is destructive and will immediately invalidate any signal channels that were already called, this should be done synchronously such that any connected producers or consumers aren't attempting to produce or consume.

```go
type Resizer interface {
    Resize(size int) (items []interface{})
}
```

EnqueueLossy is definitely a holdover from LabVIEW, but it allows you to have a finite queue while keeping the most recent data. When executed, if the queue is full, it will discard the oldest item at the front of the queue and place the new item in the back of the queue. This somewhat contrasts with EnqueueInFront().

```go
type EnqueueLossy interface {
    EnqueueLossy(item interface{}) (discardedElement interface{}, discard bool)
}
```

## Patterns

Single element queue:

```go
//TODO: write code that shows a single element queue
```

Running Average:

```go
//TODO: generate code to show a running average calculation
```

## Testing

Similar to the [github.com/antonio-alexander/go-queue](github.com/antonio-alexander/go-queue) package, some of the interfaces (or implementations) specific to the finite implementation, have tests that can be used to verify that your implementation matches the expected finite implementation.

These are the avaialble unit tests:

- Resize: can be used to verify the Resize() function
- Enqueue: can be used to verify the Enqueue() function for finite queues
- EnqueueMultiple: can be used to verify the EnqueueMultiple() function for finite queues
- EnqueueInFront: can be used to verify the EnqueueInFront() function for finite queues
- EnqueueLossy: can be used to verify the EnqueueLossy() function
