package commands

import (
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/asperling/go-todo-cli/config"
	"github.com/asperling/go-todo-cli/todos"
)

func AddAction(c *cli.Context, store *config.Store) error {
	if c.Args().Len() == 0 {
		return Exit("Usage: todo add [task]")
	}

	task := strings.Join(c.Args().Slice(), " ")

	cfg, err := store.Load()
	if err != nil {
		return Exitf("Failed to load config: %v", err)
	}

	storage := todos.StorageFromConfig(&cfg)
	list, err := storage.Load()
	if err != nil {
		return Exitf("Failed to load todos: %v", err)
	}

	if errAdd := todos.Add(&list, task); errAdd != nil {
		return Exitf("Failed to add todo: %v", errAdd)
	}

	if errSave := storage.Save(list); errSave != nil {
		return Exitf("Failed to save todos: %v", errSave)
	}

	SuccessPrintf("Added: %s", task)
	PrintList(list, cfg.ActivePackage)
	return nil
}

func AddCommand(store *config.Store) *cli.Command {
	return &cli.Command{
		Name:        "add",
		Aliases:     []string{"a"},
		Usage:       "Add a new task to your current package",
		Description: "Use this command to add a new task to your current package. You can specify the task description as an argument.",
		ArgsUsage:   "[task description]",
		Action: func(c *cli.Context) error {
			return AddAction(c, store)
		},
	}
}
