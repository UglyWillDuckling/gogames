package main

import "fmt"

func insertionSortTest(items []int) []int {
	for i := range items {
		for j := i; j > 0 && items[j-1] > items[j]; j-- {
			v := items[j]
			items[j] = items[j-1]
			items[j-1] = v
		}
	}
	return items
}

func main() {
	// var intSlice = make([]int, 10)        // when length and capacity is same
	var intSlice = []int{14, 22, 76, 88, 444, 43, 655, 43, 21, 999, 66}

	fmt.Println(insertionSortTest(intSlice))
}
