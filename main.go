package main

import (
	"gator/internal/config"
	"fmt"
)

func main() {
	config, err := Read()
	if err != nil {
		fmt.Printf("Failed to read the config file")
		return
	}

	config.SetUser("chandan")

	config, err = Read()
	if err != nil {
		fmt.Printf("Failed to update the config")
		return
	}

	fmt.Print(config)
}

