package main

import (
	"./server"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("unix", "/tmp/echo.sock")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	for {
		fd, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}
		go server.RequestHandler(fd)
	}
}
