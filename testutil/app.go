package testutil

import (
	cli "github.com/urfave/cli/v2"
)

func App(commands []*cli.Command) *cli.App {
	return &cli.App{
		Name:     "todo-test",
		Commands: commands,
	}
}
