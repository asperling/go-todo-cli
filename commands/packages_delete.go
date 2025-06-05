package commands

import (
	"github.com/urfave/cli/v2"

	"github.com/asperling/go-todo-cli/config"
	"github.com/asperling/go-todo-cli/todos"
)

func PackagesDeleteCommand(store *config.Store) *cli.Command {
	return &cli.Command{
		Name:        "delete",
		Aliases:     []string{"del", "rm"},
		Usage:       "Delete a package",
		Description: "Deletes the specified package file. The default package cannot be deleted.",
		ArgsUsage:   "[package name]",
		Action: func(c *cli.Context) error {
			return PackagesDeleteAction(c, store)
		},
	}
}

func PackagesDeleteAction(c *cli.Context, store *config.Store) error {
	if c.Args().Len() < 1 {
		return Exit("Usage: todo packages delete [package]")
	}

	name := c.Args().First()
	if !isValidPackageName(name) {
		return Exit("Invalid package name")
	}

	cfg, err := store.Load()
	if err != nil {
		return Exitf("Failed to load config: %v", err)
	}

	storage := todos.StorageFromConfig(&cfg)

	// Reset active package if deleted
	if name == cfg.ActivePackage {
		cfg.ActivePackage = config.DefaultPackage
		if errSave := store.Save(&cfg); errSave != nil {
			return Exitf("Failed to update active package: %v", errSave)
		}
	}

	if errDelete := storage.DeletePackage(name); errDelete != nil {
		return Exitf("Could not delete package: %v", errDelete)
	}

	SuccessPrintf("Deleted package: %s", name)
	return nil
}
