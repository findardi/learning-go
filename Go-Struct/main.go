package main

import (
	"fmt"
	"go-struct/user"
)

func main() {
	user, err := user.New("Ardi", "", "Luckey")
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	user.PrintOut()
	user.ClearUsername()
	user.PrintOut()
}
