package commands_test

import (
	"strings"
	"testing"

	"github.com/urfave/cli/v2"

	"github.com/asperling/go-todo-cli/commands"
	"github.com/asperling/go-todo-cli/testutil"
	"github.com/asperling/go-todo-cli/todos"
)

func TestDeleteAction_integration(t *testing.T) {
	cfg, store := testutil.ConfigSetup(t)

	storage := todos.StorageFromConfig(cfg)
	initial := []todos.Todo{
		{Task: "Task A", Completed: false},
		{Task: "Task B", Completed: false},
	}
	if err := storage.Save(initial); err != nil {
		t.Fatalf("failed to save initial todos: %v", err)
	}

	app := testutil.App([]*cli.Command{
		commands.DeleteCommand(&store),
	})

	output := testutil.Capture(func() {
		err := app.Run([]string{"todo", "delete", "1"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(output, "Deleted task") {
		t.Errorf("missing success message, got: %q", output)
	}

	// Only "Task B" should remain after deletion
	remaining, err := storage.Load()
	if err != nil {
		t.Fatalf("failed to load remaining todos: %v", err)
	}

	if len(remaining) != 1 || remaining[0].Task != "Task B" {
		t.Errorf("unexpected todo list after delete: %+v", remaining)
	}
}
