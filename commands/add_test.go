package commands_test

import (
	"strings"
	"testing"

	"github.com/asperling/go-todo-cli/commands"
	"github.com/asperling/go-todo-cli/testutil"
	"github.com/asperling/go-todo-cli/todos"
	"github.com/urfave/cli/v2"
)

func TestAddAction_integration(t *testing.T) {
	cfg, store := testutil.ConfigSetup(t)

	app := testutil.App([]*cli.Command{
		commands.AddCommand(&store),
		commands.ListCommand(&store),
	})

	output := testutil.Capture(func() {
		err := app.Run([]string{"todo", "add", "Write unit tests"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(output, "✅ Added") {
		t.Errorf("missing success message: %q", output)
	}

	storage := todos.StorageFromConfig(cfg)
	list, err := storage.Load()
	if err != nil {
		t.Fatalf("could not reload todos: %v", err)
	}
	if len(list) != 1 || list[0].Task != "Write unit tests" {
		t.Errorf("todo not saved properly: %+v", list)
	}
}
