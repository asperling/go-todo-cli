package commands

import (
	"github.com/urfave/cli/v2"

	"github.com/asperling/go-todo-cli/config"
)

func PackagesUseCommand(store *config.Store) *cli.Command {
	return &cli.Command{
		Name:        "use",
		Aliases:     []string{"u"},
		Usage:       "Switch to a different package",
		Description: "Sets the active package. If the package does not exist, it will be created automatically when adding a new task.",
		ArgsUsage:   "[package name]",
		Action: func(c *cli.Context) error {
			return PackagesUseAction(c, store)
		},
	}
}

func PackagesUseAction(c *cli.Context, store *config.Store) error {
	if c.Args().Len() < 1 {
		return Exit("Usage: todo packages use [package]")
	}

	name := c.Args().First()
	if !isValidPackageName(name) {
		return Exit("Package name may only contain letters and numbers")
	}

	cfg, err := store.Load()
	if err != nil {
		return Exitf("Failed to load config: %v", err)
	}

	cfg.ActivePackage = name
	if errSave := store.Save(&cfg); errSave != nil {
		return Exitf("Failed to update config: %v", errSave)
	}

	SuccessPrintf("Switched to package: %s", name)
	return nil
}
