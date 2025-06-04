package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/asperling/go-todo-cli/config"
	"github.com/urfave/cli/v2"
)

func InitAction(_ *cli.Context, store *config.Store) error {
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
			return cli.Exit(fmt.Sprintf("❌ could not create storage directory: %v", errMkdir), 1)
		}
		fmt.Println("✅ Directory created:", input)
	case err != nil:
		return cli.Exit(fmt.Sprintf("❌ failed to access path: %v", err), 1)
	case !info.IsDir():
		return cli.Exit(fmt.Sprintf("❌ path exists but is not a directory: %s", input), 1)
	}

	configuration := config.Config{
		StoragePath:   input,
		ActivePackage: config.DefaultPackage,
	}

	if errSave := store.Save(&configuration); errSave != nil {
		return cli.Exit(fmt.Sprintf("failed to save config: %v", errSave), 1)
	}
	fmt.Printf("✔ Configuration saved to %s\n", store.FilePath)
	return nil
}

func InitCommand(store *config.Store) *cli.Command {
	return &cli.Command{
		Name:        "init",
		Aliases:     []string{"i"},
		Usage:       "Initialize the todo CLI configuration",
		Description: "This command initializes the todo CLI by setting up the storage path for todos.",
		Action: func(c *cli.Context) error {
			return InitAction(c, store)
		},
	}
}
