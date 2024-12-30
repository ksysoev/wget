package macro

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// Update updates a macro configuration file from its source URL.
// It takes filepath of type string as input.
// It returns an error if the update fails, the file doesn't exist, the source URL is missing,
// the YAML unmarshalling fails, or the macro version is unsupported.
func Update(filepath string) error {
	// Read existing config
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("fail to open macro file: %w", err)
	}
	defer file.Close()

	cfg, err := newConfig(file)
	if err != nil {
		return fmt.Errorf("fail to read macro config: %w", err)
	}

	sourceURL := cfg.Source
	if sourceURL == "" {
		return fmt.Errorf("macro file has no source URL")
	}

	// Download latest version
	resp, err := http.Get(sourceURL) //nolint:gosec // This is a CLI tool, and the URL is from the config file
	if err != nil {
		return fmt.Errorf("fail to download macro update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fail to download macro update: %s", resp.Status)
	}

	newCfg, err := newConfig(resp.Body)
	if err != nil {
		return fmt.Errorf("fail to read updated macro: %w", err)
	}

	newCfg.SetSource(sourceURL)

	if _, err := newCfg.CreateRepo(); err != nil {
		return fmt.Errorf("fail to create commands: %w", err)
	}

	// Write updated config
	file, err = os.Create(filepath)
	if err != nil {
		return fmt.Errorf("fail to create file: %w", err)
	}

	defer func() {
		if e := file.Close(); err == nil && e != nil {
			err = fmt.Errorf("fail to close file: %w", e)
		}
	}()

	if err := newCfg.Write(file); err != nil {
		return fmt.Errorf("fail to write updated macro: %w", err)
	}

	return nil
}

// UpdateAll updates all macro configuration files in the specified directory.
// It takes macroDir of type string as input.
// It returns an error if reading the directory fails, but continues updating other files if one fails.
func UpdateAll(macroDir string) error {
	entries, err := os.ReadDir(macroDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No macro files yet
		}
		return fmt.Errorf("failed to read macro directory: %w", err)
	}

	var updateErrors []error
	for _, entry := range entries {
		if !entry.IsDir() && (filepath.Ext(entry.Name()) == ".yml" || filepath.Ext(entry.Name()) == ".yaml") {
			macroPath := filepath.Join(macroDir, entry.Name())
			if err := Update(macroPath); err != nil {
				updateErrors = append(updateErrors, fmt.Errorf("failed to update %s: %w", entry.Name(), err))
				continue
			}
		}
	}

	if len(updateErrors) > 0 {
		var errMsg string
		for _, err := range updateErrors {
			errMsg += err.Error() + "\n"
		}
		return fmt.Errorf("some updates failed:\n%s", errMsg)
	}

	return nil
}

// FindMacro looks for a macro file with the given name in the specified directory.
// It tries both .yml and .yaml extensions.
// It returns the full path to the found macro file and nil if found, or empty string and error if not found.
func FindMacro(macroDir, macroName string) (string, error) {
	extensions := []string{".yml", ".yaml"}
	for _, ext := range extensions {
		path := filepath.Join(macroDir, macroName+ext)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("macro %s not found", macroName)
}
