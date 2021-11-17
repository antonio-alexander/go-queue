# github.com/antonio-alexander/go-queue

go-queue, is a FIFO data structure that is a functional replacement for channels. It's opinion is that channels are good for synchronously signaling that there is data available, but not good at communicating that data. There are a number of API that completely separate this signaling from the destructive (or non-destructive) reading of data. I think that go-queue should be used in situations where a channel is long-lived and/or you require non-destructive access to data while maintaining the FIFO.

Here are some common situations where go-queue functionality would be advantageous to using channels:

- If you want to "peek" at the data at the head of the queue to perform work on it before removing it from the queue (e.g. if the attempted "work" on that item failed, you've already removed it from the channel so you can't put it back)
- If you want to put an item at the front of the queue when there are items in the queue
- If you want to remove all data from the queue at once (e.g. this is almost 100% necessary for high throughput then the consumer can run faster than the producer)
- You're using the producer/consumer pattern and you want to avoid polling (you can use a select case with a time.After or ticker, but you're still forced to "poll" whether the channel has data in it)
- You have need of a channel, but your business logic needs the channel's size to grow at runtime (I think this is a code smell, but I've dealt with stranger things)
- You want to know how many items are in the channel

## Queue interfaces

go-queue is separated into a high-level/common go module [github.com/antonio-alexander/go-queue](github.com/antonio-alexander/go-queue) where the interfaces (described below) and tests are defined that can be imported/used by anyone attempting to implement those interfaces.

> If it's not obvious, the goal of this separation of ownership of interfaces is used such that anyone using queues depend on the interface, not the implementation

Keep in mind that some of these functions are dependent on the underlying implementation; for example overflow and capacity will have different output depending on if the queue is finite or infinite.

Owner, similar to GarbageCollector(), defines functions that operate on the underlying pointer. The Close() function will ready the underlying pointer for garbage collection and return any items that remain in the queue.

```go
type Owner interface {
    Close() (items []interface{})
}
```

GarbageCollecter can be used to perform a kind of defragmentation of memory. Generally because the queue implementations are backed by a slice, depending on how the data is put within that slice (e.g. NOT a pointer) periodic destruction and re-creation of the slice can allow garbage collection.

```go
type GarbageCollecter interface {
    GarbageCollect()
}
```

Dequeuer can be used to destructively remove one or more items from the queue, underflow will be true if the queue is empty. In the event the queue is empty, the output of items and flush will have a length of zero. Once an item is removed, it loses its place in the fifo and it's order can't be guaranteed.

```go
type Dequeuer interface {
    Dequeue() (item interface{}, underflow bool)
    DequeueMultiple(n int) (items []interface{})
    Flush() (items []interface{})
}
```

Peeker can be used to non-destructively remove one or more items from the queue. Underflow is true if there are no items in the queue.

```go
type Peeker interface {
    Peek() (items []interface{})
    PeekHead() (item interface{}, underflow bool)
    PeekFromHead(n int) (items []interface{})
}
```

Enqueuer can be used to put one or more item in the queue, overflow is true if the queue is full.

```go
type Enqueuer interface {
    Enqueue(item interface{}) (overflow bool)
    EnqueueMultiple(items []interface{}) (itemsRemaining []interface{}, overflow bool)
}
```

EnqueueInFronter can be used to place a single item at teh front of the queue, if the queue is full overflow will be true. Note that this won't "add" an item to the queue if its full.

```go
type EnqueueInFronter interface {
    EnqueueInFront(item interface{}) (overflow bool)
}
```

Info can be used to return information about the queue such as how many items are in the queue, or the current "size" of the queue.

```go
type Info interface {
    Length() (size int)
    Capacity() (capacity int)
}
```

Event can be used to get a read-only signal channel that will signal with an empty struct whenever data is put "in" to the queue or taken "out" of the queue. These are very useful in avoiding polling in certain patterns.

```go
type Event interface {
    GetSignalIn() (signal <-chan struct{})
    GetSignalOut() (signal <-chan struct{})
}
```

## Patterns

These are a handful of patterns that can be used to get data out of and into the queue using the given interfaces. Almost all of these patterns are based on the producer/consumer design patterns and variants of it.

All of these patterns assume that the queue is of a fixed size. Some of them don't make sense for infinite queues.

This is a producer "polling" pattern, it will enqueue data at the rate of the producer ticker. The "in" is fairly straight forward, keep in mind that you don't have to perform any conversion for the data in since it's an empty interface. Just be careful about using non-scalar values, I think it's a good practice to keep items in the queue 1:1.

Pros:

- Immediate feedback if the queue is full (via overflow)

Cons:

- There's no type safety for enqueing data (be careful)

```go
var queue goqueue.Enquerer

tProduce := time.NewTicker(time.Second)
defer tProduce.Stop()
for {
    select {
    case <-tProduce.C:
        tNow := time.Now()
        if overflow := queue.Enqueue(tNow); !overflow {
            fmt.Printf("enqueued: %v\n", tNow)
        }
    }
}
```

This is a polling producer pattern that handles situations where the queue could be full meaning that the data in the queue is being produced faster than it can be consumed.

Pros:

