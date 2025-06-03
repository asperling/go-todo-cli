package commands

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/asperling/go-todo-cli/config"
	"github.com/asperling/go-todo-cli/todos"
)

func AddAction(c *cli.Context, store *config.Store) error {
	if c.Args().Len() == 0 {
		return cli.Exit("❌ Usage: todo add [task]", 1)
	}

	task := strings.Join(c.Args().Slice(), " ")

	cfg, err := store.Load()
	if err != nil {
		return cli.Exit(fmt.Sprintf("❌ Failed to load config: %v", err), 1)
	}

	storage := todos.StorageFromConfig(&cfg)
	list, err := storage.Load()
	if err != nil {
		return cli.Exit(fmt.Sprintf("❌ Failed to load todos: %v", err), 1)
	}

	if errAdd := todos.Add(&list, task); errAdd != nil {
		return cli.Exit(fmt.Sprintf("❌ Failed to add todo: %v", errAdd), 1)
	}

	if errSave := storage.Save(list); errSave != nil {
		return cli.Exit(fmt.Sprintf("❌ Failed to save todos: %v", errSave), 1)
	}

	fmt.Printf("✅ Added: %s\n", task)
	return nil
}

func AddCommand(store *config.Store) *cli.Command {
	return &cli.Command{
		Name:  "add",
		Usage: "Add a new task to your current package",
		Action: func(c *cli.Context) error {
			return AddAction(c, store)
		},
	}
}
