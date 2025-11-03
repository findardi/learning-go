package main

import (
	"bufio"
	"fmt"
	"go-esentials/files"
	"os"
)

func main() {
	fmt.Println("Handle file simple")
	message := bufio.NewReader(os.Stdin)
	fmt.Print("Masukan teks: ")
	teks, _ := message.ReadString('\n')

	files.CreateFile(teks)
	fmt.Println("Hasil:", files.ReadFile())
}
