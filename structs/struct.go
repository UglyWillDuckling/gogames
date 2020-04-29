package main

import "fmt"

// Game Struct
type Game struct {
}

// Position struct
type position struct {
	x float32
	y float32
}

type badGuy struct {
	name   string
	health int
	pos    position
}

func whereIsBadGuy(b badGuy) {
	x := b.pos.x
	y := b.pos.y

	fmt.Println("(", x, ",", y, ")")
}

func main() {

	// var p Position

	// p.x = 5
	// p.y = 4

	p := position{4, 2}

	// fmt.Println(p)
	// fmt.Println(p.x)
	// fmt.Println(p.y)

	b := badGuy{"Jabba the Hut", 100, p}

	whereIsBadGuy(b)
}
