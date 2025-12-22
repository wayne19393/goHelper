package main

//remember if you dont change the test case each time you run again it reads the result from its cache
//with -count=1 you can tell go to not cache the results

import (
	"fmt"
	"reflect"
	"testing"
)

func TestPlayer(t *testing.T) {
	expected := Player{
		Name: "Foo",
		Hp:   100,
	}
	have := Player{
		Name: "Bar",
		Hp:   95,
	}
	//this way you tell it to compare deeply according to the type which is Player
	if !reflect.DeepEqual(expected, have) {
		t.Errorf("Expected: %+v, Have: %+v", expected, have)
	}
}

func TestCalculateValues(t *testing.T) {
	var (
		expected = 10
		a        = 5
		b        = 5
	)
	have := calculateValues(a, b)
	if have != expected {
		t.Errorf("have %d, want %d", have, expected)
	}
	fmt.Println("hello from test!")
}
