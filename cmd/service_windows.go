//go:build windows

package cmd

import (
	"errors"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"

	"golang.org/x/sys/windows/svc/mgr"
)

func toCamelCase(s string) string {
	words := strings.Fields(s)
	for i := 0; i < len(words); i++ {
		words[i] = strings.Title(words[i])
	}
	return strings.Join(words, "")
}

func isWindowsAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	return true
}

func NewServiceCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:   "service",
		Short: "Manages Windows service registration",
	}

	command.AddCommand(serviceRegisterCommand(app))
	command.AddCommand(serviceRemovalCommand(app))

	return command
}

func serviceRegisterCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:   "register",
		Short: "Registers this pocketbase application as a Windows service",
		RunE: func(command *cobra.Command, args []string) error {
			serviceName := toCamelCase(app.Settings().Meta.AppName)
			serviceDisplayName := app.Settings().Meta.AppName
			serviceDescription := "Windows service for " + app.Settings().Meta.AppName + " application"

			if !isWindowsAdmin() {
				return errors.New("You need to run this command as an administrator")
			}

			exePath, err := os.Executable()
			if err != nil {
				return errors.New("Failed to get executable path")
			}

			exePath += " serve" // Add "serve" parameter to the executable path

			m, err := mgr.Connect()
			if err != nil {
				return errors.New("Failed to connect to service manager")
			}
			defer m.Disconnect()

			service, err := m.OpenService(serviceName)
			if err == nil {
				service.Close()
				return errors.New("Service already exists")
			}

			service, err = m.CreateService(serviceName, exePath, mgr.Config{
				DisplayName: serviceDisplayName,
				Description: serviceDescription,
				StartType:   mgr.StartAutomatic,
			})
			if err != nil {
				return errors.New("Failed to create service")
			}
			defer service.Close()

			color.Green("Service registered successfully")
			return nil
		},
	}

	return command
}

func serviceRemovalCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:   "remove",
		Short: "Removes the Windows service registration for this pocketbase application",
		RunE: func(command *cobra.Command, args []string) error {
			serviceName := toCamelCase(app.Settings().Meta.AppName)

			if !isWindowsAdmin() {
				return errors.New("You need to run this command as an administrator")
			}

			m, err := mgr.Connect()
			if err != nil {
				return errors.New("Failed to connect to service manager")
			}
			defer m.Disconnect()

			service, err := m.OpenService(serviceName)
			if err != nil {
				return errors.New("Service does not exist")
			}
			defer service.Close()

			err = service.Delete()
			if err != nil {
				return errors.New("Failed to remove service")
			}

			color.Green("Service removed successfully")
			return nil
		},
	}

	return command
}
