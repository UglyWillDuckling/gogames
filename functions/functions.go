package main

import "fmt"

func sayHello() {
	fmt.Println("Hello")
}

func addOne(x int) int {
	return x + 1
}

// recursion
func sayHelloABuch() {
	fmt.Println("Hello")
	sayHelloABuch()
}


func main() {

	x := 5

	x = addOne(addOne(x))
	// fmt.Println(x)

	x = addOne(addOne(x))

	fmt.Println(x)

	sayHelloABuch()
}
