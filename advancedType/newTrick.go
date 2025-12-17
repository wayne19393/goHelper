package main

import (
	"fmt"
)

// make enum
type Color int

const (
	Black Color = iota //iota increaments
	Blue
	Green
	Red
	Yellow
)

// now cause there is the answer comes back as 0 only we implement stringer interface for Color
func (c Color) String() string {
	switch c {
	case Black:
		return "Black"
	case Blue:
		return "Blue"
	case Green:
		return "Green"
	case Red:
		return "Red"
	case Yellow:
		return "Yellow"
	default:
		panic("Invalid color")
	}
}

func main() {
	fmt.Printf("%+v\n", Blue)
}
