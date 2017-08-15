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

func get(key [util.HashSize]byte) tkvsProtocol.ResponseParam {
	filename := fmt.Sprintf("%x", key)
	if fileData, err := ioutil.ReadFile(FILE_DIR + filename); err == nil {
		return tkvsProtocol.ResponseParam{tkvsProtocol.SUCCESS, uint64(len(fileData)), fileData}
	} else {
		notFound := []byte("Not Found")
		return tkvsProtocol.ResponseParam{tkvsProtocol.ERROR, uint64(len(notFound)), notFound}
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
	buf := make([]byte, util.BUF_SIZE)
	for {
		nr, err := conn.Read(buf)
		if err != nil {
			connClosed <- true
			return
		}
		//TODO case for size > BUF_SIZE
		method, size := tkvsProtocol.GetHeader(buf[:nr])
		switch tkvsProtocol.RequestMethod(method) {
		case tkvsProtocol.GET:
			c <- buf[:nr]
		case tkvsProtocol.SET:
			wholeSize := size + tkvsProtocol.HEADER_REQ_SIZE
			for total := uint64(nr); total < wholeSize; total += uint64(nr) {
				nr, err = conn.Read(buf[total:])
				if err != nil {
					connClosed <- true
					return
				}
			}
			c <- buf[:size]
		}
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
			req := tkvsProtocol.DeserializeReq(rawReq)
			response := handleReq(req)
			send(tkvsProtocol.SerializeRes(response))
		case <-timeout:
			//send timeout message
			empData := make([]byte, 0)
			sendData := tkvsProtocol.ResponseParam{tkvsProtocol.TIMEOUT, 0, empData}
			send(tkvsProtocol.SerializeRes(sendData))
			println("timeout")
			return
		case <-connClosed:
			return
		}
	}
}

func handleReq(req tkvsProtocol.RequestParam) tkvsProtocol.ResponseParam {
	empData := make([]byte, 0)
	method := req.Method
	switch method {
	case tkvsProtocol.GET:
		if fileHashMap[req.Hash] {
			return get(req.Hash)
		} else {
			NotFound := []byte("Not Found")
			return tkvsProtocol.ResponseParam{tkvsProtocol.ERROR, uint64(len(NotFound)), NotFound}
		}
	case tkvsProtocol.SET:
		if hashedData := sha256.Sum256(req.Data); hashedData == req.Hash {
			if fileHashMap[req.Hash] {
				return tkvsProtocol.ResponseParam{tkvsProtocol.FILEEXIST, 0, empData}
			} else {
				set(req.Hash, req.Data)
				fileHashMap[req.Hash] = true
				return tkvsProtocol.ResponseParam{tkvsProtocol.SUCCESS, 0, empData}
			}
		}
		// case tkvsProtocol.CLOSE:
		// 	return tkvsProtocol.ResponseParam{tkvsProtocol.CLOSE, empKey, empData}
		// case tkvsProtocol.ERROR:
		// 	return tkvsProtocol.ResponseParam{tkvsProtocol.ERROR, empKey, empData}
	}
	return tkvsProtocol.ResponseParam{tkvsProtocol.ERROR, 0, empData}
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
