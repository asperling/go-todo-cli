package commands

import (
	"fmt"

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

	// Reset active package if deleted
	if name == cfg.ActivePackage {
		cfg.ActivePackage = config.DefaultPackage
		if errSave := store.Save(&cfg); errSave != nil {
			return cli.Exit(fmt.Sprintf("❌ failed to update active package: %v", errSave), 1)
		}
	}

	if errDelete := storage.DeletePackage(name); errDelete != nil {
		return cli.Exit(fmt.Sprintf("❌ could not delete package: %v", errDelete), 1)
	}

	fmt.Printf("🗑️ Deleted package: %s\n", name)
	return nil
}
