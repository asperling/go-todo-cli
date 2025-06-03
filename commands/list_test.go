package commands_test

import (
	"strings"
	"testing"

	"github.com/asperling/go-todo-cli/commands"
	"github.com/asperling/go-todo-cli/testutil"
	"github.com/asperling/go-todo-cli/todos"
	"github.com/urfave/cli/v2"
)

func TestPrintList_empty(t *testing.T) {
	output := testutil.Capture(func() {
		commands.PrintList([]todos.Todo{}, "default")
	})

	if !strings.Contains(output, "No todos found") {
		t.Errorf("expected empty message, got: %q", output)
	}
}

func TestPrintList_withTodos(t *testing.T) {
	list := []todos.Todo{
		{Task: "Learn Go", Completed: false},
		{Task: "Write tests", Completed: true},
	}

	output := testutil.Capture(func() {
		commands.PrintList(list, "testpkg")
	})

	if !strings.Contains(output, "📝 Todo List (testpkg):") {
		t.Error("missing list header")
	}
	if !strings.Contains(output, "⬜ [1] Learn Go") {
		t.Error("missing or incorrect first task")
	}
	if !strings.Contains(output, "✅ [2] Write tests") {
		t.Error("missing or incorrect second task")
	}
}

func TestListAction_integration(t *testing.T) {
	cfg, store := testutil.ConfigSetup(t)

	storage := todos.StorageFromConfig(cfg)
	todosList := []todos.Todo{
		{Task: "A", Completed: false},
		{Task: "B", Completed: true},
	}
	if err := storage.Save(todosList); err != nil {
		t.Fatalf("failed to save todos: %v", err)
	}

	app := testutil.App([]*cli.Command{
		commands.ListCommand(&store),
	})

	output := testutil.Capture(func() {
		err := app.Run([]string{"todo", "list"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	if !strings.Contains(output, "Todo List") || !strings.Contains(output, "A") {
		t.Errorf("unexpected output: %q", output)
	}
}
