package main

// 1. make the program print out how many tries it took
// 2. See if you can tell the user is lying

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	low := 1
	high := 100

	fmt.Println("Please think of a number between", low, " and ", high)
	fmt.Println("Press enter when ready")
	scanner.Scan()

	// n == 1 to 100 O(n)

	guess := 50
	lastguess := 0
	tries := 1
	for {
		fmt.Println("I guess the number is ", guess)

		fmt.Println("Is that :")
		fmt.Println("(a) too high?")
		fmt.Println("(b) too low?")
		fmt.Println("(c) correct")

		scanner.Scan()
		response := scanner.Text()

		if response == "a" {
			high = guess
			guess = (guess + low) / 2
		} else if response == "b" {
			low = guess
			guess = (guess + high) / 2
		} else if response == "c" {
			fmt.Println("I won!", "Number of tries ", tries)
			break
		} else {
			fmt.Println("Invalid response, try again.")
		}

		if guess == lastguess {
			fmt.Println("Liar!!!")
			break
		}

		lastguess = guess
		tries++
	}
}
