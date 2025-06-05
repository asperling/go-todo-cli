package commands

import (
	"fmt"
	"strconv"

	"github.com/urfave/cli/v2"

	"github.com/asperling/go-todo-cli/config"
	"github.com/asperling/go-todo-cli/todos"
)

func DeleteAction(c *cli.Context, store *config.Store) error {
	if c.Args().Len() < 1 {
		return cli.Exit("❌ Usage: todo delete [index]", 1)
	}

	index, err := strconv.Atoi(c.Args().First())
	if err != nil {
		return cli.Exit("❌ Task number must be a valid integer", 1)
	}

	cfg, err := store.Load()
	if err != nil {
		return cli.Exit(fmt.Sprintf("❌ Failed to load config: %v", err), 1)
	}

	storage := todos.StorageFromConfig(&cfg)
	list, err := storage.Load()
	if err != nil {
		return cli.Exit(fmt.Sprintf("❌ Failed to load todos: %v", err), 1)
	}

	if errDelete := todos.Delete(&list, index); errDelete != nil {
		return cli.Exit(fmt.Sprintf("❌ %v", errDelete), 1)
	}

	if errSave := storage.Save(list); errSave != nil {
		return cli.Exit(fmt.Sprintf("❌ Failed to save todos: %v", errSave), 1)
	}

	fmt.Printf("🗑️  Deleted task at index %d\n", index)
	PrintList(list, cfg.ActivePackage)
	return nil
}

func DeleteCommand(store *config.Store) *cli.Command {
	return &cli.Command{
		Name:        "delete",
		Aliases:     []string{"del", "rm"},
		Usage:       "Delete a task by its index",
		Description: "Removes the task at the specified position in your current package.",
		ArgsUsage:   "[task number]",
		Action: func(c *cli.Context) error {
			return DeleteAction(c, store)
		},
	}
}
