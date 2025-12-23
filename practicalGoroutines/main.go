package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type UserProfile struct {
	ID       int
	Comments []string
	Likes    int
	Friends  []int
}

func main() {
	start := time.Now()
	userProfile, err := handleGetUserProfile(10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User Profile: %+v\n", userProfile)
	fmt.Printf("Fetching the user Profile took: %s\n", time.Since(start))
}

type Response struct {
	data any
	err  error
}

func handleGetUserProfile(id int) (*UserProfile, error) {
	var (
		responsech = make(chan Response, 3)
		wg         = &sync.WaitGroup{}
	)
	// we are doing three requests inside their own goroutine
	wg.Add(3)
	go getLikes(id, responsech, wg)
	go getFriends(id, responsech, wg)
	go getComments(id, responsech, wg)
	//add three to wait group

	wg.Wait() // block until wg counter is zero
	close(responsech)

	// keep ranging when to stop?
	userProfile := &UserProfile{}
	for resp := range responsech {
		if resp.err != nil {
			return nil, resp.err
		}
		switch msg := resp.data.(type) {
		case int:
			userProfile.Likes = msg
		case []int:
			userProfile.Friends = msg
		case []string:
			userProfile.Comments = msg
		}
		fmt.Println(resp.data)
	}
	return userProfile, nil
}
func getComments(id int, responsech chan Response, wg *sync.WaitGroup) {
	time.Sleep(200 * time.Millisecond)

	comments := []string{
		"This is a comment",
		"This is another comment",
		"Ow, I didnt know that",
	}
	responsech <- Response{
		data: comments,
		err:  nil,
	}
	wg.Done()
}
func getLikes(id int, responsech chan Response, wg *sync.WaitGroup) {
	time.Sleep(200 * time.Millisecond)
	responsech <- Response{
		data: 120,
		err:  nil,
	}
	wg.Done()
}
func getFriends(id int, responsech chan Response, wg *sync.WaitGroup) {
	time.Sleep(100 * time.Millisecond)
	friendIds := []int{11, 12, 13, 845, 534}
	responsech <- Response{
		data: friendIds,
		err:  nil,
	}
	wg.Done()
}
