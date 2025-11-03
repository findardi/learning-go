package main

import "fmt"

func main() {
	number := 32
	fmt.Printf("original value: %d, address: %p\n", number, &number)

	noPointer := addNumber(number)
	fmt.Printf("No Pointer (result): %d\n", noPointer)

	pointer := addNumberPointer(&number)
	fmt.Printf("Ponter (result): %d\n", pointer)
}

func addNumber(number int) int {
	fmt.Printf("no pointer address: %p\n", &number)
	return number + 10
}

func addNumberPointer(number *int) int {
	fmt.Printf("pointer address: %p\n", &number)
	return *number + 10
}
