package main

import (
	"fmt"

	"github.com/hlpd-pham/chirpy/server"
)

func main() {
	s := server.NewServer()
	err := s.Start()
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
