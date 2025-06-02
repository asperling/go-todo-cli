package cli

import (
	"errors"
	"fmt"

	"github.com/asperling/go-todo-cli/todos"
)

func Run(args []string) error {
	if len(args) < 2 {
		return errors.New("Usage: todo [add|list|done|undone|move|delete] [...]")
	}

	command := args[1]
	todoList, err := todos.Load()
	if err != nil {
		return fmt.Errorf("error loading todos: %w", err)
	}

	switch command {
	case "add":
		err = handleAdd(args, &todoList)
	case "list":
		err = handleList(&todoList)
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

	if command != "list" {
		if err := todos.Save(todoList); err != nil {
			return err
		}
		//nolint:forbidigo // Print a success message only for commands that modify the todo list
		fmt.Println("Todos updated successfully.")
		todos.List(todoList)
	}

	return nil
}

func handleAdd(args []string, todosRef *[]todos.Todo) error {
	if _, err := ValidateArgs(args, 3, []int{}, "Usage: todo add '[task]'"); err != nil {
		return err
	}
	return todos.Add(todosRef, args[2])
}

func handleList(todosRef *[]todos.Todo) error {
	todos.List(*todosRef)
	return nil
}

func handleMove(args []string, todosRef *[]todos.Todo) error {
	ints, err := ValidateArgs(args, 4, []int{2, 3}, "Usage: todo move [from] [to]")
	if err != nil {
		return err
	}
	return todos.Move(todosRef, ints[0], ints[1])
}

func handleDone(args []string, todosRef *[]todos.Todo, markDone bool) error {
	ints, err := ValidateArgs(args, 3, []int{2}, "Usage: todo done|undone [task number]")
	if err != nil {
		return err
	}
	return todos.Done(todosRef, ints[0], markDone)
}

func handleDelete(args []string, todosRef *[]todos.Todo) error {
	ints, err := ValidateArgs(args, 3, []int{2}, "Usage: todo delete [task number]")
	if err != nil {
		return err
	}
	return todos.Delete(todosRef, ints[0])
}
