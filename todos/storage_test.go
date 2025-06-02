package todos

import (
	"os"
	"testing"
)

func TestSaveAndLoad(t *testing.T) {
	tmp := t.TempDir()
	tmpFile := tmp + "/test-todos.json"

	origTodoFile := todoFile
	todoFile = func() string {
		return tmpFile
	}
	defer func() { todoFile = origTodoFile }()

	input := []Todo{
		{Task: "Learn Go", Completed: true},
		{Task: "Write tests", Completed: false},
	}

	if err := Save(input); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	output, err := Load()
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
	tmp := t.TempDir()
	tmpFile := tmp + "/does-not-exist.json"

	orig := todoFile
	todoFile = func() string { return tmpFile }
	defer func() { todoFile = orig }()

	todos, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(todos) != 0 {
		t.Fatalf("expected empty list, got %v", todos)
	}
}

func TestLoadWithInvalidJSON(t *testing.T) {
	tmp := t.TempDir()
	tmpFile := tmp + "/bad.json"

	if err := os.WriteFile(tmpFile, []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}

	orig := todoFile
	todoFile = func() string { return tmpFile }
	defer func() { todoFile = orig }()

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoadReadError(t *testing.T) {
	tmpFile := t.TempDir() + "/no-read-permission.json"

	// Datei anlegen und Rechte entziehen
	err := os.WriteFile(tmpFile, []byte(""), 0o000)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	defer func() { _ = os.Chmod(tmpFile, 0o644) }()

	orig := todoFile
	defer func() { todoFile = orig }()
	todoFile = func() string { return tmpFile }

	_, err = Load()
	if err == nil {
		t.Fatal("expected error from unreadable file, got nil")
	}
}