- This ensures that even if data is being consumed slower than it's being produced, you don't lose any data (but you can't produce as fast as you can consume...)

Cons:

- Because this uses "polling", it can only check as fast as the ticker, so you could hypothetically sacrifice CPU cycles for data integrity.

```go
var queue goqueue.Enquerer

tProduce := time.NewTicker(time.Second)
defer tProduce.Stop()
<-start
for {
    select {
    case <-tProduce.C:
        tNow := time.Now()
        for overflow := queue.Enqueue(tNow); !overflow; {
            fmt.Println("overflow occured")
            <-time.After(time.Millisecond)
            overflow = queue.Enqueue(tNow)
        }
    case <-stopper:
        return
    }
}
```

Alternatively, this is an event-based producer pattern that handles situations where the queue could be full and is a little more efficient in terms of cpu usage; just keep in mind that if there are multiple producers, there's no guarantee that once you get the signal the queue won't be full.

Cons:

- This means that you're producing faster than you can consume, this only makes sense in a go routine, but it generally means that you should increase the size of your queue
- Has the potential to block forever, make sure that you have some way to stop it (e.g. a stopper signal channel)

```go
var queue interface{
    goqueue.Enquerer
    goqueue.Event
}

signal := queue.GetSignalIn()
tProduce := time.NewTicker(time.Second)
defer tProduce.Stop()
<-start
for {
    select {
    case <-tProduce.C:
        tNow := time.Now()
        for overflow := queue.Enqueue(tNow); !overflow; {
            fmt.Println("overflow occured")
            <-signal
            overflow = queue.Enqueue(tNow)
        }
    case <-stopper:
        return
    }
}
```

This is a consumer polling pattern, it will dequeue data at the rate of the consumer ticker. This out is rather annoying in that it outputs an empty interface, and you need to know how to cast that into the appropriate data type. Type switch case is the most elegant solution when you have more than one data types.

Be careful to NOT use anonymous structs that travel between package boundaries (they aren't always equivalent).

Although this works, this has the down-side that you're limited at being able to consume the data no faster than 1Hz. If data is produced any faster, it could fall significantly behind. If data is produced any slower, you waste cpu cycles with an underflow. Its simply not super efficient if the producer and consumer aren't at approximately the same speed.

```go
var queue goqueue.Dequeuer

tConsume := time.NewTicker(time.Second)
defer tConsume.Stop()
for {
    select {
    case <-tConsume.C:
        if item, underflow := queue.Dequeue(); !underflow {
            switch v := item.(type) {
                default:
                    fmt.Printf("unsupported type: %T\n", v)
                case time.Time, *time.Time:
                    fmt.Printf("dequeued: %v\n", v)
            }
        }
    }
}
```

This is a consumer event-based pattern, where it will wait until data is placed INTO the queue, each time data is placed into the queue, a signal is received, which will then "trigger" the logic to dequeue. This works really well and lets the loop run as fast as data is being produced.

This works really well, but it has the down-side that it could miss a signal if data is produced faster than it can be consumed.

```go
var queue interface{
    goqueue.Dequeuer
    goqueue.Event
}

signal := queue.GetSignalIn()
tConsume := time.NewTicker(time.Second)
defer tConsume.Stop()
for {
    select {
      case <-signal:
        if item, underflow := queue.Dequeue(); !underflow {
            switch v := item.(type) {
                default:
                    fmt.Printf("unsupported type: %T\n", v)
                case time.Time, *time.Time:
                    fmt.Printf("dequeued: %v\n", v)
            }
        }
    }
}
```

This is a high throughput design pattern; it handles the "scalar" problem by being able to consume significantly more data per cycle; it also handles the "I missed a signal" problem, by ALSO using polling to consume the data. The Flush() function allows you to process all available data at once.

This pattern ensures that no data is lost and that you consume data faster or as fast as you produce it.

```go
var queue interface{
    goqueue.Dequeuer
    goqueue.Event
}

signal := queue.GetSignalIn()
consumeFx := func(items []interface{}) {
    for _, item := range items {
        switch v := item.(type) {
            default:
                fmt.Printf("unsupported type: %T\n", v)
            case time.Time, *time.Time:
                fmt.Printf("dequeued: %v\n", v)
        }
    }
}
tConsume := time.NewTicker(time.Second)
defer tConsume.Stop()
for {
    select {
      case <-tConsume.C:
        items := queue.Flush()
        consumeFx(items)
      case <-signal:
        items := queue.Flush()
        consumeFx(items)
    }
}
```

## Testing

The existing tests are implemented as "code" and can be used within your implementation's tests to "confirm" that they implement the interfaces as expected by "this" version of the go-queue package.

Take note that there isn't an enqueue test, this is because that's fairly specific to the implementation.

These are the avaialble unit tests:

- New: can be used to verify the constructor
- GarbageCollect: can be used to verify garbage collection
- Dequeue: can be used to verify dequeue
- DequeueMultiple: can be used to verify dequeue multiple
- Flush: can be used to verify flush
- Peek: can be used to verify peek
- PeekFromHead: can be used to verify peek from head

These are the available function/integration tests:

- Event: can be used to verify that event signals work as expected
- Info: can be used to verify that info works as expected (finite leaning)
- Queue: can be used to verify that queue works as expected (in general)
- Async: can be used to verify if safe for concurrent usage

To use one of the tests, you can use the following code snippet. Keep in mind that in order to test, the queue/constructor needs to implement ALL of the interfaces expected by the test (and by association they need to implement those interfaces as expected).

```go
import (
    "testing"

    goqueue "github.com/antonio-alexander/go-queue"
    finite "github.com/antonio-alexander/go-queue/finite"

    goqueue_tests "github.com/antonio-alexander/go-queue/tests"
)

func TestQueue(t *testing.T) {
    goqueue_tests.Queue(t, func(size int) interface {
        goqueue.Owner
        goqueue.Enqueuer
        goqueue.Dequeuer
        goqueue.Info
    } {
        return finite.New(size)
    })
}
```

## Finite Queue

This is a fixed size fifo, when the queue is full, it won't allow you to place any more items inside the queue (sans EnqueueLossy). For more information, look at this [README.md](./finite/README.md).

## Infinite Queue

This is a queue that starts with a fixed size, but when that queue fills up, it'll grow by the initially configured grow size. For more information, look at this [README.md](./infinite/README.md).
