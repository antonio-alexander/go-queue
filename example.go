package goqueue

import (
	"encoding/json"
	"fmt"
	"math/rand"
)

type Example struct {
	Int    int     `json:"int,omitempty"`
	Float  float64 `json:"float,omitempty"`
	String string  `json:"string,omitempty"`
}

func (v *Example) MarshalBinary() ([]byte, error) {
	return json.Marshal(v)
}

func (v *Example) UnmarshalBinary(bytes []byte) error {
	return json.Unmarshal(bytes, v)
}

func ExampleConvertSingle(item interface{}) *Example {
	switch v := item.(type) {
	default:
		fmt.Printf("Arf! Arf! unsupported type: %T", v)
		return nil
	case *Example:
		return v
	case Bytes:
		data := &Example{}
		if err := data.UnmarshalBinary(v); err != nil {
			fmt.Printf("Arf! Arf! failure unmarshalling data: %s", err)
			return nil
		}
		return data
	case []byte:
		data := &Example{}
		if err := data.UnmarshalBinary(v); err != nil {
			fmt.Printf("Arf! Arf! failure unmarshalling data: %s", err)
			return nil
		}
		return data
	}
}

func ExampleConvertMultiple(items []interface{}) []*Example {
	values := make([]*Example, 0, len(items))
	for _, item := range items {
		value := ExampleConvertSingle(item)
		if value == nil {
			continue
		}
		values = append(values, value)
	}
	return values
}

func ExampleClose(queue Owner) []*Example {
	items := queue.Close()
	return ExampleConvertMultiple(items)
}

func ExampleEnqueueMultiple(queue Enqueuer, values []*Example) ([]*Example, bool) {
	items := make([]interface{}, 0, len(values))
	for _, value := range values {
		items = append(items, value)
	}
	itemsRemaining, overflow := queue.EnqueueMultiple(items)
	if !overflow {
		return nil, false
	}
	return ExampleConvertMultiple(itemsRemaining), true
}

func ExamplePeek(queue Peeker) []*Example {
	items := queue.Peek()
	return ExampleConvertMultiple(items)
}

func ExamplePeekHead(queue Peeker) (*Example, bool) {
	item, underflow := queue.PeekHead()
	if underflow {
		return nil, true
	}
	return ExampleConvertSingle(item), false
}

func ExamplePeekFromHead(queue Peeker, n int) []*Example {
	items := queue.PeekFromHead(n)
	return ExampleConvertMultiple(items)
}

func ExampleDequeue(queue Dequeuer) (*Example, bool) {
	item, underflow := queue.Dequeue()
	if underflow {
		return nil, true
	}
	return ExampleConvertSingle(item), false
}

func ExampleDequeueMultiple(queue Dequeuer, n int) []*Example {
	items := queue.DequeueMultiple(n)
	return ExampleConvertMultiple(items)
}

func ExampleFlush(queue Dequeuer) []*Example {
	items := queue.Flush()
	return ExampleConvertMultiple(items)
}

//ExampleGenFloat64 will generate a random number of random float values if n is equal to 0
// not to exceed the constant TestMaxExamples, if n is provided, it will generate that many items
func ExampleGenFloat64(n int) []*Example {
	if n <= 0 {
		n = int(rand.Float64() * 1000)
	}
	values := make([]*Example, 0, n)
	for i := 0; i < n; i++ {
		values = append(values, &Example{Float: rand.Float64()})
	}
	return values
}
