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

## Testing

These are the avaialble unit tests:

- GarbageCollect
- Enqueue
- EnqueueMultiple
- EnqueueInFront

These are the available function tests:

- Queue: can be used to verify that queue works as expected (for infinite queues)

## Event-based operations

```go
type Event interface {
    GetSignalIn() (signal <-chan struct{})
    GetSignalOut() (signal <-chan struct{})
}
```

```go
//SendSignal will perform a non-blocking send with or without
// a timeout depending on whether ConfigSignalTimeout is greater
// than 0
func SendSignal(signal chan struct{}, timeout time.Duration) bool {
    if timeout > 0 {
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
```

infinite implements the Event interface from go-queue. Since the implementation uses an unbuffered channel (since we can never know the total size of the queue) there's no guarantee that attempts to send won't block forever, as a result, the "timeout" for waiting for signals to be "read" can be configured via the package variable "ConfigSignalTimeout". It's implemented such that the timeout is used (if configured) and it'll wait until the time.After() signals or an entities reads from the unbuffered channel.

As a result, the following assumptions can be made:

- If no timeout is configured, if no-one is reading the channel, the signal will be missed
- If a timeout is configured, and no-one reads the channel before timing out
- Event based dequeue/enqueue operations should take into account the possility that a signal in/out may be missed and you'll need to have a ticker that dequeues with time, or some other logic to capture situations where you've missed signals and need to dequeue
