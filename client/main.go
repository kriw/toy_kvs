package main

import (
	"../tkvs_protocol"
	"./query"
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func readUsr(s *bufio.Scanner, iCh chan string, isClosed chan bool) {
	for s.Scan() {
		if err := s.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
		iCh <- s.Text()
	}
	isClosed <- true
}

func readSrv(r io.Reader, srvInput chan string, isClosed chan bool) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			return
		}
		res := tkvs_protocol.Deserialize(buf[0:n])
		switch res.Method {
		case tkvs_protocol.CLOSE:
			srvInput <- "Connection has been closed"
			isClosed <- true
		case tkvs_protocol.OK:
			if res.Data == "" {
				srvInput <- "OK"
			} else {
				srvInput <- res.Data
			}
		case tkvs_protocol.ERROR:
			srvInput <- "ERROR"
		}
	}
}

func handleQuery(queryStr string) tkvs_protocol.Protocol {
	q := query.Parse(queryStr)
	switch q.Op {
	case query.GET:
		return tkvs_protocol.Protocol{tkvs_protocol.GET, q.Args[0]}
	case query.SET:
		key, value := q.Args[0], q.Args[1]
		return tkvs_protocol.Protocol{tkvs_protocol.SET, key + "," + value}
	default:
		return tkvs_protocol.Protocol{}
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	c, err := net.Dial("unix", "/tmp/echo.sock")
	isClosed := make(chan bool)
	srvInput := make(chan string)
	usrInput := make(chan string)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	go readSrv(c, srvInput, isClosed)
	go readUsr(scanner, usrInput, isClosed)

	for {
		print("> ")
		select {
		case <-isClosed:
			return
		case input := <-srvInput:
			println(input)
		case input := <-usrInput:
			q := handleQuery(input)
			p := tkvs_protocol.Serialize(q)
			if _, err := c.Write(p); err != nil {
				log.Fatal("write error:", err)
				break
			}
		}
	}
}
