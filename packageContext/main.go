package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

// basically package context is when you want to use or call another service and want to give it a time out and handle the time your application should waint for the external service to respond
func main() {
	start := time.Now()
	ctx := context.WithValue(context.Background(), "Username", "kianush-keykhosravi") //defining a parent context so you can seperate your context in goroutines so race condition doesnt ruin your program
	userID, err := fetchUserId(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("The response took %v -> %+v\n", time.Since(start), userID)
}
func fetchUserId(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond) // everything that takes more than 100 milliseconds gets cancelled
	defer cancel()                                                //manual cancellation
	val := ctx.Value("Username")
	fmt.Println("the user is:", val)
	type result struct {
		userId string
		err    error
	}
	resultch := make(chan result, 1)
	go func() {
		res, err := thirdpartyHTTPCall()
		resultch <- result{
			userId: res,
			err:    err,
		}
	}()
	select {
	//Done()
	// 1. the context timeout is exceeded
	// 2. the context has been manually canceled -> Cancel()
	case <-ctx.Done():
		return "", ctx.Err()
	case result := <-resultch:
		return result.userId, result.err
	}
}
func thirdpartyHTTPCall() (string, error) {
	time.Sleep(10 * time.Millisecond) //if you set the response time to lesser than 100 milliSeconds you get no error but if its set to higher you get timeout exceeded
	return "some response", nil
}
