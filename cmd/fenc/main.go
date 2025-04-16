package main

import (
	"fmt"
	"log"

	f "github.com/nielsjaspers/fenc/internal/filehandling"
)

func main() {
	fmt.Println("fenc - file encryption & compression")
	_, err := f.OpenFile("~/Desktop/test.md")
	if err != nil {
		log.Fatalf("Error while opening file: %v\n", err)
	}
}
