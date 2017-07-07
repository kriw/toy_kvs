package server

import (
	"../formData"
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
			return
		}
		c <- string(buf[0:nr])
	}
}

func requestHandler(conn net.Conn) {
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
			sendData := formData.FormData{formData.OK, response}
			if _, err := conn.Write(formData.Serialize(sendData)); err != nil {
				return
			}
		case <-timeout:
			//send timeout message
			sendData := formData.FormData{formData.CLOSE, ""}
			if _, err := conn.Write(formData.Serialize(sendData)); err != nil {
				return
			}
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

func Serve(connType, laddr string) {
	l, err := net.Listen(connType, laddr)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	for {
		fd, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}
		go requestHandler(fd)
	}
}
