package client

import (
	"../../tkvsProtocol"
	"../../util"
	"../query"
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
)

const BUF_SIZE = 20 * 1024 * 1024

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
	buf := make([]byte, BUF_SIZE)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			return
		}
		res := tkvsProtocol.Deserialize(buf[0:n])
		switch res.Method {
		case tkvsProtocol.CLOSE:
			srvInput <- "Connection has been closed"
			isClosed <- true
		case tkvsProtocol.OK:
			if len(res.Data) == 0 {
				srvInput <- "OK"
			} else {
				srvInput <- string(res.Data)
			}
		case tkvsProtocol.ERROR:
			srvInput <- "ERROR"
		}
	}
}

func checkFileSize(data []byte) bool {
	return len(data) <= BUF_SIZE
}

func handleQuery(queryStr string) tkvsProtocol.Protocol {
	q := query.Parse(queryStr)
	switch q.Op {
	case query.GET:
		if len(q.Args) == 1 {
			if key, err := hex.DecodeString(string(q.Args[0])); err == nil {
				var key32bit [util.HashSize]byte
				copy(key32bit[:], key)
				return tkvsProtocol.Protocol{tkvsProtocol.GET, key32bit, make([]byte, 0)}
			}
		}
	case query.SET:
		if len(q.Args) == 1 {
			filename := string(q.Args[0])
			if filedata, err := ioutil.ReadFile(filename); err == nil {
				key := sha256.Sum256(filedata)
				fmt.Printf("key: %x\n", key)
				return tkvsProtocol.Protocol{tkvsProtocol.SET, key, filedata}
			}
		}
	case query.SAVE:
		if len(q.Args) == 1 {
			filename := q.Args[0]
			return tkvsProtocol.Protocol{tkvsProtocol.SAVE, [util.HashSize]byte{}, filename}
		}
	}

	return tkvsProtocol.Protocol{tkvsProtocol.ERROR, [util.HashSize]byte{}, make([]byte, 0)}
}

func ClientMain(r io.Reader, endpoint string) {
	scanner := bufio.NewScanner(r)

	c, err := net.Dial("unix", endpoint)
	defer c.Close()
	if err != nil {
		panic(err)
	}

	isClosed := make(chan bool)
	srvInput := make(chan string)
	usrInput := make(chan string)

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
			if q := handleQuery(input); q.Method == tkvsProtocol.ERROR {
				println("Error Input: " + input)
			} else {
				p := tkvsProtocol.Serialize(q)
				if _, err := c.Write(p); err != nil {
					log.Fatal("write error:", err)
					break
				}
			}
		}
	}
}
