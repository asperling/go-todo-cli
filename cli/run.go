package cli

import (
	"errors"
	"fmt"

	"github.com/asperling/go-todo-cli/todos"
)

const (
	minArgs        = 2
	addArgs        = 3
	moveArgs       = 4
	doneArgs       = 3
	deleteArgs     = 3
	argIndexFrom   = 2
	argIndexTo     = 3
	argIndexNumber = 2
	initCommand    = "init"
)

func Run(args []string) error {
	if len(args) < minArgs {
		return errors.New("usage: todo [add|list|done|undone|move|delete] [...]")
	}

	command := args[1]

	// Handle init command before the configuration is loaded as it might only be available after initialization.
	if command == initCommand {
		return Init()
	}

	storage, configError := todos.StorageFromConfig()
	if configError != nil {
		return fmt.Errorf("error loading storage configuration: %w, did you run `todo init`?", configError)
	}
	todoList, err := storage.Load()
	if err != nil {
		return fmt.Errorf("error loading todos: %w", err)
	}

	switch command {
	case "add":
		err = handleAdd(args, &todoList)
	case "list":
		handleList(&todoList)
	case "move":
		err = handleMove(args, &todoList)
	case "done", "undone":
		err = handleDone(args, &todoList, command == "done")
	case "delete":
		err = handleDelete(args, &todoList)
	default:
		err = fmt.Errorf("unknown command: %s", command)
	}

	if err != nil {
		return err
	}

	if command != "list" && command != initCommand {
		if saveErr := storage.Save(todoList); saveErr != nil {
			return saveErr
		}
		fmt.Println("Todos updated successfully.")
		todos.List(todoList)
	}

	return nil
}

func handleAdd(args []string, todosRef *[]todos.Todo) error {
	if _, err := ValidateArgs(args, addArgs, []int{}, "Usage: todo add '[task]'"); err != nil {
		return err
	}
	return todos.Add(todosRef, args[2])
}

func handleList(todosRef *[]todos.Todo) {
	todos.List(*todosRef)
}

func handleMove(args []string, todosRef *[]todos.Todo) error {
	ints, err := ValidateArgs(args, moveArgs, []int{argIndexFrom, argIndexTo}, "Usage: todo move [from] [to]")
	if err != nil {
		return err
	}
	return todos.Move(todosRef, ints[0], ints[1])
}

func handleDone(args []string, todosRef *[]todos.Todo, markDone bool) error {
	ints, err := ValidateArgs(args, doneArgs, []int{2}, "Usage: todo done|undone [task number]")
	if err != nil {
		return err
	}
	return todos.Done(todosRef, ints[0], markDone)
}

func handleDelete(args []string, todosRef *[]todos.Todo) error {
	ints, err := ValidateArgs(args, deleteArgs, []int{2}, "Usage: todo delete [task number]")
	if err != nil {
		return err
	}
	return todos.Delete(todosRef, ints[0])
}
