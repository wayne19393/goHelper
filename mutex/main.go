package main

import "sync"

// state -> set OR update OR delete
type State struct {
	mu    sync.Mutex
	count int
}

func (s *State) setState(i int) {
	s.mu.Lock()
	s.count = i
	s.mu.Unlock() //free up the mutex for the next in the line
}

func main() {
}
