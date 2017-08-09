package main

import (
	"./client"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		endpoint := os.Args[1]
		client.ClientMain(os.Stdin, endpoint)
	} else {
		fmt.Printf("Usage: %s <ip-address:port>\n", os.Args[0])
	}
}
