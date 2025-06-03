package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/asperling/go-todo-cli/config"
	"github.com/asperling/go-todo-cli/todos"
)

func UsePackage(args []string) error {
	if len(args) < UseArgs {
		return errors.New("usage: todo use [package]")
	}
	name := strings.TrimSpace(args[2])
	if !todos.IsValidPackageName(name) {
		return fmt.Errorf("invalid package name: %q (only letters, digits, '-', and '_' allowed)", name)
	}

	cfg, errLoad := config.Load()
	if errLoad != nil {
		return fmt.Errorf("failed to load config: %w", errLoad)
	}
	cfg.ActivePackage = name
	if errSave := cfg.Save(); errSave != nil {
		return fmt.Errorf("failed to save config: %w", errSave)
	}

	fmt.Printf("✅ Switched to package: %s\n", name)
	return nil
}
