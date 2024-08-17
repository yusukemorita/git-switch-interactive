package branchmenu

import (
	"slices"

	"github.com/yusukemorita/git-switch-interactive/internal/git"
)

func New(current git.Branch, others []git.Branch) BranchMenu {
	return BranchMenu{
		Current: current,
		Others:  others,
	}

}

type BranchMenu struct {
	Current           git.Branch
	Others            []git.Branch
	SelectedForDelete []git.Branch
	cursorIndex       int
}

func (menu *BranchMenu) CursorUp() {
	menu.cursorIndex = menu.cursorIndex - 1
	if menu.cursorIndex < 0 {
		menu.cursorIndex += len(menu.Others)
	}
}

func (menu *BranchMenu) CursorDown() {
	menu.cursorIndex = (menu.cursorIndex + 1) % len(menu.Others)
}

func (menu *BranchMenu) SelectedBranch() git.Branch {
	return menu.Others[menu.cursorIndex]
}

func (menu *BranchMenu) BranchCount() uint {
	// increment by 1 to count the current branch as well
	return uint(len(menu.Others) + 1)
}

func   (menu *BranchMenu) ToggleCurrentForDelete() {
	if slices.Contains(menu.SelectedForDelete, menu.SelectedBranch()) {
		var newSelectedForDelete []git.Branch
		for _, branch := range menu.SelectedForDelete {
			if branch != menu.SelectedBranch() {
				newSelectedForDelete = append(newSelectedForDelete, branch)
			}
		}
		menu.SelectedForDelete = newSelectedForDelete
	} else {
		menu.SelectedForDelete = append(menu.SelectedForDelete, menu.SelectedBranch())
	}
}

func (menu *BranchMenu) HasBranchesSelectedForDelete() bool {
	return len(menu.SelectedForDelete) > 0
}
