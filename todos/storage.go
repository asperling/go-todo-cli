package todos

import (
	"encoding/json"
	"os"
)

var todoFile = func() string {
	const fileName = ".aws-todos.json"
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fileName
	}
	return homeDir + "/" + fileName
}

func Load() ([]Todo, error) {
	if _, err := os.Stat(todoFile()); os.IsNotExist(err) {
		return []Todo{}, nil
	}

	data, err := os.ReadFile(todoFile())
	if err != nil {
		return nil, err
	}

	var todos []Todo
	err = json.Unmarshal(data, &todos)

	return todos, err
}

func Save(todos []Todo) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(todoFile(), data, 0o644)
}
