package main

import (
	"fmt"
	"os"

	"github.com/asperling/go-todo-cli/cli"
	"github.com/asperling/go-todo-cli/todos"
)

func main() {
	todoList, err := todos.Load()
	if err != nil {
		fmt.Println("Error loading todos:", err)
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		if _, errArgs := cli.ValidateArgs(3, []int{}, "Usage: todo add '[task]'"); errArgs != nil {
			fmt.Println(errArgs)
			return
		}

		if err := todos.Add(&todoList, os.Args[2]); err != nil {
			fmt.Println("Error adding todo:", err)
			return
		}

	case "list":
		todos.List(todoList)
	case "move":
		intArgs, errArgs := cli.ValidateArgs(4, []int{2, 3}, "Usage: todo move [from] [to]")
		if errArgs != nil {
			fmt.Println(errArgs)
			return
		}
		if err := todos.Move(&todoList, intArgs[0], intArgs[1]); err != nil {
			fmt.Println("Error moving todo:", err)
			return
		}

	case "done", "undone":
		intArgs, errArgs := cli.ValidateArgs(3, []int{2}, "Usage: todo done [task number] or todo undone [task number]")
		if errArgs != nil {
			fmt.Println(errArgs)
			return
		}

		if err := todos.Done(&todoList, intArgs[0], command == "done"); err != nil {
			fmt.Println("Error marking todo as done/undone:", err)
			return
		}

	case "delete":
		intArgs, errArgs := cli.ValidateArgs(3, []int{2}, "Usage: todo delete [task number]")
		if errArgs != nil {
			fmt.Println(errArgs)
			return
		}
		if err := todos.Delete(&todoList, intArgs[0]); err != nil {
			fmt.Println("Error deleting todo:", err)
			return
		}

	default:
		fmt.Println("Unknown command:", command)
		return
	}

	if command != "list" {
		if err := todos.Save(todoList); err != nil {
			fmt.Println("Error saving todos:", err)
			return
		}
		fmt.Println("Todos updated successfully.")
		todos.List(todoList)
	}
}
