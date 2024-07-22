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

func ListBranches() (current Branch, other []Branch, err error) {
	command := exec.Command("git", "branch")
	outputBytes, err := command.Output()
	if err != nil {
		return Branch{}, nil, fmt.Errorf("error when running git branch: %s", err.Error())
	}

	lines := strings.Split(string(outputBytes), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "*") {
			lineWithoutAsterisk := strings.Replace(line, "* ", "", 1)
			current = Branch{
				Name:      strings.TrimSpace(lineWithoutAsterisk),
				IsCurrent: true,
			}
			continue
		}

		other = append(other, Branch{Name: strings.TrimSpace(line)})
	}

	return
}

func Switch(branch Branch) error {
	command := exec.Command("git", "switch", branch.Name)
	outputBytes, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error when running git branch.\n%s", string(outputBytes))
	}

	return nil
}
