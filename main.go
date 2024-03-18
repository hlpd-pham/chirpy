package main

import (
	"fmt"
)

func main() {
	server := NewServer()
	err := server.Start()
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
