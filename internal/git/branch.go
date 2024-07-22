package git

import (
	"fmt"
	"os/exec"
	"strings"
)

type Branch struct {
	Name      string
	IsCurrent bool
}

func ListBranches() ([]Branch, error) {
	command := exec.Command("git", "branch")
	outputBytes, err := command.Output()
	if err != nil {
		return nil, fmt.Errorf("error when running git branch: %s", err.Error())
	}

	lines := strings.Split(string(outputBytes), "\n")

	var branches []Branch

	for _, line := range lines {
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "*") {
			lineWithoutAsterisk := strings.Replace(line, "* ", "", 1)
			branches = append(branches, Branch{
				Name:      strings.TrimSpace(lineWithoutAsterisk),
				IsCurrent: true,
			})
			continue
		}

		branches = append(branches, Branch{Name: strings.TrimSpace(line)})
	}

	return branches, nil
}

func Switch(branch Branch) error {
	command := exec.Command("git", "switch", branch.Name)
	outputBytes, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error when running git branch.\n%s", string(outputBytes))
	}

	return nil
}
