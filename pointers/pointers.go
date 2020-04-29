package main


import "fmt"


func addOne(num *int) {
	*num = *num + 1
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

func whereIsBadGuy(guy *badGuy) {
	x := guy.pos.x
	y := guy.pos.y

	fmt.Println("(", x, ",", y, ")")
}


func main()  {

	// x := 5
	// fmt.Println(x)

	// var xPtr *int = &x

	// fmt.Println(xPtr, x)
	// fmt.Println(*xPtr, x)


	// xPtr = 665


	p := position{44, 52}
	b := badGuy{"Jabba the Hut", 100, p}


	whereIsBadGuy(&b)

	// addOne(xPtr)
	// fmt.Println(x)
}