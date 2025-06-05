package commands

import (
	"regexp"

	"github.com/asperling/go-todo-cli/config"
	"github.com/urfave/cli/v2"
)

func isValidPackageName(name string) bool {
	matched, _ := regexp.MatchString(`^[A-Za-z0-9]+$`, name)
	return matched
}

func PackagesCommand(store *config.Store) *cli.Command {
	return &cli.Command{
		Name:        "packages",
		Usage:       "Manage todo packages. See `todo packages help` for more information.",
		Description: "Manages your task packages. Packages are named sets of todos for different topics or tickets.",
		ArgsUsage:   "",
		Aliases:     []string{"pkg", "p"},
		Subcommands: []*cli.Command{
			PackagesListCommand(store),
			PackagesUseCommand(store),
			PackagesDeleteCommand(store),
		},
		Action: func(_ *cli.Context) error {
			// Default action: show list
			return PackagesListAction(store)
		},
	}
}
