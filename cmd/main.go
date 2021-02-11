package main

import (
	"fmt"

	radish "github.com/Artamus/radish"
)

func main() {
	server, err := radish.NewRadishServer(6379)
	if err != nil {
		fmt.Printf("Failed to start Radish server, %v", err)
	}

	server.Listen()
}
