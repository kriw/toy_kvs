package server

import (
	"log"
	"net"
	"strings"
	"time"
)

func receiver(conn net.Conn, c chan string) {
	for {
		buf := make([]byte, 512)
		nr, err := conn.Read(buf)
		if err != nil {
			log.Fatal("Read: ", err, "\n", "Connection closed")
			return
		}
		c <- string(buf[0:nr])
	}
}

func RequestHandler(conn net.Conn) {
	rx := make(chan string)
	go receiver(conn, rx)
	for {
		//set timeout
		timeout := make(chan bool, 1)
		go func() {
			time.Sleep(10 * time.Second)
			timeout <- true
		}()

		//send prefix
		if _, err := conn.Write([]byte("> ")); err != nil {
			log.Fatal("Write: ", err)
			return
		}

		select {
		case query := <-rx:
			response := handleQuery(query)
			if _, err := conn.Write([]byte(response)); err != nil {
				log.Fatal("Write: ", err)
				return
			}
		case <-timeout:
			println("timeout")
			return
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
