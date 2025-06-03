package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Store struct {
	FilePath string
}

func DefaultStore() Store {
	home, err := os.UserHomeDir()
	if err != nil {
		// fallback: relative path
		return Store{FilePath: filepath.Join(".", FolderName, FileName)}
	}
	path := filepath.Join(home, FolderName, FileName)
	return Store{FilePath: path}
}

func (s Store) Load() (Config, error) {
	data, err := os.ReadFile(s.FilePath)
	if err != nil {
		return Config{}, err
	}
	var config Config
	if errUnmarshal := json.Unmarshal(data, &config); errUnmarshal != nil {
		return Config{}, errUnmarshal
	}
	return config, nil
}

func (s Store) Save(cfg *Config) error {
	dir := filepath.Dir(s.FilePath)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.FilePath, data, 0o600)
}
