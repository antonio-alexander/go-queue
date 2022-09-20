package goqueue

import (
	"encoding/json"
	"math/rand"
	"reflect"
)

//REFERENCE: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func randomString(nLetters ...int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	nLetter := 20
	if len(nLetters) > 0 {
		nLetter = nLetters[0]
	}
	b := make([]rune, nLetter)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

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
		return nil
	case *Example:
		return v
	case Bytes:
		data := &Example{}
		if err := data.UnmarshalBinary(v); err != nil {
			return nil
		}
		return data
	case []byte:
		data := &Example{}
		if err := data.UnmarshalBinary(v); err != nil {
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
	return ExampleConvertMultiple(queue.DequeueMultiple(n))
}

func ExampleFlush(queue Dequeuer) []*Example {
	return ExampleConvertMultiple(queue.Flush())
}

//ExampleGenFloat64 will generate a random number of random float values if n is equal to 0
// not to exceed the constant TestMaxExamples, if n is provided, it will generate that many items
func ExampleGenFloat64(sizes ...int) []*Example {
	size := int(rand.Float64() * 100)
	if len(sizes) > 0 {
		size = sizes[0]
	}
	values := make([]*Example, 0, size)
	for i := 0; i < size; i++ {
		values = append(values, &Example{Float: rand.Float64()})
	}
	return values
}

func ExampleGenInt(sizes ...int) []*Example {
	size := int(rand.Float64() * 100)
	if len(sizes) > 0 {
		size = sizes[0]
	}
	values := make([]*Example, 0, size)
	for i := 0; i < size; i++ {
		values = append(values, &Example{Int: rand.Int()})
	}
	return values
}

func ExampleGenString(sizes ...int) []*Example {
	size := int(rand.Float64() * 100)
	if len(sizes) > 0 {
		size = sizes[0]
	}
	values := make([]*Example, 0, size)
	for i := 0; i < size; i++ {
		values = append(values, &Example{String: randomString()})
	}
	return values
}

func ExampleGen(sizes ...int) []*Example {
	size := int(rand.Float64() * 100)
	if len(sizes) > 0 {
		size = sizes[0]
	}
	values := make([]*Example, 0, size)
	for i := 0; i < size; i++ {
		values = append(values, &Example{
			Int:    rand.Int(),
			Float:  rand.Float64(),
			String: randomString(),
		})
	}
	return values
}

func AssertExamples(example *Example, examples []*Example) func() bool {
	return func() bool {
		for _, e := range examples {
			if reflect.DeepEqual(*e, *example) {
				return true
			}
		}
		return false
	}
}
