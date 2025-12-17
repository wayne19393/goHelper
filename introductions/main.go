package main

import "fmt"

// global variables
var (
	firstName string = "Foo"
	lastName  string = "Foo"
)

// general types
var (
	floatVar32 float32 = 0.1
	floatVar64 float64 = 0.1 //you always use float64
	name       string  = "Foo"
	intVar32   int32   = 1
	intVar64   int64   = 12345
	intVar     int     = 10 //always use this basically
	uintVar    uint    = 1
	uintVar32  uint32  = 1
	uintVar64  uint64  = 1
	uintVar8   uint8   = 0x2
	byteVar    byte    = 0x1
	runeVar    rune    = 'a' //represents unicode characters
)

type Player struct {
	//used to create a collection of members of different data types, into a single variable
	name        string
	attackPower float64
	health      int
}

//func getHealth(player Player) int {
//	// always set name of your argument inside function(player) then the type (Player)
//	return player.health //you can attach part of your struct to a function and we call that a function reciever which we will show in the lower section
//}

func (player Player) getHealth() int {
	//reciever function for the Player struct
	return player.health
}

// custom types
type Weapon string

func getWeapon(weapon Weapon) string {
	return string(weapon) //you cant just set weapon cause you are saying the return type of this fun is string although
	// Weapon is a string still cause golang is strongly typed you should convert to type string
}
func main() {
	version := 1 //infer int
	player := Player{
		name:        "Captain kia",
		attackPower: 100,
		health:      200,
	} //usage of a struct
	fmt.Println(version)
	fmt.Printf("this is a player: %+v\n", player)
	//fmt.Printf("this is a player health: %+v\n", getHealth(player))
	fmt.Printf("Health: %+v\n", player.getHealth()) //because getHealth function unlike the previous implementation has access to Player struct we dont need to give it an argument
	//function recievers are commonly used and are really strong why?
	//Improves Code Organization: By attaching functions to specific types, your code becomes more organized and readable.
	//2. Encapsulation: Methods with receivers allow you to encapsulate logic and keep related functions with the type they operate on
	users := map[string]int{} //empty map
	users["kia"] = 10
	users["bar"] = 11
	notEmptyMap := map[string]int{
		"kia": 10,
	}
	ageKia := notEmptyMap["kia"]
	age := users["bar"]
	//how to check if a key exists in the map?
	ageExist, ok := users["bar"]
	if !ok {
		fmt.Println("baz not exist in map")
	} else {
		fmt.Println("exist in map:", ageExist)
	}
	//delete from map
	delete(notEmptyMap, "kia")
	//range over maps
	for k, v := range users {
		fmt.Printf("the key:%s and the value:%d\n", k, v)
	}

	//slices
	numbers := []int{1, 2, 3}
	//looping over slices are like maps
	//append
	numbers = append(numbers, 10)
	//arrays
	arrayNumbers := [2]int{1, 2}

	fmt.Printf("notEmptyMap: %+v\n", notEmptyMap)
	fmt.Printf("users: %+v\n", users)
	fmt.Printf("age: %d\n", age)
	fmt.Printf("ageKia: %d\n", ageKia)
	fmt.Printf("numbers: %+v\n", numbers)
	fmt.Printf("arrayNumbers: %+v\n", arrayNumbers)

}
