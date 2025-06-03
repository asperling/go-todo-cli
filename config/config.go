package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func Path() string {
	paths := []string{FolderName, FileName}
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

func (cfg *Config) Save() error {
	dir := filepath.Dir(Path())
	if errMkdir := os.MkdirAll(dir, 0o700); errMkdir != nil {
		return errMkdir
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(Path(), data, 0o600)
}
