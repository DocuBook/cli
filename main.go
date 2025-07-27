package main

import (
	"os"

	cmd "github.com/DocuBook/cli/cmd/docubook"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		cmd.PrintVersion()
		return
	}

	cmd.Execute()
}
