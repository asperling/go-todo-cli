package todos

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Storage struct {
	FilePath string
}

func DefaultStorage() Storage {
	home, err := os.UserHomeDir()
	const filename = ".aws-todos.json"
	if err != nil {
		return Storage{FilePath: filename}
	}
	return Storage{FilePath: filepath.Join(home, filename)}
}

func (s Storage) Load() ([]Todo, error) {
	if _, err := os.Stat(s.FilePath); os.IsNotExist(err) {
		return []Todo{}, nil
	}
	data, err := os.ReadFile(s.FilePath)
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
	return os.WriteFile(s.FilePath, data, 0o600)
}
