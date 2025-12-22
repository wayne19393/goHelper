package main

import (
	"fmt"
	"time"
)

//there is no parallel concept in go routines

func main() {
	//results := fetchResources()
	//fmt.Println(results)

	//think you have some computations like a call to a library or some resource
	//you basically do the under one
	fetchResources()
	fetchResources()
	fetchResources()
	fetchResources()

	//what happens it takes too long cause you fetch them one by one that will take 2 secs
	//here you can put them in a goroutine to be ran conccurant
	//you tell go fetchResource is going to be scheduled and its async call
	go fetchResources(1)
	//now with the above code you can use the alternative syntax
	go func() {
		result := fetchResources(1)
		fmt.Println("result:", result)
	}()
	//now the problem is you wont get result cause as we said goroutin is async meaning although
	//you went inside that function it doesnt wait for the func to return anything
	//it just schedules the result to get whenever its finished and then it continues
	//since there are no other logics here the code finishes before it gets the result
	//whats the solution?
	// well the solution is channels
	//always append ch at the end of the variables for channel
	//you can make a channel of any fucking type you want
	resultch := make(chan string)      // -> unbuffered channel
	resultch1 := make(chan string, 10) // buffered channel
	resultch <- "foo"                  //writing into channel
	result := <-resultch               //reading from the resultch into result
	println(result)
	//this will result into a deadlock the reason is that you used unbuffered channel
	//a channel in golang will always block if its full
	resultch1 <- "bar"
	result1 := <-resultch1
	println(result1)
	//now this one wont block it cause there is space left in channel resultch1 to insert bar into
	//if you want to make it work in unbuffered channel
	//you can create a cookie box give it to that box then accept new cookie
	go func() {
		result := <-resultch1
		fmt.Println("result:", result)
	}()
	resultch <- "foo"
	// now always remmember if youre producer is faster than your consumer you will face the problem of deadlock so buffer it bigger when this scenario is true
	msgch := make(chan string, 128)
	msgch <- "A"
	msgch <- "B"
	msgch <- "C"
	msgch <- "D"
	close(msgch) // why we close? just read on and you'll understand
	// as you remember for one producer we had one reader but this is weird in multiple msgch
	msg := <-msgch
	println(msg)
	// we should range over channels
	for msg := range msgch {
		fmt.Println("the message is:", msg)
	}
	fmt.Println("done reading messages from the channel")
	// again you'll encounter a deadlock here
	// because we are writing no logic to stop the range so our consumer don't know we stopped producing messages?
	// so you should close the msg channel

	// now if you use only a for loop and not with range if you close before the for the for will loop for ever and doesnt have any break mechanism
	// you should have ok which is a boolean to notify the break in that case
	for {
		msg, ok := <-msgch
		if !ok {
			break
		}
		fmt.Println("the message is:", msg)
	}

	// our next concern is control flow
	// we will use a pattern check goldenPatern directory

}

// when you want to use goroutine its better to have your fun with the n int and
func fetchResources(n int) string {
	time.Sleep(1 * time.Second)
	return fmt.Sprintf("result %d", n)
}
