package main

import (
	"log"

	"github.com/golang-travel/tour/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("cmd.Execute: %v", err)
	}
}
