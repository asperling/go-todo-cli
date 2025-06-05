package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func Exit(message string) error {
	return cli.Exit(fmt.Sprintf("❌ %s", message), 1)
}

func Exitf(format string, args ...any) error {
	return cli.Exit(fmt.Sprintf("❌ "+format, args...), 1)
}

func SuccessPrint(message string) {
	fmt.Printf("✔ %s\n", message)
}

func SuccessPrintf(format string, args ...any) {
	fmt.Printf("✔ "+format+"\n", args...)
}
