package server

import (
	"../formData"
	"../query"
	"log"
	"net"
	"time"
)

func backgroundRead(conn net.Conn, c chan string, connClosed chan bool) {
	for {
		buf := make([]byte, 512)
		nr, err := conn.Read(buf)
		if err != nil {
			connClosed <- true
			return
		}
		c <- string(buf[0:nr])
	}
}

func requestHandler(conn net.Conn) {
	rx := make(chan string)
	connClosed := make(chan bool)
	go backgroundRead(conn, rx, connClosed)
	send := func(msg []byte) {
		if _, err := conn.Write(msg); err != nil {
			log.Fatal("Write: ", err)
		}
	}
	for {
		//set timeout
		timeout := make(chan bool, 1)
		go func() {
			time.Sleep(10 * time.Second)
			timeout <- true
		}()

		//send prefix
		send([]byte("> "))

		select {
		case query := <-rx:
			response := handleQuery(query)
			sendData := formData.FormData{formData.OK, response}
			send(formData.Serialize(sendData))
		case <-timeout:
			//send timeout message
			sendData := formData.FormData{formData.CLOSE, ""}
			send(formData.Serialize(sendData))
			println("timeout")
			return
		case <-connClosed:
			return
		}
	}
}

func handleQuery(queryStr string) string {
	q := query.Parse(queryStr)
	switch q.Op {
	case query.GET:
		if len(q.Args) == 0 {
			return "Error\n"
		} else {
			return get(q.Args[0]) + "\n"
		}
	case query.SET:
		if len(q.Args) <= 1 {
			return "Error\n"
		} else {
			key, value := q.Args[0], q.Args[1]
			set(key, value)
			return "OK\n"
		}
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
