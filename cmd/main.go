package main

import (
	"fmt"
	"log"


	"example.com/git-switch-interactive/internal/git"
)


func main() {
	branches, err := git.ListBranches()
	if err != nil {
		log.Fatalf(err.Error())
	}

	for _, branch := range branches {
		if branch.IsCurrent {
			fmt.Printf("- # %s\n", branch.Name)
			continue
		}

		fmt.Printf("- %s\n", branch.Name)
	}
}
