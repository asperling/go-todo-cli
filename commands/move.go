package commands

import (
	"fmt"
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
		return cli.Exit("❌ Usage: todo move [from] [to]", 1)
	}

	from := c.Args().First()
	to := c.Args().Get(1)

	fromIndex, errFrom := strconv.Atoi(from)
	toIndex, errTo := strconv.Atoi(to)
	if errFrom != nil || errTo != nil {
		return cli.Exit("❌ Both 'from' and 'to' must be valid integers", 1)
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

	if errAdd := todos.Move(&list, fromIndex, toIndex); errAdd != nil {
		return cli.Exit(fmt.Sprintf("❌ Failed to move todo: %v", errAdd), 1)
	}

	if errSave := storage.Save(list); errSave != nil {
		return cli.Exit(fmt.Sprintf("❌ Failed to save todos: %v", errSave), 1)
	}

	fmt.Println("✔ Successfully moved todo")
	PrintList(list, cfg.ActivePackage)
	return nil
}

func MoveCommand(store *config.Store) *cli.Command {
	return &cli.Command{
		Name:    "move",
		Aliases: []string{"mv", "m"},
		Usage:   "Change the position of a task in your current package",
		Action: func(c *cli.Context) error {
			return MoveAction(c, store)
		},
	}
}
