package main

import (
	"fmt"
	"time"
)

type Server struct {
	quitch chan struct{} // make the channel 0 bytes channels that only need to communicate binary statements
	msgch  chan string
}

func NewServer() *Server {
	return &Server{
		quitch: make(chan struct{}),
		msgch:  make(chan string, 128),
	}
}

func (s *Server) start() {
	fmt.Println("starting server")
	s.loop() //block cause we have a for loop that runs forever so initialize your for loop
}

func (s *Server) stop() {
	fmt.Println("stopping server")
}
func (s *Server) sendMessage(msg string) {
	s.msgch <- msg
}

func (s *Server) quit() {
	s.quitch <- struct{}{}
}
func (s *Server) loop() {
mainloop:
	for {
		select {
		case <-s.quitch:
			//do someStuff
			fmt.Println("quiting server")
			break mainloop //if you don't name the loop you will only break the select statement and not the for loop
		case msg := <-s.msgch:
			// do someStuff when we have a message
			s.handleMessage(msg)
		default: //since you are handling the loop break correctly you can even delete the default part
		}
	}
	fmt.Println("stopping server gracefully")
}
func (s *Server) handleMessage(msg string) {
	fmt.Println("received message:", msg)
}
func main() {
	server := NewServer()
	server.start() //start is blocking when there is no default in loop function cause it knows it doesnt do anything
	go func() {
		time.Sleep(5 * time.Second)
		server.quit()
	}()

	//	for i := 0; i < 10; i++ {
	//		server.sendMessage(fmt.Sprintf("handling thess numbers: %d", i))
	//	}
	//	server.sendMessage("sending a message")
	//	time.Sleep(5 * time.Second)
}
