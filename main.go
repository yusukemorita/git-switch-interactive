package main

import (
	"fmt"
	"log"

	"github.com/yusukemorita/git-switch-interactive/internal/git"
	"github.com/pkg/term"
)

var (
	// keycodes
	up        []byte = []byte{27, 91, 65}
	down      []byte = []byte{27, 91, 66}
	enter     []byte = []byte{13, 0, 0}
	escape    []byte = []byte{27, 0, 0}
	control_c []byte = []byte{3, 0, 0}

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
		keycode, err := readInput()
		if err != nil {
			log.Fatal(err.Error())
		}

		if keycodeMatches(keycode, escape) || keycodeMatches(keycode, control_c) {
			break
		}

		if keycodeMatches(keycode, up) {
			currentIndex = currentIndex - 1
			if currentIndex < 0 {
				currentIndex += len(otherBranches)
			}
			drawBranches(currentBranch, otherBranches, currentIndex, true)
		}

		if keycodeMatches(keycode, down) {
			currentIndex = (currentIndex + 1) % len(otherBranches)
			drawBranches(currentBranch, otherBranches, currentIndex, true)
		}

		if keycodeMatches(keycode, enter) {
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

func keycodeMatches(a, b []byte) bool {
	return a[0] == b[0] && a[1] == b[1] && a[2] == b[2]
}
