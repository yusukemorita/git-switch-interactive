package main

import (
	"fmt"
	"log"

	"example.com/git-switch-interactive/internal/git"
	"github.com/pkg/term"
)

var (
	up        []byte = []byte{27, 91, 65}
	down      []byte = []byte{27, 91, 66}
	enter     []byte = []byte{13, 0, 0}
	escape    []byte = []byte{27, 0, 0}
	control_c []byte = []byte{3, 0, 0}
)

func main() {
	branches, err := git.ListBranches()
	if err != nil {
		log.Fatalf(err.Error())
	}

	currentIndex := 0
	drawBranches(branches, currentIndex, false)

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
				currentIndex += len(branches)
			}
			drawBranches(branches, currentIndex, true)
		}

		if keycodeMatches(keycode, down) {
			currentIndex = (currentIndex + 1) % len(branches)
			drawBranches(branches, currentIndex, true)
		}

		if keycodeMatches(keycode, enter) {
			err = git.Switch(branches[currentIndex])
			if err != nil {
				log.Fatal(err.Error())
			}
			break
		}
	}
}

func drawBranches(branches []git.Branch, currentIndex int, redraw bool) {
	if redraw {
		// Move the cursor up n lines where n is the number of options, setting the new
		// location to start printing from, effectively redrawing the option list
		//
		// This is done by sending a VT100 escape code to the terminal
		// @see http://www.climagic.org/mirrors/VT100_Escape_Codes.html
		fmt.Printf("\033[%dA", len(branches))
	}

	for index, branch := range branches {
		if index == currentIndex {
			fmt.Printf("> %s\n", branch.Name)
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
