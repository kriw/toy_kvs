package main

import (
	"../tkvs_protocol"
	"./query"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func readInputFromUsr(s *bufio.Scanner, iCh chan string, isClosed chan bool) {
	for s.Scan() {
		if err := s.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
		iCh <- s.Text()
	}
	isClosed <- true
}

func readMsgFromSrv(r io.Reader, srvInput chan string, isClosed chan bool) {
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
			if len(res.Data) == 0 {
				srvInput <- "OK"
			} else {
				srvInput <- string(res.Data)
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
		if len(q.Args) == 1 {
			key, data := q.Args[0], make([]byte, 0)
			return tkvs_protocol.Protocol{tkvs_protocol.GET, key, data}
		}
	case query.SET:
		if len(q.Args) == 2 {
			key, data := q.Args[0], q.Args[1]
			return tkvs_protocol.Protocol{tkvs_protocol.SET, key, data}
		}
	case query.SETFILE:
		if len(q.Args) == 2 {
			key, filename := q.Args[0], string(q.Args[1])
			if filedata, err := ioutil.ReadFile(filename); err == nil {
				return tkvs_protocol.Protocol{tkvs_protocol.SET, key, filedata}
			}
		}
	}
	b := make([]byte, 0)
	return tkvs_protocol.Protocol{tkvs_protocol.ERROR, b, b}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	c, err := net.Dial("unix", "/tmp/echo.sock")
	if err != nil {
		panic(err)
	}
	isClosed := make(chan bool)
	srvInput := make(chan string)
	usrInput := make(chan string)
	defer c.Close()

	go readMsgFromSrv(c, srvInput, isClosed)
	go readInputFromUsr(scanner, usrInput, isClosed)

	for {
		print("> ")
		select {
		case <-isClosed:
			return
		case input := <-srvInput:
			println(input)
		case input := <-usrInput:
			if q := handleQuery(input); q.Method == tkvs_protocol.ERROR {
				println("Error Input: " + input)
			} else {
				p := tkvs_protocol.Serialize(q)
				if _, err := c.Write(p); err != nil {
					log.Fatal("write error:", err)
					break
				}
			}
		}
	}
}
