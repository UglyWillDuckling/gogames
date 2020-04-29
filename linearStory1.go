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
	page2 := storyPage{"You are alone, and you need to find the sacred helment before the bad guys do.", nil}
	page3 := storyPage{"You see a troll ahead", nil}

	page1.nextPage = &page2
	page2.nextPage = &page3

	// playStory(&page1)

	page25 := storyPage{"I'm an insert page, add me afer page two, just before the troll comes.", nil}

	insertPage(&page25, &page2)

	deletePage(&page3)

	insertPage(&storyPage{"I'm the new last page, instead of the Troll there is a bear.", nil}, &page25)

	playStory(&page1)

	// Functions - has a return value - may also execute
	// Procedures - has no return value, just executes
	// Methods - functions attached to a struct/objet
}
