package commands

import (
	"errors"
	"fmt"

	"github.com/asperling/go-todo-cli/config"
	"github.com/asperling/go-todo-cli/todos"
)

func Run(args []string) error {
	if len(args) < MinArgs {
		return errors.New("usage: todo [add|list|done|undone|move|delete] [...]")
	}

	command := args[1]

	// Handle init command before the configuration is loaded as it might only be available after initialization.
	if command == InitCommand {
		return Init()
	}

	if command == "use" {
		if errUsePackage := UsePackage(args); errUsePackage != nil {
			return errUsePackage
		}
		command = ListCommand // After using a package, we default to listing todos.
	}

	config, configError := config.Load()
	if configError != nil {
		return fmt.Errorf("error loading storage configuration: %w, did you run `todo init`?", configError)
	}

	// early routing for non-todo commands
	switch command {
	case ListPackagesCommand:
		storage := todos.StorageFromConfig(&config)
		return handleListPackages(storage)

	case DeletePackageCommand:
		storage := todos.StorageFromConfig(&config)
		return handleDeletePackage(args, &storage)
	}

	storage := todos.StorageFromConfig(&config)
	todoList, err := storage.Load()
	if err != nil {
		return fmt.Errorf("error loading todos: %w", err)
	}

	switch command {
	case "add":
		err = handleAdd(args, &todoList)
	case ListCommand:
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

	if command != ListCommand && command != InitCommand {
		if saveErr := storage.Save(todoList); saveErr != nil {
			return saveErr
		}
		fmt.Println("Todos updated successfully.")
		todos.List(todoList)
	}

	return nil
}

func handleAdd(args []string, todosRef *[]todos.Todo) error {
	if _, err := ValidateArgs(args, AddArgs, []int{}, "Usage: todo add '[task]'"); err != nil {
		return err
	}
	return todos.Add(todosRef, args[2])
}

func handleList(todosRef *[]todos.Todo) {
	todos.List(*todosRef)
}

func handleMove(args []string, todosRef *[]todos.Todo) error {
	ints, err := ValidateArgs(args, MoveArgs, []int{ArgIndexFrom, ArgIndexTo}, "Usage: todo move [from] [to]")
	if err != nil {
		return err
	}
	return todos.Move(todosRef, ints[0], ints[1])
}

func handleDone(args []string, todosRef *[]todos.Todo, markDone bool) error {
	ints, err := ValidateArgs(args, DoneArgs, []int{2}, "Usage: todo done|undone [task number]")
	if err != nil {
		return err
	}
	return todos.Done(todosRef, ints[0], markDone)
}

func handleDelete(args []string, todosRef *[]todos.Todo) error {
	ints, err := ValidateArgs(args, DeleteArgs, []int{2}, "Usage: todo delete [task number]")
	if err != nil {
		return err
	}
	return todos.Delete(todosRef, ints[0])
}

func handleListPackages(storage todos.Storage) error {
	pkgs, active, err := storage.ListPackages()
	if err != nil {
		return fmt.Errorf("failed to list packages: %w", err)
	}

	fmt.Println("Available packages:")
	for _, name := range pkgs {
		mark := " "
		if name == active {
			mark = "*"
		}
		fmt.Printf("  %s %s\n", mark, name)
	}
	return nil
}

func handleDeletePackage(args []string, storage *todos.Storage) error {
	if _, err := ValidateArgs(args, DeletePackageArgs, []int{}, "Usage: todo delete-package [name]"); err != nil {
		return err
	}

	name := args[2]
	if err := storage.DeletePackage(name); err != nil {
		return fmt.Errorf("failed to delete package: %w", err)
	}

	fmt.Printf("✅ Deleted package: %s\n", name)
	return nil
}
