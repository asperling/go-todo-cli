//nolint:reassign // os.Stdout and os.Stderr are reassigned to capture output
package commands_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/urfave/cli/v2"

	"github.com/asperling/go-todo-cli/commands"
	"github.com/asperling/go-todo-cli/config"
	"github.com/asperling/go-todo-cli/testutil"
)

func TestInitAction_interactive(t *testing.T) {
	tmpDir := t.TempDir()
	targetPath := filepath.Join(tmpDir, "my-todos")
	storePath := filepath.Join(tmpDir, "config.json")

	// simulate user input (targetPath + newline)
	input := targetPath + "\n"
	r, w, _ := os.Pipe()
	_, _ = w.WriteString(input)
	_ = w.Close()
	origStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = origStdin }()

	store := &config.Store{FilePath: storePath}
	app := testutil.App([]*cli.Command{
		commands.InitCommand(store),
	})

	output := testutil.Capture(func() {
		if err := app.Run([]string{"todo", "init"}); err != nil {
			t.Fatalf("init failed: %v", err)
		}
	})

	if !strings.Contains(output, "✅ Configuration saved") {
		t.Errorf("expected success message, got: %q", output)
	}

	// Check that the directory was created
	if _, err := os.Stat(targetPath); err != nil {
		t.Errorf("expected directory to be created, got error: %v", err)
	}

	// Check config values
	cfg, err := store.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	if cfg.StoragePath != targetPath {
		t.Errorf("unexpected storage path: %s", cfg.StoragePath)
	}
	if cfg.ActivePackage != config.DefaultPackage {
		t.Errorf("unexpected default package: %s", cfg.ActivePackage)
	}
}

func TestInitAction_empty_input_uses_default(t *testing.T) {
	tmp := t.TempDir()
	storePath := filepath.Join(tmp, "config.json")

	// Simuliere [Enter] → leere Eingabe
	r, w, _ := os.Pipe()
	_, _ = w.WriteString("\n")
	_ = w.Close()
	origStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = origStdin }()

	store := &config.Store{FilePath: storePath}
	app := testutil.App([]*cli.Command{
		commands.InitCommand(store),
	})

	_ = testutil.Capture(func() {
		if err := app.Run([]string{"todo", "init"}); err != nil {
			t.Fatalf("init failed: %v", err)
		}
	})

	cfg, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(cfg.StoragePath, config.FolderName) {
		t.Errorf("expected default path to end in %q, got: %s", config.FolderName, cfg.StoragePath)
	}
}

// That was fiddly, as en error causes cli to exit immediately,
// so we need to capture the error handler.
func TestInitAction_path_is_file(t *testing.T) {
	tmp := t.TempDir()

	// Eine Datei anstelle eines Verzeichnisses erzeugen
	file := filepath.Join(tmp, "somefile")
	if err := os.WriteFile(file, []byte("dummy"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Simulate user input: file + \n
	r, w, _ := os.Pipe()
	_, _ = w.WriteString(file + "\n")
	_ = w.Close()
	origStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = origStdin }()

	store := &config.Store{FilePath: filepath.Join(tmp, "cfg.json")}
	app := testutil.App([]*cli.Command{
		commands.InitCommand(store),
	})

	var exitErr error
	app.ExitErrHandler = func(_ *cli.Context, err error) {
		exitErr = err
	}

	_ = app.Run([]string{"todo", "init"})

	if exitErr == nil {
		t.Fatal("expected cli.Exit error, got nil")
	}

	// This should contain the error message about not being a directory
	msg := exitErr.Error()

	if !strings.Contains(msg, "not a directory") {
		t.Errorf("expected error to mention 'not a directory', got: %s", msg)
	}
}

func TestInitAction_unwritable_directory(t *testing.T) {
	parent := t.TempDir()
	noWrite := filepath.Join(parent, "nowrite")
	if err := os.Mkdir(noWrite, 0o400); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Chmod(noWrite, 0o700)
	}()

	target := filepath.Join(noWrite, "todo")

	r, w, _ := os.Pipe()
	_, _ = w.WriteString(target + "\n")
	_ = w.Close()
	origStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = origStdin }()

	store := &config.Store{FilePath: filepath.Join(parent, "cfg.json")}
	app := testutil.App([]*cli.Command{
		commands.InitCommand(store),
	})

	var exitErr error
	app.ExitErrHandler = func(_ *cli.Context, err error) {
		exitErr = err
	}

	_ = app.Run([]string{"todo", "init"})

	if exitErr == nil {
		t.Fatal("expected cli.Exit error, got nil")
	}

	msg := exitErr.Error()
	if !strings.Contains(msg, "could not create") && !strings.Contains(msg, "failed to access path") {
		t.Errorf("expected creation error, got: %s", msg)
	}
}
