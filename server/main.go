package main

import (
	"./malScan"
	"./server"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		malScan.ConstructRules()
		server.Serve("unix", os.Args[1])
	} else {
		fmt.Printf("Usage: %s <endpoint>", os.Args[0])
	}
}
