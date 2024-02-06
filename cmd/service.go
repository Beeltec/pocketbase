//go:build !windows

package cmd

import (
	"os"
	"runtime"

	"github.com/fatih/color"
	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
)

func NewServiceCommand(app core.App) *cobra.Command {
	if runtime.GOOS != "windows" {
		color.Yellow("This command can only run on Windows")
		os.Exit(1)
	}

	command := &cobra.Command{
		Use:   "service",
		Short: "Manages Windows service registration",
	}

	return command
}
