package commands_test

import (
	"strings"
	"testing"

	"github.com/urfave/cli/v2"

	"github.com/asperling/go-todo-cli/commands"
	"github.com/asperling/go-todo-cli/testutil"
	"github.com/asperling/go-todo-cli/todos"
)

func TestDoneAction_integration(t *testing.T) {
	cfg, store := testutil.ConfigSetup(t)

	storage := todos.StorageFromConfig(cfg)
	list := []todos.Todo{
		{Task: "Write doc", Completed: false},
	}
	if err := storage.Save(list); err != nil {
		t.Fatalf("failed to save todos: %v", err)
	}

	app := testutil.App([]*cli.Command{
		commands.DoneCommand(&store),
		commands.UndoneCommand(&store),
	})

	// ✅ mark as done
	doneOutput := testutil.Capture(func() {
		if err := app.Run([]string{"todo", "done", "1"}); err != nil {
			t.Fatalf("done failed: %v", err)
		}
	})

	if !strings.Contains(doneOutput, "completed") {
		t.Errorf("expected 'completed' in output: %q", doneOutput)
	}

	todosAfterDone, err := storage.Load()
	if err != nil {
		t.Fatalf("failed to reload after done: %v", err)
	}
	if !todosAfterDone[0].Completed {
		t.Errorf("expected task to be marked as done, got: %+v", todosAfterDone[0])
	}

	// ❎ mark as undone
	undoneOutput := testutil.Capture(func() {
		if errRun := app.Run([]string{"todo", "undone", "1"}); errRun != nil {
			t.Fatalf("undone failed: %v", errRun)
		}
	})

	if !strings.Contains(undoneOutput, "reopened") {
		t.Errorf("expected 'reopened' in output: %q", undoneOutput)
	}

	todosAfterUndone, err := storage.Load()
	if err != nil {
		t.Fatalf("failed to reload after undone: %v", err)
	}
	if todosAfterUndone[0].Completed {
		t.Errorf("expected task to be reopened, got: %+v", todosAfterUndone[0])
	}
}
