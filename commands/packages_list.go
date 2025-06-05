package commands

import (
	"fmt"

	"github.com/asperling/go-todo-cli/config"
	"github.com/asperling/go-todo-cli/todos"
	"github.com/urfave/cli/v2"
)

func PackagesListCommand(store *config.Store) *cli.Command {
	return &cli.Command{
		Name:        "list",
		Aliases:     []string{"ls", "l"},
		Usage:       "List available packages",
		Description: "Shows all available packages. The active one is marked.",
		ArgsUsage:   "",
		Action:      func(_ *cli.Context) error { return PackagesListAction(store) },
	}
}

func PackagesListAction(store *config.Store) error {
	cfg, err := store.Load()
	if err != nil {
		return cli.Exit(fmt.Sprintf("❌ Failed to load config: %v", err), 1)
	}

	storage := todos.StorageFromConfig(&cfg)
	pkgs, active, err := storage.ListPackages()
	if err != nil {
		return cli.Exit(fmt.Sprintf("❌ Failed to list packages: %v", err), 1)
	}

	fmt.Println("📦 Available packages:")
	for _, p := range pkgs {
		if p == active {
			fmt.Printf("• %s (active)\n", p)
		} else {
			fmt.Printf("• %s\n", p)
		}
	}
	return nil
}
