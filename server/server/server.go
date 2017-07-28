package server

import (
	"../../tkvs_protocol"
	"../../util"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

const BUF_SIZE = 1024 * 1024 * 1024

var database = make(map[[util.HashSize]byte][]byte)

func save(filename string) {
	toBytes := func(data map[[util.HashSize]byte][]byte) []byte {
		b := bytes.Buffer{}
		e := gob.NewEncoder(&b)
		if err := e.Encode(data); err != nil {
			fmt.Println(`failed gob Encode`, err)
		}
		return b.Bytes()
	}

	content := toBytes(database)
	ioutil.WriteFile(filename, content, os.ModePerm)
}

func get(key [util.HashSize]byte) []byte {
	return database[key]
}

func set(key [util.HashSize]byte, value []byte) {
	database[key] = value
}

func backgroundRead(conn net.Conn, c chan []byte, connClosed chan bool) {
	buf := make([]byte, BUF_SIZE)
	for {
		nr, err := conn.Read(buf)
		if err != nil {
			connClosed <- true
			return
		}
		c <- buf[0:nr]
	}
}

func requestHandler(conn net.Conn) {
	rx := make(chan []byte)
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
			time.Sleep(100 * time.Second)
			timeout <- true
		}()

		select {
		case rawReq := <-rx:
			req := tkvs_protocol.Deserialize(rawReq)
			response := handleReq(req)
			send(tkvs_protocol.Serialize(response))
		case <-timeout:
			//send timeout message
			empKey := [util.HashSize]byte{}
			empData := make([]byte, 0)
			sendData := tkvs_protocol.Protocol{tkvs_protocol.CLOSE, empKey, empData}
			send(tkvs_protocol.Serialize(sendData))
			println("timeout")
			return
		case <-connClosed:
			return
		}
	}
}

func handleReq(req tkvs_protocol.Protocol) tkvs_protocol.Protocol {
	empKey := [util.HashSize]byte{}
	empData := make([]byte, 0)
	method := req.Method
	switch method {
	case tkvs_protocol.GET:
		res := get(req.Key)
		return tkvs_protocol.Protocol{tkvs_protocol.OK, empKey, res}
	case tkvs_protocol.SET:
		if hashedData := sha256.Sum256(req.Data); hashedData == req.Key {
			fmt.Printf("%x", hashedData)
			set(req.Key, req.Data)
			return tkvs_protocol.Protocol{tkvs_protocol.OK, empKey, empData}
		}
	case tkvs_protocol.SAVE:
		save(string(req.Data))
		return tkvs_protocol.Protocol{tkvs_protocol.OK, empKey, empData}
	case tkvs_protocol.CLOSE:
		return tkvs_protocol.Protocol{tkvs_protocol.CLOSE, empKey, empData}
	}
	return tkvs_protocol.Protocol{tkvs_protocol.ERROR, empKey, empData}
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
