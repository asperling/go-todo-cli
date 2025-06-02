package todos_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/asperling/go-todo-cli/todos"
)

func TestList(t *testing.T) {
	todoList := []todos.Todo{
		{Task: "one", Completed: false},
		{Task: "two", Completed: true},
	}

	// redirect stdout to a buffer
	var buffer bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	//nolint:reassign // redirect os.Stdout to the pipe writer
	os.Stdout = w

	todos.List(todoList)

	// close the writer and restore stdout
	_ = w.Close()
	//nolint:reassign // restore os.Stdout to its original value
	os.Stdout = stdout
	_, _ = buffer.ReadFrom(r)
	output := buffer.String()

	want := "Todo List:\n⬜ [1] one\n✅ [2] two\n"
	if output != want {
		t.Errorf("unexpected output:\n--- got:\n%s\n--- want:\n%s", output, want)
	}
}

func TestListEmpty(t *testing.T) {
	// stdout abfangen
	var buf bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	//nolint:reassign // redirect os.Stdout to the pipe writer
	os.Stdout = w

	todos.List([]todos.Todo{})

	_ = w.Close()
	//nolint:reassign // restore os.Stdout to its original value
	os.Stdout = stdout
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	want := "No todos found.\n"
	if output != want {
		t.Errorf("unexpected output for empty list:\n--- got:\n%q\n--- want:\n%q", output, want)
	}
}

func TestAdd(t *testing.T) {
	todoList := []todos.Todo{}
	task := "Test task"

	if err := todos.Add(&todoList, task); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(todoList) != 1 {
		t.Fatalf("Expected 1 todo, got %d", len(todoList))
	}

	if todoList[0].Task != task {
		t.Fatalf("Expected task '%s', got '%s'", task, todoList[0].Task)
	}
}

func TestDelete(t *testing.T) {
	list := []todos.Todo{
		{Task: "eins"},
		{Task: "zwei"},
		{Task: "drei"},
	}

	err := todos.Delete(&list, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(list) != 2 || list[0].Task != "eins" || list[1].Task != "drei" {
		t.Errorf("unexpected list after delete: %+v", list)
	}
}

func TestDeleteInvalid(t *testing.T) {
	list := []todos.Todo{{Task: "eins"}}
	err := todos.Delete(&list, 5)

	if err == nil {
		t.Fatal("expected error for invalid delete position")
	}
}

func TestMove(t *testing.T) {
	list := []todos.Todo{
		{Task: "eins"},
		{Task: "zwei"},
		{Task: "drei"},
	}

	err := todos.Move(&list, 3, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if list[0].Task != "drei" || list[1].Task != "eins" || list[2].Task != "zwei" {
		t.Errorf("unexpected order after move: %+v", list)
	}
}

func TestMoveInvalid(t *testing.T) {
	list := []todos.Todo{{Task: "eins"}, {Task: "zwei"}}
	err := todos.Move(&list, 0, 5)

	if err == nil {
		t.Fatal("expected error for invalid move")
	}
}

func TestDone(t *testing.T) {
	list := []todos.Todo{
		{Task: "eins", Completed: false},
	}

	err := todos.Done(&list, 1, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !list[0].Completed {
		t.Errorf("expected task to be marked done")
	}
}

func TestDoneInvalid(t *testing.T) {
	list := []todos.Todo{{Task: "eins"}}
	err := todos.Done(&list, 5, true)

	if err == nil {
		t.Fatal("expected error for invalid task number")
	}
}
