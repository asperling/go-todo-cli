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
			return Exitf("Could not create storage directory: %v", errMkdir)
		}
		SuccessPrintf("Directory created: %s", input)
	case err != nil:
		return Exitf("Failed to access path: %v", err)
	case !info.IsDir():
		return Exitf("Path exists but is not a directory: %s", input)
	}

	configuration := config.Config{
		StoragePath:   input,
		ActivePackage: config.DefaultPackage,
	}

	if errSave := store.Save(&configuration); errSave != nil {
		return Exitf("Failed to save config: %v", errSave)
	}
	SuccessPrintf("Configuration saved to %s", store.FilePath)
	return nil
}

func InitCommand(store *config.Store) *cli.Command {
	return &cli.Command{
		Name:        "init",
		Aliases:     []string{"i"},
		Usage:       "Initialize the todo CLI configuration",
		Description: "This command initializes the todo CLI by setting up the storage path for todos. This is usually the first command you run.",
		ArgsUsage:   "",
		Action: func(c *cli.Context) error {
			return InitAction(c, store)
		},
	}
}
