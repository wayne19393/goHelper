package main

// these make your programs more complex
// avoid if you can
// cases are generics are useful
// you want to make a custom map
type CustomMap[K comparable, V any] struct {
	data map[K]V
}

func (m *CustomMap[K, V]) Insert(k K, v V) error {
	m.data[k] = v
	return nil
}

func NewCustomMap[K comparable, V any]() *CustomMap[K, V] {
	return &CustomMap[K, V]{
		data: make(map[K]V),
	}
}

func main() {
	m1 := NewCustomMap[string, int]()
	m1.Insert("Foo", 1)
	m1.Insert("Bar", 2)
}
