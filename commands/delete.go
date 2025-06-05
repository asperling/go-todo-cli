package commands

import (
	"strconv"

	"github.com/urfave/cli/v2"

	"github.com/asperling/go-todo-cli/config"
	"github.com/asperling/go-todo-cli/todos"
)

func DeleteAction(c *cli.Context, store *config.Store) error {
	if c.Args().Len() < 1 {
		return Exit("Usage: todo delete [index]")
	}

	index, err := strconv.Atoi(c.Args().First())
	if err != nil {
		return Exit("Task number must be a valid integer")
	}

	cfg, err := store.Load()
	if err != nil {
		return Exitf("Failed to load config: %v", err)
	}

	storage := todos.StorageFromConfig(&cfg)
	list, err := storage.Load()
	if err != nil {
		return Exitf("Failed to load todos: %v", err)
	}

	if errDelete := todos.Delete(&list, index); errDelete != nil {
		return Exitf("%v", errDelete)
	}

	if errSave := storage.Save(list); errSave != nil {
		return Exitf("Failed to save todos: %v", errSave)
	}

	SuccessPrintf("Deleted task at index %d", index)
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
