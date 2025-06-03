package todos_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/asperling/go-todo-cli/config"
	"github.com/asperling/go-todo-cli/todos"
)

func newTestStorage(t *testing.T, packageName string) todos.Storage {
	t.Helper()
	return todos.Storage{
		Config: &config.Config{
			StoragePath:   t.TempDir(),
			ActivePackage: packageName,
		},
	}
}

func TestSaveAndLoad(t *testing.T) {
	storage := newTestStorage(t, "test")

	input := []todos.Todo{
		{Task: "Learn Go", Completed: true},
		{Task: "Write tests", Completed: false},
	}

	if err := storage.Save(input); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	output, err := storage.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(output) != len(input) {
		t.Fatalf("Loaded list has wrong length: got %d, want %d", len(output), len(input))
	}

	for i := range input {
		if input[i].Task != output[i].Task || input[i].Completed != output[i].Completed {
			t.Errorf("Mismatch at %d: got %+v, want %+v", i, output[i], input[i])
		}
	}
}

func TestLoadWhenFileDoesNotExist(t *testing.T) {
	storage := newTestStorage(t, "missing")

	todos, err := storage.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(todos) != 0 {
		t.Fatalf("expected empty list, got %v", todos)
	}
}

func TestLoadWithInvalidJSON(t *testing.T) {
	storage := newTestStorage(t, "bad")
	// Create invalid JSON manually
	path := filepath.Join(storage.Config.StoragePath, storage.Config.ActivePackage+".json")
	if err := os.WriteFile(path, []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := storage.Load()
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoadReadError(t *testing.T) {
	storage := newTestStorage(t, "restricted")

	path := filepath.Join(storage.Config.StoragePath, storage.Config.ActivePackage+".json")
	if err := os.WriteFile(path, []byte(""), 0o000); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	defer func() { _ = os.Chmod(path, 0o644) }()

	_, err := storage.Load()
	if err == nil {
		t.Fatal("expected error from unreadable file, got nil")
	}
}
