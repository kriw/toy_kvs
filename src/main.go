package main

import (
	"./server"
	"log"
	"net"
	"strings"
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
		go server.KeyValueServer(fd, proc)
	}
}

func proc(query string) string {
	s := strings.Split(query, " ")
	op, arg := s[0], s[1:]
	switch op {
	case "get":
		return server.Get(strings.TrimSpace(arg[0])) + "\n"
	case "set":
		key, value := strings.TrimSpace(arg[0]), strings.TrimSpace(arg[1])
		server.Set(key, value)
		return "OK\n"
	default:
		return "Unknown query.\n"
	}
}
