package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ksysoev/wsget/pkg/repo/macro"
	"github.com/spf13/cobra"
)

func newUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [macro]",
		Short: "Update macro files from their sources",
		Long: `Update command checks for updates in installed macro files and updates them if newer versions are available.
For example:
  wsget update            # Updates all installed macro files
  wsget update [macro]   # Updates a specific macro file`,
		RunE: runUpdate,
	}

	return cmd
}

func runUpdate(cmd *cobra.Command, args []string) error {
	configDir, err := cmd.Root().PersistentFlags().GetString("config-dir")
	if err != nil {
		return fmt.Errorf("failed to get config directory flag: %w", err)
	}

	// If configDir is not set, use default location
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		configDir = filepath.Join(homeDir, "wsget")
	}

	macroDir := filepath.Join(configDir, "macro")
	if err := os.MkdirAll(macroDir, 0755); err != nil {
		return fmt.Errorf("failed to create macro directory: %w", err)
	}

	// If a specific macro is provided, update only that macro
	if len(args) > 0 {
		macroName := args[0]
		return updateMacro(macroDir, macroName)
	}

	// Otherwise, update all macros
	return updateAllMacros(macroDir)
}

func updateMacro(macroDir, macroName string) error {
	macroPath, err := macro.FindMacro(macroDir, macroName)
	if err != nil {
		return err
	}

	fmt.Printf("Updating macro: %s\n", macroName)
	if err := macro.Update(macroPath); err != nil {
		return fmt.Errorf("failed to update macro %s: %w", macroName, err)
	}

	fmt.Printf("Successfully updated macro: %s\n", macroName)
	return nil
}

func updateAllMacros(macroDir string) error {
	fmt.Println("Updating all macros...")
	if err := macro.UpdateAll(macroDir); err != nil {
		return fmt.Errorf("failed to update macros: %w", err)
	}
	fmt.Println("Successfully updated all macros")
	return nil
}
