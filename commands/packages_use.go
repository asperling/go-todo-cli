package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/asperling/go-todo-cli/config"
)

func PackagesUseCommand(store *config.Store) *cli.Command {
	return &cli.Command{
		Name:      "use",
		Aliases:   []string{"u"},
		Usage:     "Switch to a different package",
		ArgsUsage: "[package name]",
		Action: func(c *cli.Context) error {
			return PackagesUseAction(c, store)
		},
	}
}

func PackagesUseAction(c *cli.Context, store *config.Store) error {
	if c.Args().Len() < 1 {
		return cli.Exit("❌ Usage: todo packages use [package]", 1)
	}

	name := c.Args().First()
	if !isValidPackageName(name) {
		return cli.Exit("❌ Package name may only contain letters and numbers", 1)
	}

	cfg, err := store.Load()
	if err != nil {
		return cli.Exit(fmt.Sprintf("❌ Failed to load config: %v", err), 1)
	}

	cfg.ActivePackage = name
	if errSave := store.Save(&cfg); errSave != nil {
		return cli.Exit(fmt.Sprintf("❌ Failed to update config: %v", errSave), 1)
	}

	fmt.Printf("📦 Switched to package: %s\n", name)
	return nil
}
