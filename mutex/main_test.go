package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestMutex(t *testing.T) {
	state := State{}
	wg := sync.WaitGroup{} //wait group
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			state.count = i + 1
			state.setState(i + 1)
			wg.Done()
		}(i)
	}
	wg.Wait() //what it does it waits until those 10 times doing the work it waits for the wg.Done until its back to 0
	//as you see every time the result will be difference the reason is its concurrent you have faced your first race condition
	// --race will show you the data race
	// multiple people try to write to the same value which is the count of state
	// well that can make problems
	// so we need to find a way to synchronize
	//we do that with sync mutex you will see it in the main part
	fmt.Printf("%+v/n", state)
}
