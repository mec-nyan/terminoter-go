package notes

import (
	"encoding/json"
	"fmt"
	"os"
)

type Note struct {
	// The text of the note.
	Content string `json:"content"`
	// To be used to save some data (i.e. date and time, etc).
	Metadata map[string]any `json:"metadata"`
}

type Data struct {
	Notes []Note `json:"notes"`
}

func LoadNotes(file string) (*Data, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to load file %s: %w", file, err)
	}

	var data Data
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &data, nil
}

func SaveNotes(data *Data, file string) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(file, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

type Loader struct {
	Data
	Err error
}
