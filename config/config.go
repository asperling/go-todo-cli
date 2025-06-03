package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func Path() string {
	paths := []string{".aws-todo", "config.json"}
	home, err := os.UserHomeDir()
	if err != nil {
		// fallback: relative path
		return filepath.Join(".", filepath.Join(paths...))
	}
	return filepath.Join(append([]string{home}, paths...)...)
}

func Load() (Config, error) {
	path := Path()
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var config Config
	if errUnmarshal := json.Unmarshal(data, &config); errUnmarshal != nil {
		return Config{}, errUnmarshal
	}
	return config, nil
}

func Save(config Config) error {
	path := Path()
	dir := filepath.Dir(path)
	if errMkdir := os.MkdirAll(dir, 0o700); errMkdir != nil {
		return errMkdir
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}
