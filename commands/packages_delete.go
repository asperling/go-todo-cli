package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/asperling/go-todo-cli/config"
	"github.com/asperling/go-todo-cli/todos"
)

func PackagesDeleteCommand(store *config.Store) *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Usage:     "Delete a package",
		ArgsUsage: "[package name]",
		Action: func(c *cli.Context) error {
			return PackagesDeleteAction(c, store)
		},
	}
}

func PackagesDeleteAction(c *cli.Context, store *config.Store) error {
	if c.Args().Len() < 1 {
		return cli.Exit("❌ Usage: todo packages delete [package]", 1)
	}

	name := c.Args().First()
	if !isValidPackageName(name) {
		return cli.Exit("❌ Invalid package name", 1)
	}

	cfg, err := store.Load()
	if err != nil {
		return cli.Exit(fmt.Sprintf("❌ Failed to load config: %v", err), 1)
	}

	storage := todos.StorageFromConfig(&cfg)

	if errDelete := storage.DeletePackage(name); errDelete != nil {
		return cli.Exit(fmt.Sprintf("❌ %v", errDelete), 1)
	}

	return nil
}
