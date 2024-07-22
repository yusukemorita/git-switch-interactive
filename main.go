package main

import (
	"fmt"
	"log"

	"github.com/pkg/term"
	"github.com/yusukemorita/git-switch-interactive/internal/git"
	"github.com/yusukemorita/git-switch-interactive/internal/keycode"
)

var (
	// colours
	colourReset          = "\033[0m"
	selectedBranchColour = "\033[34m" // blue
	currentBranchColour  = "\033[32m" // green
)

func main() {
	currentBranch, otherBranches, err := git.ListBranches()
	if err != nil {
		log.Fatalf(err.Error())
	}

	currentIndex := 0
	drawBranches(currentBranch, otherBranches, currentIndex, false)

	for {
		input, err := readInput()
		if err != nil {
			log.Fatal(err.Error())
		}

		if keycode.Matches(input, keycode.ESCAPE, keycode.CONTROL_C) {
			break
		}

		if keycode.Matches(input, keycode.UP, keycode.K) {
			currentIndex = currentIndex - 1
			if currentIndex < 0 {
				currentIndex += len(otherBranches)
			}
			drawBranches(currentBranch, otherBranches, currentIndex, true)
		}

		if keycode.Matches(input, keycode.DOWN, keycode.J) {
			currentIndex = (currentIndex + 1) % len(otherBranches)
			drawBranches(currentBranch, otherBranches, currentIndex, true)
		}

		if keycode.Matches(input, keycode.ENTER) {
			err = git.Switch(otherBranches[currentIndex])
			if err != nil {
				log.Fatal(err.Error())
			}
			break
		}
	}
}

func drawBranches(current git.Branch, otherBranches []git.Branch, currentIndex int, redraw bool) {
	if redraw {
		// Move the cursor up n lines where n is the number of options, setting the new
		// location to start printing from, effectively redrawing the option list
		//
		// This is done by sending a VT100 escape code to the terminal
		// @see http://www.climagic.org/mirrors/VT100_Escape_Codes.html
		// ref: https://medium.com/@nexidian/writing-an-interactive-cli-menu-in-golang-d6438b175fb6
		fmt.Printf("\033[%dA", len(otherBranches)+1)
	}

	fmt.Printf("%s  %s (current)%s\n", currentBranchColour, current.Name, colourReset)

	for index, branch := range otherBranches {
		if index == currentIndex {
			fmt.Printf("%s> %s%s\n", selectedBranchColour, branch.Name, colourReset)
		} else {
			fmt.Printf("  %s\n", branch.Name)
		}
	}
}

func readInput() ([]byte, error) {
	terminal, err := term.Open("/dev/tty")
	if err != nil {
		log.Fatal(err.Error())
	}

	err = terminal.SetRaw()
	if err != nil {
		log.Fatal(err.Error())
	}

	readBytes := make([]byte, 3)
	_, err = terminal.Read(readBytes)

	terminal.Restore()
	terminal.Close()

	return readBytes, nil
}
