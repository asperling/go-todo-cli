//nolint:forbidigo // This is a CLI tool, so using fmt.Println is acceptable here
package todos

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func List(todos []Todo) {
	if len(todos) == 0 {
		fmt.Println("No todos found.")
		return
	}
	fmt.Println("Todo List:")
	for index, todo := range todos {
		status := "⬜"
		if todo.Completed {
			status = "✅"
		}
		fmt.Printf("%s [%d] %s\n", status, index+1, todo.Task)
	}
}

func Add(todos *[]Todo, task string) error {
	if task == "" {
		return errors.New("task cannot be empty")
	}
	todo := Todo{ID: uuid.New(), Task: task, Completed: false}
	*todos = append(*todos, todo)
	return nil
}

func Delete(todos *[]Todo, taskNumber int) error {
	if taskNumber < 1 || taskNumber > len(*todos) {
		return errors.New("invalid task number")
	}
	*todos = append((*todos)[:taskNumber-1], (*todos)[taskNumber:]...)
	return nil
}

func Move(todos *[]Todo, from, to int) error {
	if from < 1 || from > len(*todos) || to < 1 || to > len(*todos) || from == to {
		return errors.New("invalid positions for moving todo")
	}
	item := (*todos)[from-1]
	*todos = append((*todos)[:from-1], (*todos)[from:]...)
	if to > len(*todos) {
		*todos = append(*todos, item)
	} else {
		*todos = append((*todos)[:to-1], append([]Todo{item}, (*todos)[to-1:]...)...)
	}
	return nil
}

func Done(todos *[]Todo, taskNumber int, done bool) error {
	if taskNumber < 1 || taskNumber > len(*todos) {
		return errors.New("invalid task number")
	}
	(*todos)[taskNumber-1].Completed = done
	return nil
}
