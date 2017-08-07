package main

import (
	"./client"
	"os"
)

func main() {
	endpoint := "/tmp/tmp.sock"
	if len(os.Args) > 1 {
		endpoint = os.Args[1]
	}
	client.ClientMain(os.Stdin, endpoint)
}
