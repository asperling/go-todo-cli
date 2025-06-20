package commands

import (
	"fmt"

	"github.com/asperling/go-todo-cli/config"
	"github.com/asperling/go-todo-cli/todos"
	"github.com/urfave/cli/v2"
)

func PrintList(list []todos.Todo, activePackage string) {
	if len(list) == 0 {
		fmt.Println("📭 No todos found.")
		return
	}

	fmt.Printf("📝 Todo List (%s):\n", activePackage)
	for i, todo := range list {
		status := "⬜"
		if todo.Completed {
			status = "✅"
		}
		fmt.Printf("%s [%d] %s\n", status, i+1, todo.Task)
	}
}

func ListAction(store *config.Store) error {
	cfg, err := store.Load()
	if err != nil {
		return Exitf("Failed to load config: %v", err)
	}

	storage := todos.StorageFromConfig(&cfg)
	list, err := storage.Load()
	if err != nil {
		return Exitf("Failed to load todos: %v", err)
	}

	PrintList(list, cfg.ActivePackage)
	return nil
}

func ListCommand(store *config.Store) *cli.Command {
	return &cli.Command{
		Name:        "list",
		Aliases:     []string{"ls", "l"},
		Usage:       "List all todos in the active package",
		Description: "Displays all tasks in the currently active package. Completed tasks are shown with a checkmark.",
		ArgsUsage:   "",
		Action: func(_ *cli.Context) error {
			return ListAction(store)
		},
	}
}
