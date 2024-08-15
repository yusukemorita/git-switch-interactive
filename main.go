package main

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/pkg/term"
	"github.com/yusukemorita/git-switch-interactive/internal/branchmenu"
	"github.com/yusukemorita/git-switch-interactive/internal/git"
	"github.com/yusukemorita/git-switch-interactive/internal/keycode"
)

const (
	COLOUR_RESET           = "\033[0m"
	COLOUR_SELECTED_BRANCH = "\033[34m" // blue
	COLOUR_CURRENT_BRANCH  = "\033[32m" // green
)

func main() {
	currentBranch, otherBranches, err := git.ListBranches()
	if err != nil {
		log.Fatalf(err.Error())
	}

	var isDeleteMode bool
	branchMenu := branchmenu.New(currentBranch, otherBranches)

	drawBranches(branchMenu, false, isDeleteMode)

	for {
		input, err := readInput()
		if err != nil {
			log.Fatal(err.Error())
		}

		// exit
		if keycode.Matches(input, keycode.ESCAPE, keycode.CONTROL_C) {
			break
		}

		// move cursor up
		if keycode.Matches(input, keycode.UP, keycode.K) {
			branchMenu.CursorUp()
			drawBranches(branchMenu, true, isDeleteMode)
		}

		// move cursor down
		if keycode.Matches(input, keycode.DOWN, keycode.J) {
			branchMenu.CursorDown()
			drawBranches(branchMenu, true, isDeleteMode)
		}

		// switch branch
		if !isDeleteMode && keycode.Matches(input, keycode.ENTER) {
			err = git.Switch(branchMenu.SelectedBranch())
			if err != nil {
				log.Fatal(err.Error())
			}
			break
		}

		// select branch for deletion
		if keycode.Matches(input, keycode.D) {
			isDeleteMode = true
			branchMenu.SelectCurrentForDelete()
			drawBranches(branchMenu, true, isDeleteMode)
		}

		// delete selected branches
		if isDeleteMode && keycode.Matches(input, keycode.ENTER) {
			fmt.Printf("Are you sure you want to delete the selected branches? [y/n]\n")

			input, err := readInput()
			if err != nil {
				log.Fatal(err.Error())
			}

			if keycode.Matches(input, keycode.Y) {
				var deletedBranchNames []string
				for _, branch := range branchMenu.SelectedForDelete {
					deletedBranchNames = append(deletedBranchNames, branch.Name)
					err = git.Delete(branch)
					if err != nil {
						log.Fatal(err.Error())
					}
				}

				fmt.Printf("deleted branches: %s\n", strings.Join(deletedBranchNames, ", "))
			} else {
				fmt.Println("Input does not match \"y\", ignoring")
			}
			break
		}
	}
}

func drawBranches(branchMenu branchmenu.BranchMenu, redraw bool, isDeleteMode bool) {
	if redraw {
		// Move the cursor up n lines where n is the number of options, setting the new
		// location to start printing from, effectively redrawing the option list
		//
		// This is done by sending a VT100 escape code to the terminal
		// @see http://www.climagic.org/mirrors/VT100_Escape_Codes.html
		// ref: https://medium.com/@nexidian/writing-an-interactive-cli-menu-in-golang-d6438b175fb6
		fmt.Printf("\033[%dA", branchMenu.BranchCount())
	}

	fmt.Printf("%s  %s (current)%s\n", COLOUR_CURRENT_BRANCH, branchMenu.Current.Name, COLOUR_RESET)

	for _, branch := range branchMenu.Others {
		line := ""

		if branch == branchMenu.SelectedBranch() {
			line += ">"
		} else {
			line += " "
		}

		if isDeleteMode && slices.Contains(branchMenu.SelectedForDelete, branch) {
			line += "🗑️ "
		} else {
			line += "  "
		}

		line += branch.Name

		if branch == branchMenu.SelectedBranch() {
			line = fmt.Sprintf("%s%s%s", COLOUR_SELECTED_BRANCH, line, COLOUR_RESET)
		}

		fmt.Println(line)
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
