package server

import (
	"../../tkvsProtocol"
	"../../util"
	"../malScan"
	"../scanLog"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

const BUF_SIZE = 1024 * 1024 * 1024
const FILE_DIR = "./files/"

var fileHashMap = make(map[[util.HashSize]byte]bool)

func save(filename string, fileContent []byte) {
	toBytes := func(data []byte) []byte {
		b := bytes.Buffer{}
		e := gob.NewEncoder(&b)
		if err := e.Encode(data); err != nil {
			fmt.Println(`failed gob Encode`, err)
		}
		return b.Bytes()
	}

	content := toBytes(fileContent)
	ioutil.WriteFile(FILE_DIR+filename, content, os.ModePerm)
}

func get(key [util.HashSize]byte) tkvsProtocol.Protocol {
	filename := fmt.Sprintf("%x", key)
	empKey := [util.HashSize]byte{}
	if filedata, err := ioutil.ReadFile(FILE_DIR + filename); err == nil {
		return tkvsProtocol.Protocol{tkvsProtocol.OK, empKey, filedata}
	} else {
		return tkvsProtocol.Protocol{tkvsProtocol.ERROR, empKey, []byte("Not Found")}
	}
}

func registerFiles() {
	f := func(fileName string) {
		hash := [util.HashSize]byte{}
		decoded, err := hex.DecodeString(fileName)
		if err != nil {
			log.Fatal(err)
		}
		copy(hash[:], decoded[:util.HashSize])
		fileHashMap[hash] = true
	}
	util.FilesMap(FILE_DIR, f)
}

func set(key [util.HashSize]byte, value []byte) {
	match := malScan.Scan(value)
	for _, m := range match {
		scanLog.Write(m.Rule, key)
	}
	fileHashMap[key] = true
	save(fmt.Sprintf("%x", key[:]), value)
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
			req := tkvsProtocol.Deserialize(rawReq)
			response := handleReq(req)
			send(tkvsProtocol.Serialize(response))
		case <-timeout:
			//send timeout message
			empKey := [util.HashSize]byte{}
			empData := make([]byte, 0)
			sendData := tkvsProtocol.Protocol{tkvsProtocol.CLOSE, empKey, empData}
			send(tkvsProtocol.Serialize(sendData))
			println("timeout")
			return
		case <-connClosed:
			return
		}
	}
}

func handleReq(req tkvsProtocol.Protocol) tkvsProtocol.Protocol {
	empKey := [util.HashSize]byte{}
	empData := make([]byte, 0)
	method := req.Method
	switch method {
	case tkvsProtocol.GET:
		if fileHashMap[req.Key] {
			return get(req.Key)
		} else {
			return tkvsProtocol.Protocol{tkvsProtocol.ERROR, empKey, []byte("Not Found")}
		}
	case tkvsProtocol.SET:
		if hashedData := sha256.Sum256(req.Data); hashedData == req.Key {
			if fileHashMap[req.Key] {
				return tkvsProtocol.Protocol{tkvsProtocol.FILEEXIST, empKey, empData}
			} else {
				set(req.Key, req.Data)
				fileHashMap[req.Key] = true
				return tkvsProtocol.Protocol{tkvsProtocol.OK, empKey, empData}
			}
		}
	case tkvsProtocol.CLOSE:
		return tkvsProtocol.Protocol{tkvsProtocol.CLOSE, empKey, empData}
	case tkvsProtocol.ERROR:
		return tkvsProtocol.Protocol{tkvsProtocol.ERROR, empKey, empData}
	}
	return tkvsProtocol.Protocol{tkvsProtocol.ERROR, empKey, empData}
}

func Serve(connType, laddr string) {
	registerFiles()
	malScan.ConstructRules()
	go malScan.RunRuleWatcher()

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
