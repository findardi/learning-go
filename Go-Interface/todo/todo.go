package todo

import (
	"encoding/json"
	"fmt"
	"os"
)

type Todo struct {
	Text string `json:"text"`
}

func New(text string) *Todo {
	return &Todo{
		Text: text,
	}
}

func (t *Todo) Save() error {
	filename := "todo.json"

	result, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("failed to marshaling :%w", err)
	}

	if err := os.WriteFile(filename, result, 0644); err != nil {
		return fmt.Errorf("failed to save file :%w", err)
	}

	return nil
}
