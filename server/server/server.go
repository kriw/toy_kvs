package server

import (
	"log"
	"net"
	"strings"
)

func RequestHandler(c net.Conn) {
	for {
		if _, err := c.Write([]byte("> ")); err != nil {
			log.Fatal("Write: ", err)
		}
		buf := make([]byte, 512)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}
		query := buf[0:nr]
		response := handleQuery(string(query))
		if _, err := c.Write([]byte(response)); err != nil {
			log.Fatal("Write: ", err)
		}
	}
}

func handleQuery(query string) string {
	s := strings.Split(query, " ")
	op, arg := s[0], s[1:]
	switch op {
	case "get":
		return get(strings.TrimSpace(arg[0])) + "\n"
	case "set":
		key, value := strings.TrimSpace(arg[0]), strings.TrimSpace(arg[1])
		set(key, value)
		return "OK\n"
	default:
		return "Unknown query.\n"
	}
}
