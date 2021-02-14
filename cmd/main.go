package main

import (
	"fmt"

	radish "github.com/Artamus/radish"
)

func main() {
	storage := make(map[string]string)

	server, err := radish.NewServer(6379, storage)
	if err != nil {
		fmt.Printf("Failed to start Radish server, %v", err)
	}

	server.Listen()
}
