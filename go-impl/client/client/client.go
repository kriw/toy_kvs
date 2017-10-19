package client

import (
	"../../proto"
	"../../tkvsProtocol"
	"../../util"
	"../query"
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
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
	buf := make([]byte, util.BUF_SIZE)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			return
		}
		res := tkvsProtocol.DeserializeRes(buf[0:n])
		switch res.Response {
		case proto.TIMEOUT:
			srvInput <- "Timeout: Connection has been closed"
			isClosed <- true
		case proto.FILEEXIST:
			fallthrough
		case proto.SUCCESS:
			if len(res.Data) == 0 {
				srvInput <- "OK"
			} else {
				srvInput <- string(res.Data)
			}
		case proto.ERROR:
			srvInput <- fmt.Sprintf("ERROR: %s", res.Data)
		}
	}
}

func checkFileSize(data []byte) bool {
	return len(data) <= util.BUF_SIZE
}

func handleQuery(queryStr string) proto.RequestParam {
	q := query.Parse(queryStr)
	switch q.Op {
	case query.GET:
		if len(q.Args) == 1 {
			if key, err := hex.DecodeString(string(q.Args[0])); err == nil {
				var key32bit [proto.HashSize]byte
				copy(key32bit[:], key)
				return proto.RequestParam{proto.GET, 0, key32bit, make([]byte, 0)}
			}
		}
	case query.SET:
		if len(q.Args) == 1 {
			filename := string(q.Args[0])
			if filedata, err := util.ReadFile(filename); err == nil {
				key := sha256.Sum256(filedata)
				fmt.Printf("key: %x\n", key)
				return proto.RequestParam{proto.SET, uint64(len(filedata)), key, filedata}
			}
		}
	}

	return proto.RequestParam{proto.GET, 0, [proto.HashSize]byte{}, make([]byte, 0)}
}

func ClientMain(r io.Reader, endpoint string) {
	scanner := bufio.NewScanner(r)

	c, err := net.Dial("tcp", endpoint)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	isClosed := make(chan bool)
	srvInput := make(chan string)
	usrInput := make(chan string)

	go readMsgFromSrv(c, srvInput, isClosed)
	go readInputFromUsr(scanner, usrInput, isClosed)

	for {
		print("> ")
		select {
		case <-isClosed:
			q := proto.RequestParam{proto.CLOSE_CLI, 0, [proto.HashSize]byte{}, make([]byte, 0)}
			p := tkvsProtocol.SerializeReq(q)
			_, _ = c.Write(p)
			return
		case input := <-usrInput:
			//FIXME
			if q := handleQuery(input); q.Method == proto.ERROR_INPUT {
				println("Error Input: " + input)
			} else {
				p := tkvsProtocol.SerializeReq(q)
				if _, err := c.Write(p); err != nil {
					log.Fatal("write error:", err)
					break
				}
			}
			res := <-srvInput
			println(res)
		}
	}
}
