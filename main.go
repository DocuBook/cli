package main

import (
	"log"
	"os"

	"github.com/DocuBook/cli/cmd"
)

func main() {
	// Set up logging
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	// Execute the root command and handle any errors
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
