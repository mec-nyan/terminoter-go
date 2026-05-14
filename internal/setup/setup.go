package setup

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/mec-nyan/terminoter-go/internal/args"
	"github.com/mec-nyan/terminoter-go/internal/notes"
)

const (
	AppName                  = "terminoter"
	DefaultConfigurationFile = "config.toml"
)

func Setup(opts *args.Options) error {
	// Required: file to persist our notes.
	if opts.File == "" {
		opts.File = defaultSaveLocation()
	}

	err := ensureFileExists(opts.File)
	if err != nil {
		return err
	}

	// Optional: configuration file.
	// If it doesn't exist, we'll use the default configuration.
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = filepath.Join(os.Getenv("HOME"), ".config")
	}

	configFile := filepath.Join(dir, AppName, "config.toml")

	if _, err := os.Stat(configFile); err != nil {
		if os.IsNotExist(err) {
			// Fallback to default configuration.
			return nil
		} else {
			return err
		}
	} else {
		// We've got a config file.
		content, err := os.ReadFile(configFile)
		if err != nil {
			return err
		}

		if len(content) == 0 {
			// Fallback to default configuration.
			return nil
		}

		// TODO: Parse config file and load options.
	}

	return nil
}

func defaultSaveLocation() string {
	dir := os.Getenv("XDG_DATA_HOME")
	if dir == "" {
		home := os.Getenv("HOME")
		if home == "" {
			return filepath.Join(".", AppName+".json")
		}
		dir = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(dir, AppName, AppName+".json")
}

func ensureFileExists(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}

	// If file doesn't exist, create it.
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o600)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	// If the file is empty, initialise it.
	if info.Size() == 0 {
		data := notes.Data{}

		content, err := json.Marshal(data)
		if err != nil {
			return err
		}

		err = os.WriteFile(path, content, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
