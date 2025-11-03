package main

import (
	"bufio"
	"fmt"
	"go-interface/note"
	"go-interface/todo"
	"os"
)

type saver interface {
	Save() error
}

func main() {
	todoText := userInput("input todo: ")
	noteText := userInput("input note: ")

	todo := todo.New(todoText)
	note := note.New(noteText)

	if err := saveDate(todo); err != nil {
		return
	}

	if err := saveDate(note); err != nil {
		return
	}

	anyValue(100)
	anyValue("Nothinggg")
}

func saveDate(data saver) error {
	if err := data.Save(); err != nil {
		fmt.Printf("%s", err.Error())
		return err
	}

	fmt.Println("Success Saved")
	return nil
}

func userInput(prompt string) string {
	fmt.Printf("%v ", prompt)
	reader := bufio.NewReader(os.Stdin)
	teks, _ := reader.ReadString('\n')
	return teks
}

func anyValue(value interface{}) {
	switch value.(type) {
	case int:
		fmt.Println("int value: ", value)
	case string:
		fmt.Println("string value: ", value)
	default:
		fmt.Println("another value type: ", value)
	}
}
