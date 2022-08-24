package main

import (
	"fmt"
	"os"

	"github.com/abhi-g80/chipku/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Printf("main: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}
