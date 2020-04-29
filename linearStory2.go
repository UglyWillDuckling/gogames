// add a function to insert new page after a page
// delete a page from the list

package main

import (
	"fmt"
)

type storyPage struct {
	text     string
	nextPage *storyPage
}

func playStory(page *storyPage) {
	for page != nil {
		fmt.Println(page.text)
		page = page.nextPage
	}
}

func insertPage(page *storyPage, beforePage *storyPage) {
	if beforePage.nextPage != nil {
		page.nextPage = beforePage.nextPage
	}

	beforePage.nextPage = page
}

func (page *storyPage) addToEnd(text string) {
	pageToAdd := &storyPage{text, nil}

	for page.nextPage != nil {
		page = page.nextPage
	}

	page.nextPage = pageToAdd
}

func (page *storyPage) addAfter(text string) {
	newPage := &storyPage{text, page.nextPage}

	page.nextPage = newPage
}

// Delete

// delete one page from the book and keep the coherency
func deletePage(page *storyPage) {
	if page.nextPage != nil {
		*page = *page.nextPage
		return
	}

	// just set the value to the zeroth value object
	*page = storyPage{}
}

func main() {

	page1 := storyPage{"It was a dark and stormy night.", nil}

	page1.addToEnd("You are alone, and you need to find the sacred helment before the bad guys do.")
	page1.addToEnd("You see a troll ahead")

	page1.addAfter("Testing")

	playStory(&page1)

	// Functions - has a return value - may also execute
	// Procedures - has no return value, just executes
	// Methods - functions attached to a struct/object/etc
}
