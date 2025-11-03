package files

import (
	"fmt"
	"os"
)

func CreateFile(message string) {
	fmt.Println("This is will save to file")
	if err := os.WriteFile("files/example.txt", []byte(message), 0644); err != nil {
		fmt.Printf("failed to saved file: %v", err)
	}
}

func ReadFile() string {
	data, err := os.ReadFile("files/example.txt")
	if err != nil {
		fmt.Printf("failed to read file: %v", err)
	}
	return string(data)
}
