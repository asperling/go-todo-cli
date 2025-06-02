package todos_test

import (
	"os"
	"testing"

	"github.com/asperling/go-todo-cli/todos"
)

func TestSaveAndLoad(t *testing.T) {
	storage := todos.Storage{FilePath: t.TempDir() + "test.json"}

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
	storage := todos.Storage{FilePath: t.TempDir() + "/does-not-exist.json"}

	todos, err := storage.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(todos) != 0 {
		t.Fatalf("expected empty list, got %v", todos)
	}
}

func TestLoadWithInvalidJSON(t *testing.T) {
	storage := todos.Storage{FilePath: t.TempDir() + "/bad.json"}
	if err := os.WriteFile(storage.FilePath, []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := storage.Load()
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoadReadError(t *testing.T) {
	storage := todos.Storage{FilePath: t.TempDir() + "/no-read-permission.json"}

	// Datei anlegen und Rechte entziehen
	err := os.WriteFile(storage.FilePath, []byte(""), 0o000)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	defer func() { _ = os.Chmod(storage.FilePath, 0o644) }()

	_, err = storage.Load()
	if err == nil {
		t.Fatal("expected error from unreadable file, got nil")
	}
}
