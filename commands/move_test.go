package commands_test

import (
	"strings"
	"testing"

	"github.com/urfave/cli/v2"

	"github.com/asperling/go-todo-cli/commands"
	"github.com/asperling/go-todo-cli/testutil"
	"github.com/asperling/go-todo-cli/todos"
)

func TestMoveAction_integration(t *testing.T) {
	cfg, store := testutil.ConfigSetup(t)

	// Seed todos: A, B, C
	initial := []todos.Todo{
		{Task: "A"},
		{Task: "B"},
		{Task: "C"},
	}
	storage := todos.StorageFromConfig(cfg)
	if err := storage.Save(initial); err != nil {
		t.Fatalf("failed to seed todos: %v", err)
	}

	app := testutil.App([]*cli.Command{
		commands.MoveCommand(&store),
	})

	output := testutil.Capture(func() {
		err := app.Run([]string{"todo", "move", "3", "1"}) // move C to top
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(output, "Successfully moved") {
		t.Errorf("expected success message, got: %q", output)
	}
	if !strings.Contains(output, "[1] C") || !strings.Contains(output, "[2] A") {
		t.Errorf("unexpected order after move: %q", output)
	}

	// Re-load and assert persisted order
	after, err := storage.Load()
	if err != nil {
		t.Fatalf("failed to reload todos: %v", err)
	}
	if len(after) != 3 || after[0].Task != "C" || after[1].Task != "A" || after[2].Task != "B" {
		t.Errorf("unexpected order in storage: %+v", after)
	}
}
