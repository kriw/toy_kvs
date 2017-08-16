package main

import (
	"./server"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		fmt.Printf("Listening: %s\n", os.Args[1])
		server.Serve("tcp", os.Args[1])
	} else {
		fmt.Printf("Usage: %s <ip-address:port>\n", os.Args[0])
	}
}
