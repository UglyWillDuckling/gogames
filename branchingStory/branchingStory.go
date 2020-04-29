package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	Info = Teal
	Warn = Yellow
	Err  = Red
)

var (
	Black   = Color("\033[1;30m%s\033[0m")
	Red     = Color("\033[1;31m%s\033[0m")
	Green   = Color("\033[1;32m%s\033[0m")
	Yellow  = Color("\033[1;33m%s\033[0m")
	Purple  = Color("\033[1;34m%s\033[0m")
	Magenta = Color("\033[1;35m%s\033[0m")
	Teal    = Color("\033[1;36m%s\033[0m")
	White   = Color("\033[1;37m%s\033[0m")
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

type storyNode struct {
	text    string
	yesPath *storyNode
	noPath  *storyNode
}

func (node *storyNode) play() {
	fmt.Println(Info(node.text))

	if node.yesPath != nil && node.noPath != nil {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			scanner.Scan()

			answer := scanner.Text()

			if answer == "yes" {
				node.yesPath.play()
				break
			} else if answer == "no" {
				node.noPath.play()
				break
			} else {
				fmt.Println(Err("that answer was not an option!"))
			}
		}
	}
}

func (node *storyNode) printStory(depth int) {
	for i := 0; i < depth; i++ {
		fmt.Print("  ")
	}
	fmt.Print(node.text)
	fmt.Println()

	if node.yesPath != nil {
		node.yesPath.printStory(depth + 1)
	}
	if node.yesPath != nil {
		node.noPath.printStory(depth + 1)
	}
	// fmt.Println((node.text))
	// node.yesPath.printStory()
	// node.noPath.printStory()
}

func main() {

	root := storyNode{"You are at the entrance to a dark cave. Do you want to go into the cave?", nil, nil}

	winning := storyNode{"You have won!", nil, nil}
	losing := storyNode{"You have lost!", nil, nil}

	root.yesPath = &losing
	root.noPath = &winning

	// root.play()

	root.printStory(0)
}
