package main

import (
	"fmt"
	"log"
	"os"

	cmd "github.com/DocuBook/cli/cmd/docubook"
)

func main() {
	// Set up logging
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	// Check for version flag
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		fmt.Printf("DocuBook CLI\nVersion: %s\n", cmd.Version)
		return
	}

	// Execute the root command and handle any errors
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
