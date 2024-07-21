package main

import (
	"fmt"
	"log"
	"os/exec"
)


func main() {
	fmt.Print("hello")
	command := exec.Command("git branch")

	output, err := command.Output()
	if err != nil {
		log.Fatalf("error when running git branch: %s", err.Error())
	}

	log.Print(output)
}
