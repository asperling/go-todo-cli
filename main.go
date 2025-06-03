package main

import (
	"log"
	"os"

	"github.com/asperling/go-todo-cli/commands"
	"github.com/asperling/go-todo-cli/config"
	"github.com/urfave/cli/v2"
)

func main() {
	store := config.DefaultStore()

	app := &cli.App{
		Name:    "todo",
		Usage:   "Manage your todos from the command line",
		Version: "0.1.0",
		Commands: []*cli.Command{
			commands.InitCommand(&store),
			commands.ListCommand(&store),
			commands.AddCommand(&store),
			commands.DeleteCommand(&store),
			commands.DoneCommand(&store),
			commands.UndoneCommand(&store),
			commands.PackagesCommand(&store),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
