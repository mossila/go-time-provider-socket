package main

import (
	"fmt"
	"os"

	"github.com/mossila/go-time-provider-socket/provider"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage\ngo run main.go 1234")
		return
	}
	port := os.Args[1]
	provider.TimeProvider(port)
}
