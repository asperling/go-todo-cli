package todos

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/asperling/go-todo-cli/config"
)

type Storage struct {
	Config *config.Config
}

var validPackageNameRegex = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

func IsValidPackageName(name string) bool {
	return validPackageNameRegex.MatchString(name)
}

func StorageFromConfig(cfg *config.Config) Storage {
	return Storage{Config: cfg}
}

func (s Storage) filePath() string {
	pkg := s.Config.ActivePackage
	if pkg == "" {
		pkg = config.DefaultPackage
	}
	return filepath.Join(s.Config.StoragePath, pkg+".json")
}

func (s Storage) Load() ([]Todo, error) {
	if _, err := os.Stat(s.filePath()); os.IsNotExist(err) {
		return []Todo{}, nil
	}
	data, err := os.ReadFile(s.filePath())
	if err != nil {
		return nil, err
	}
	var todos []Todo
	if unserializeErr := json.Unmarshal(data, &todos); unserializeErr != nil {
		return nil, unserializeErr
	}
	return todos, nil
}

func (s Storage) Save(todos []Todo) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath(), data, 0o600)
}

func (s Storage) ListPackages() ([]string, string, error) {
	entries, err := os.ReadDir(s.Config.StoragePath)
	if err != nil {
		return nil, "", err
	}

	var pkgs []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			name := strings.TrimSuffix(entry.Name(), ".json")
			pkgs = append(pkgs, name)
		}
	}

	active := s.Config.ActivePackage
	if active == "" {
		active = config.DefaultPackage
	}
	return pkgs, active, nil
}

func (s Storage) DeletePackage(name string) error {
	if name == config.DefaultPackage {
		return errors.New("default package cannot be deleted")
	}

	if !IsValidPackageName(name) {
		return errors.New("invalid package name")
	}

	path := filepath.Join(s.Config.StoragePath, name+".json")
	if err := os.Remove(path); err != nil {
		return err
	}
	return nil
}
