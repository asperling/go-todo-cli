package commands_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/urfave/cli/v2"

	"github.com/asperling/go-todo-cli/commands"
	"github.com/asperling/go-todo-cli/testutil"
	"github.com/asperling/go-todo-cli/todos"
)

func TestPackagesList_and_Use_and_Delete(t *testing.T) {
	cfg, store := testutil.ConfigSetup(t)
	storage := todos.StorageFromConfig(cfg)

	// seed: mehrere packages
	for _, name := range []string{"default", "demo", "work"} {
		dummy := []todos.Todo{{Task: "test"}}
		path := filepath.Join(cfg.StoragePath, name+".json")
		if err := os.WriteFile(path, []byte(`[{"task":"test"}]`), 0o600); err != nil {
			t.Fatalf("failed to seed package %s: %v", name, err)
		}
		if err := storage.Save(dummy); err != nil {
			t.Fatalf("failed to save package %s: %v", name, err)
		}
	}

	app := testutil.App([]*cli.Command{
		commands.PackagesCommand(&store),
	})

	output := testutil.Capture(func() {
		err := app.Run([]string{"todo", "packages"})
		if err != nil {
			t.Fatalf("packages list failed: %v", err)
		}
	})

	if !strings.Contains(output, "default") || !strings.Contains(output, "demo") {
		t.Errorf("expected package names in output, got: %q", output)
	}

	output = testutil.Capture(func() {
		err := app.Run([]string{"todo", "packages", "use", "demo"})
		if err != nil {
			t.Fatalf("packages use failed: %v", err)
		}
	})

	if !strings.Contains(output, "Switched to package: demo") {
		t.Errorf("unexpected use output: %q", output)
	}

	newCfg, err := store.Load()
	if err != nil {
		t.Fatalf("failed to reload config: %v", err)
	}
	if newCfg.ActivePackage != "demo" {
		t.Errorf("expected active package to be 'demo', got: %s", newCfg.ActivePackage)
	}

	_ = testutil.Capture(func() {
		errRun := app.Run([]string{"todo", "packages", "delete", "demo"})
		if errRun != nil {
			t.Fatalf("packages delete failed: %v", errRun)
		}
	})

	if _, errStat := os.Stat(filepath.Join(cfg.StoragePath, "demo.json")); errStat == nil {
		t.Error("expected demo.json to be deleted")
	}
}
