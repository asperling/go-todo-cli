package testutil

import (
	"path/filepath"
	"testing"

	"github.com/asperling/go-todo-cli/config"
)

func ConfigSetup(t *testing.T) (*config.Config, config.Store) {
	t.Helper()

	tmp := t.TempDir()
	configPath := filepath.Join(tmp, "config.json")

	cfg := &config.Config{
		StoragePath:   tmp,
		ActivePackage: "default",
	}
	loader := config.Store{FilePath: configPath}

	if err := loader.Save(cfg); err != nil {
		t.Fatalf("failed to save test config: %v", err)
	}

	return cfg, loader
}
