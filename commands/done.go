package commands

import (
	"strconv"

	"github.com/urfave/cli/v2"

	"github.com/asperling/go-todo-cli/config"
	"github.com/asperling/go-todo-cli/todos"
)

func DoneAction(c *cli.Context, store *config.Store) error {
	if c.Args().Len() < 1 {
		return Exit("Usage: todo done|undone [task number]")
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

	markDone := c.Command.Name == "done"

	if errDone := todos.Done(&list, index, markDone); errDone != nil {
		return Exitf("%v", errDone)
	}

	if errSave := storage.Save(list); errSave != nil {
		return Exitf("Failed to save todos: %v", errSave)
	}

	action := "completed"
	if !markDone {
		action = "reopened"
	}
	SuccessPrintf("Task %d %s", index, action)
	PrintList(list, cfg.ActivePackage)
	return nil
}

func DoneCommand(store *config.Store) *cli.Command {
	return &cli.Command{
		Name:        "done",
		Aliases:     []string{"d"},
		Usage:       "Mark a task as completed",
		Description: "Marks the task at the given position as completed.",
		ArgsUsage:   "[task number]",
		Action: func(c *cli.Context) error {
			return DoneAction(c, store)
		},
	}
}

func UndoneCommand(store *config.Store) *cli.Command {
	return &cli.Command{
		Name:        "undone",
		Aliases:     []string{"u"},
		Usage:       "Reopen a completed task",
		Description: "Marks the task at the given position as not completed.",
		ArgsUsage:   "[task number]",
		Action: func(c *cli.Context) error {
			return DoneAction(c, store)
		},
	}
}
