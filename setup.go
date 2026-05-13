package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func setup(opts *Options) error {
	if opts.file == "" {
		opts.file = defaultSaveLocation()
	}

	err := ensureFileExists(opts.file)
	if err != nil {
		return err
	}

	return nil
}

func defaultSaveLocation() string {
	file := "terminoter.json"
	app := "terminoter"

	dir := os.Getenv("XDG_DATA_HOME")
	if dir == "" {
		home := os.Getenv("HOME")
		if home == "" {
			return filepath.Join(".", file)
		}
		dir = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(dir, app, file)
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
		data := Data{}

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
