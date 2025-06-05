package commands

import (
	"strconv"

	"github.com/urfave/cli/v2"

	"github.com/asperling/go-todo-cli/config"
	"github.com/asperling/go-todo-cli/todos"
)

const (
	minArgs = 2 // Minimum number of arguments for move command
)

func MoveAction(c *cli.Context, store *config.Store) error {
	if c.Args().Len() < minArgs {
		return Exit("Usage: todo move [from] [to]")
	}

	from := c.Args().First()
	to := c.Args().Get(1)

	fromIndex, errFrom := strconv.Atoi(from)
	toIndex, errTo := strconv.Atoi(to)
	if errFrom != nil || errTo != nil {
		return Exit("Both 'from' and 'to' must be valid integers")
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

	if errAdd := todos.Move(&list, fromIndex, toIndex); errAdd != nil {
		return Exitf("Failed to move todo: %v", errAdd)
	}

	if errSave := storage.Save(list); errSave != nil {
		return Exitf("Failed to save todos: %v", errSave)
	}

	SuccessPrint("Successfully moved todo")
	PrintList(list, cfg.ActivePackage)
	return nil
}

func MoveCommand(store *config.Store) *cli.Command {
	return &cli.Command{
		Name:        "move",
		Aliases:     []string{"mv", "m"},
		Usage:       "Change the position of a task in your current package",
		Description: "Changes the position of a task within the current list. Useful for reordering priorities.",
		ArgsUsage:   "[from position] [to position]",
		Action: func(c *cli.Context) error {
			return MoveAction(c, store)
		},
	}
}
