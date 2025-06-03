package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/asperling/go-todo-cli/config"
)

func Init() error {
	reader := bufio.NewReader(os.Stdin)

	home, _ := os.UserHomeDir()
	defaultPath := filepath.Join(home, config.FolderName)

	fmt.Printf("Where would you like to store your todos? [%s]: ", defaultPath)
	input, _ := reader.ReadString('\n')
	input = strings.ReplaceAll(strings.TrimSpace(input), "\\", "")
	if input == "" {
		input = defaultPath
	}

	info, err := os.Stat(input)

	switch {
	case os.IsNotExist(err):
		fmt.Println("📁 Directory does not exist – trying to create it…")
		if errMkdir := os.MkdirAll(input, 0o700); errMkdir != nil {
			return fmt.Errorf("❌ failed to create directory: %w", errMkdir)
		}
		fmt.Println("✅ Directory created:", input)
	case err != nil:
		return fmt.Errorf("❌ failed to access path: %w", err)
	case !info.IsDir():
		return fmt.Errorf("❌ path exists but is not a directory: %s", input)
	}

	configuration := config.Config{
		StoragePath:   input,
		ActivePackage: config.DefaultPackage,
	}
	if errSave := configuration.Save(); errSave != nil {
		return fmt.Errorf("failed to save config: %w", errSave)
	}
	fmt.Printf("✅ Configuration saved to %s\n", config.Path())
	return nil
}
