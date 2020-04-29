package main

// NPCs, talk to them, fight, move around the graph
// Items that can be picked up and used
// build your own games

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Player struct
// the player cannot attack, only engage in conversation
type Player struct {
	health   int
	nickname string
}

// bad guy can talk to the player
// after the dialog is finished, the bad guy will decide if it will attack the character
type badguy struct {
	desc        string
	health      int
	attackPower int
	story       *storyNode
}

func (guy *badguy) interact(player *Player) {
	// for starters, the monster will just start the conversation
	guy.talk(player)
}

func (guy *badguy) talk(player *Player) {
	guy.story.play()
}

// attacks the main player
func (guy *badguy) Attack(player Player) {
	fmt.Println("Grue is attacking you")

	player.health = player.health - guy.attackPower

	if player.health < 0 {
		player.health = 0
	}

	fmt.Println("Player health", player.health)
	fmt.Println("The monster leaves", player.health)
}

type storyNode struct {
	text string
	// as many nodes or no nodes
	choices *choices
	monster *badguy
	event   func(Player)
}

type choices struct {
	cmd         string
	description string
	nextNode    *storyNode
	nextChoice  *choices
}

func (node *storyNode) addChoice(cmd string, description string, nextNode *storyNode) {
	choice := &choices{cmd, description, nextNode, nil}

	if node.choices == nil {
		node.choices = choice
	} else {
		currentChoice := node.choices

		for currentChoice.nextChoice != nil {
			currentChoice = currentChoice.nextChoice
		}
		currentChoice.nextChoice = choice
	}
}

func (node *storyNode) render() {
	fmt.Println(node.text)
	currentChoice := node.choices

	for currentChoice != nil {
		fmt.Println(currentChoice.cmd, ":", currentChoice.description)
		currentChoice = currentChoice.nextChoice
	}
}

func (node *storyNode) executeCmd(cmd string) *storyNode {
	currentChoice := node.choices

	for currentChoice.nextChoice != nil {
		if strings.TrimSpace(cmd) == strings.TrimSpace(currentChoice.cmd) {
			return currentChoice.nextNode
		}
		currentChoice = currentChoice.nextChoice
	}

	if strings.TrimSpace(cmd) == strings.TrimSpace(currentChoice.cmd) {
		return currentChoice.nextNode
	}

	fmt.Println("Sorry, I didn't understand that.")
	return node
}

var scanner *bufio.Scanner

func (node *storyNode) play() {
	if node.monster != nil {
		node.monster.story.play()
	}

	node.render()

	if node.event != nil {
		node.event(player1)
	}

	if node.choices != nil {
		scanner.Scan()
		node.executeCmd(scanner.Text()).play()
	}
}

func buildNewGrue() *badguy {
	grue := &badguy{health: 120, desc: "this is a grue monster", attackPower: 40, story: nil}

	buildGrueStory(grue)

	return grue
}

// attaches a new story node to the grue
func buildGrueStory(badguy *badguy) {
	start := storyNode{text: "A monster appears infront of you. I Grue. I eat you now"}

	pretty := storyNode{text: "Hmmmmmm. Interesting, more"}             // dialog opens up
	disgust := storyNode{text: "I kill you now", event: grue.Attack}    // beast should attack at this point, we pass in the attack event
	ignore := storyNode{text: "You weird. I am strong, I wil hit you."} // dialog opens up

	start.addChoice("P", "You are pretty", &pretty)
	start.addChoice("D", "You disgust me you ugly monster", &disgust)
	start.addChoice("B", "Bow down beast, I cannot be bother by the likes of you", &ignore)

	prettyHeart := storyNode{text: "Awwww, I kiss you now"}
	prettyKidding := storyNode{text: "You bad. Dieeeee..."}
	prettyLie := storyNode{text: "That nice. I have many grue friends too."}

	pretty.addChoice("H", "Your ears are like two arrows pointing at my heart", &prettyHeart)
	pretty.addChoice("K", "Just kidding you uggly beast. You will burn in hell", &prettyKidding)
	pretty.addChoice("Lie", "I have many Grue friends just like you", &prettyLie)

	ignorePride := storyNode{text: "Atttackkkkkkk!!!!"}
	ignoreQuest := storyNode{text: "You so weird. Die now."}

	ignore.addChoice("pride", "Try, a creature like you is no match for a knight of Dorne", &ignorePride)
	ignore.addChoice("ignore", "my quest is too important, out of my way", &ignoreQuest)
	ignore.addChoice("pretty", "On second thought, you are quite lovely", &pretty)

	badguy.story = &start
}

var player1 Player
var grue badguy

func main() {
	scanner = bufio.NewScanner(os.Stdin)

	player1 = Player{nickname: "player 1", health: 100}
	firstmonster := buildNewGrue()

	start := storyNode{text: `
		You are in a large chamber, deep underground.
		You see three passages leading out. A north passage leads into darkness.
		To the south, a passage appears to head upward. The eastern passages appears
		flat and well traveled`}

	darkRoom := storyNode{text: "It is pitch black. You cannot see a thing."}
	darkRoomLit := storyNode{text: "The dark passage is now lit by your lantern. You can continue north or head back south"}

	grue := storyNode{text: "While stumbling around in the darkness, you are eaten by a grue"}
	trap := storyNode{text: "You head down the well and traveled path when suddenly a trap door opens and you fall into a pit."}
	treasure := storyNode{text: "You arrive at a samll chamber, filled with treasure!"}

	start.addChoice("N", "Go North", &darkRoom)
	start.addChoice("S", "Go South", &darkRoom)
	start.addChoice("E", "Go East", &trap)

	darkRoom.addChoice("S", "Try to go back", &grue)
	darkRoom.addChoice("O", "Turn on lantern", &darkRoomLit)

	darkRoomLit.monster = firstmonster
	darkRoomLit.addChoice("N", "Go North", &treasure)
	darkRoomLit.addChoice("S", "Go South", &start)

	start.play()

	fmt.Println()
	fmt.Println("The End.")
}
