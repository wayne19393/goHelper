package main

import (
	"fmt"
)

// Storage
// Stor(er)
// io.Read(er)
// io.Writ(er)
// In golang interfaces usually come with a -er at the end its idiom
type Numberstorer interface {
	//inside interface we describe what it does
	//example in football we have a footballer(football player) and you know he can football
	//but you dont know how he does it same with an interface
	//storer interface you know it can store but how it does it you dont know
	GetAll() ([]int, error)
	Put(int) error
}
type PostGressNumberStore struct {
	// postgress numbers (db connections)
}

func (p PostGressNumberStore) GetAll() ([]int, error) {
	//reciever function
	return []int{1, 2, 3, 4, 5, 6, 7}, nil
}
func (p PostGressNumberStore) Put(num int) error {
	fmt.Println("Store the Number into the mongoDB storage")
	return nil
}

type MongoDBNumberStore struct {
	// some values
}

func (m MongoDBNumberStore) GetAll() ([]int, error) {
	//reciever function
	return []int{1, 2, 3}, nil
}
func (m MongoDBNumberStore) Put(num int) error {
	fmt.Println("Store the Number into the mongoDB storage")
	return nil
}

type ApiServer struct {
	numberStore Numberstorer
}

func main() {
	apiServer := ApiServer{
		numberStore: PostGressNumberStore{},
		//if suddenly your CTO says we will no longer use Postgress we want to use
		// MongoDb you can just change the PostGressNumberStore and implement it like we did and change is to MongoDBNumberStore
	}
	//Logic
	err := apiServer.numberStore.Put(1000)
	if err != nil {
		panic(err)
	}
	numbers, err := apiServer.numberStore.GetAll()
	if err != nil {
		panic(err)
	}
	fmt.Println(numbers)
}
