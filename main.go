package main

import (
	"fmt"
	"os"

	"github.com/asperling/go-todo-cli/commands"
)

func main() {
	if err := commands.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
