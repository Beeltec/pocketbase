//go:build !windows

package cmd

import (
	"github.com/fatih/color"
	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
)

func NewServiceCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:   "service",
		Short: "Manages Windows service registration",
		RunE: func(command *cobra.Command, args []string) error {
			color.Yellow("This command can only be run on Windows")
			return nil
		},
	}

	return command
}
