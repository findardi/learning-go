package note

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Note struct {
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

func New(text string) *Note {
	return &Note{
		Text:      text,
		CreatedAt: time.Now(),
	}
}

func (n *Note) Save() error {
	filename := "note.json"

	result, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("failed to marshaling :%w", err)
	}

	if err := os.WriteFile(filename, result, 0644); err != nil {
		return fmt.Errorf("failed to save file :%w", err)
	}

	return nil
}
