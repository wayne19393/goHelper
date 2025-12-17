package main

import "fmt"

type Position struct {
	x, y float64
}
type Entity struct {
	name    string
	id      string
	version string
	Position
}

// now if you want to make another Entity that is SpectialEntity
// we want all Entity plus extra in SpecialEntity
// we will have duplicated data how to solve it?
// we use struct embeding
// we delete everything in our SpecialEntity except the specialFields
// then we embed the Entity in our new struct
type SpecialEntity struct {
	Entity       //embeding the Entity type
	specialField float64
}

func main() {
	e := SpecialEntity{
		specialField: 88.88,
		// be aware when you have embeding types you should initialize the arguments this way
		Entity: Entity{
			name:    "MyEntity",
			id:      "id 1",
			version: "version 1.1",
			Position: Position{
				x: 220.0,
				y: 110.0,
			},
		},
	}
	// you can access them directly this way too
	e.id = "id 2"
	e.x = 33.0
	e.version = "version 1.2"
	fmt.Printf("%+v\n", e)
}
