package main

import (
	"./server"
)

func main() {
	server.Serve("unix", "/tmp/echo.sock")
}
