package server

import (
	"../../proto"
	"../../tkvsProtocol"
	"../../util"
	"../malScan"
	"../scanLog"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

const FILE_DIR = "./files/"

var fileHashMap = make(map[[proto.HashSize]byte]bool)

func save(filename string, fileContent []byte) {
	ioutil.WriteFile(FILE_DIR+"/"+filename, fileContent, os.ModePerm)
}

func get(key [proto.HashSize]byte) proto.ResponseParam {
	filename := fmt.Sprintf("%x", key)
	if fileData, err := ioutil.ReadFile(FILE_DIR + "/" + filename); err == nil {
		return proto.ResponseParam{proto.SUCCESS, uint64(len(fileData)), fileData}
	} else {
		notFound := []byte("Not Found")
		return proto.ResponseParam{proto.ERROR, uint64(len(notFound)), notFound}
	}
}

func registerFiles() {
	f := func(fileName string) {
		hash := [proto.HashSize]byte{}
		decoded, err := hex.DecodeString(fileName)
		if err != nil {
			log.Fatal(err)
		}
		copy(hash[:], decoded[:proto.HashSize])
		fileHashMap[hash] = true
	}
	util.FilesMap(FILE_DIR, f)
}

func set(key [proto.HashSize]byte, value []byte) {
	match := malScan.Scan(value)
	for _, m := range match {
		scanLog.Write(m.Rule, key)
	}
	fileHashMap[key] = true
	save(fmt.Sprintf("%x", key[:]), value)
}

func backgroundRead(conn net.Conn, c chan []byte, connClosed chan bool) {
	headerBuf := make([]byte, util.BUF_HEADER_SIZE)
	for {
		nr, err := conn.Read(headerBuf)
		if err != nil {
			connClosed <- true
			return
		}
		//TODO case for size > BUF_SIZE
		method, size := tkvsProtocol.GetHeader(headerBuf[:nr])
		switch proto.RequestMethod(method) {
		case proto.GET:
			c <- headerBuf[:nr]
		case proto.SET:
			wholeSize := size + proto.HEADER_REQ_SIZE
			buf := make([]byte, wholeSize)
			copy(buf, headerBuf)
			var total uint64
			for total = uint64(nr); total < wholeSize; total += uint64(nr) {
				nr, err = conn.Read(buf[total:])
				if err != nil {
					connClosed <- true
					return
				}
				if nr == 0 {
					break
				}
			}
			c <- buf[:total]
		case proto.CLOSE_CLI:
			connClosed <- true
			return
		}
	}
}

func requestHandler(conn net.Conn) {
	defer conn.Close()
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
		needClose := make(chan bool, 1)
		go func(needClose chan bool) {
			for counter := 0; counter < 100; counter += 1 {
				time.Sleep(1 * time.Second)
				select {
				case <-needClose:
					break
				default:
				}
			}
			timeout <- true
		}(needClose)

		select {
		case rawReq := <-rx:
			req := tkvsProtocol.DeserializeReq(rawReq)
			util.RequestLog(req)
			response := handleReq(req)
			sendData := tkvsProtocol.SerializeRes(response)
			util.ResponseLog(response)
			send(sendData)
			needClose <- true
		case <-timeout:
			//send timeout message
			empData := make([]byte, 0)
			sendData := proto.ResponseParam{proto.TIMEOUT, 0, empData}
			util.ResponseLog(sendData)
			send(tkvsProtocol.SerializeRes(sendData))
			println("timeout")
			return
		case <-connClosed:
			return
		}
	}
}

func handleReq(req proto.RequestParam) proto.ResponseParam {
	empData := make([]byte, 0)
	method := req.Method
	switch method {
	case proto.GET:
		if fileHashMap[req.Hash] {
			return get(req.Hash)
		} else {
			NotFound := []byte("Not Found")
			return proto.ResponseParam{proto.ERROR, uint64(len(NotFound)), NotFound}
		}
	case proto.SET:
		if hashedData := sha256.Sum256(req.Data); hashedData == req.Hash {
			if fileHashMap[req.Hash] {
				return proto.ResponseParam{proto.FILEEXIST, 0, empData}
			} else {
				set(req.Hash, req.Data)
				fileHashMap[req.Hash] = true
				return proto.ResponseParam{proto.SUCCESS, 0, empData}
			}
		}
	}
	return proto.ResponseParam{proto.ERROR, 0, empData}
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
