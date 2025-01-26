package main

import (
	"fmt"
	"log"
	"os"
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

const HELP_TEXT = `
git-switch-interactive is an intuitive CLI tool that makes switching between
and deleting branches a breeze!
  
* Switch between branches
Use up/down arrow buttons (or k/j for vimmers) to move the cursor between branches.
Switch to the branch using ENTER.

* Delete multiple branches
Use d to select the branch with the cursor for deletion.
Delete all selected branches using ENTER.

`

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--help" {
		fmt.Print(HELP_TEXT)
		return
	}

	if len(os.Args) != 1 {
		fmt.Println("Unsupported arguments. Supported arguments: \"--help\"")
		return
	}

	currentBranch, otherBranches, err := git.ListBranches()
	if err != nil {
		log.Fatal(err.Error())
	}

	branchMenu := branchmenu.New(currentBranch, otherBranches)

	drawBranches(branchMenu, false)

	for {
		input, err := readInput()
		if err != nil {
			log.Fatal(err.Error())
		}

		// exit
		if keycode.Matches(input, keycode.ESCAPE, keycode.CONTROL_C, keycode.Q) {
			break
		}

		// move cursor up
		if keycode.Matches(input, keycode.UP, keycode.K) {
			branchMenu.CursorUp()
			drawBranches(branchMenu, true)
		}

		// move cursor down
		if keycode.Matches(input, keycode.DOWN, keycode.J) {
			branchMenu.CursorDown()
			drawBranches(branchMenu, true)
		}

		// switch branch
		if !branchMenu.HasBranchesSelectedForDelete() && keycode.Matches(input, keycode.ENTER) {
			err = git.Switch(branchMenu.SelectedBranch())
			if err != nil {
				log.Fatal(err.Error())
			}
			break
		}

		// select/unselect branch for deletion
		if keycode.Matches(input, keycode.D) {
			branchMenu.ToggleCurrentForDelete()
			drawBranches(branchMenu, true)
		}

		// delete selected branches
		if branchMenu.HasBranchesSelectedForDelete() && keycode.Matches(input, keycode.ENTER) {
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

func drawBranches(branchMenu branchmenu.BranchMenu, redraw bool) {
	if redraw {
		// Move the cursor up n lines where n is the number of options, setting the new
		// location to start printing from, effectively redrawing the option list
		//
		// This is done by sending a VT100 escape code to the terminal
		// @see http://www.climagic.org/mirrors/VT100_Escape_Codes.html
		// ref: https://medium.com/@nexidian/writing-an-interactive-cli-menu-in-golang-d6438b175fb6
		fmt.Printf("\033[%dA", branchMenu.BranchCount()+1)
	}

	if branchMenu.HasBranchesSelectedForDelete() {
		fmt.Println("press ENTER to delete the selected branches, press \"d\" to select/unselect a branch for deletion.")
	} else {
		fmt.Println("press ENTER to switch to the selected branch, press \"d\" to select a branch for deletion.")
	}

	fmt.Printf("%s   %s (current)%s\n", COLOUR_CURRENT_BRANCH, branchMenu.Current.Name, COLOUR_RESET)

	for _, branch := range branchMenu.Others {
		line := ""

		if branch == branchMenu.SelectedBranch() {
			line += ">"
		} else {
			line += " "
		}

		if slices.Contains(branchMenu.SelectedForDelete, branch) {
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
	if err != nil {
		log.Fatal(err)
	}

	err = terminal.Restore()
	if err != nil {
		log.Fatal(err)
	}
	terminal.Close()

	return readBytes, nil
}
