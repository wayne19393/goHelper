package main

type Storage interface {
	//now in golang the best practice for an interface is
	//to be as small as possible
	//but as you can see an interface can become a really big one
	Get(id int) (any, error)
	Put(id int, data any) error
}
type Server struct {
	store Storage
}
type FooStorage struct{}

func (s Server) Get(id int) (any, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) Put(id int, data any) error {
	//TODO implement me
	panic("implement me")
}
func (f *FooStorage) Get(id int) (any, error) {
	return nil, nil
}
func (f *FooStorage) Put(id int, data any) error {
	return nil
}
func updateValue(id int, value any) error {
	store := &FooStorage{} //hard dependency of FooStorage
	//unless you ise the power interfaces
	//dependency injection
	return store.Put(id, value)
}

// with this new storage Storage you have a dependancy injection
func updateValue2(id int, value any, storage Storage) error {
	store := &FooStorage{}
	return store.Put(id, value)
}

// you can make this even better storage in the above function
// is a really big dependency
// imagine that the Storage interface had 20 implementation for it then what whould you do?
// you create putter

func main() {
	s := &Server{
		store: &FooStorage{},
	}
	updateValue2(1, "Foo", s)
	s.store.Get(1)
	s.store.Put(1, "hello")
}
