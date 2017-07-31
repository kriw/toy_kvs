package main

import (
	"./server"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		server.Serve("unix", os.Args[1])
	} else {
		fmt.Printf("Usage: %s <endpoint>", os.Args[0])
	}
}
