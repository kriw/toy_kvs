package main

import (
	"../proto"
	"../tkvsProtocol"
	"../util"
	"crypto/sha256"
	"flag"
	"fmt"
	"net"
	"time"
)

const (
	endpoint  = "127.0.0.1:8000"
	sock      = "tcp"
	clientMax = 8
)

var (
	fileDir        string
	repeats        int
	clientNum      int
	clientParallel int
	dataSet        [][]byte
	keys           [][proto.HashSize]byte
)

func client(ch chan bool, data []byte, key [proto.HashSize]byte) {
	c, err := net.Dial(sock, endpoint)
	buf := [util.BUF_SIZE]byte{}
	if err != nil {
		panic(err)
	}
	defer c.Close()

	q := proto.RequestParam{proto.SET, uint64(len(data)), key, data}
	p := tkvsProtocol.SerializeReq(q)
	if _, err := c.Write(p); err != nil {
		println("Write Error")
		return
	}
	if _, err := c.Read(buf[:]); err != nil {
		println("Read Error")
		return
	}
	q = proto.RequestParam{proto.CLOSE_CLI, 0, [proto.HashSize]byte{}, make([]byte, 0)}
	p = tkvsProtocol.SerializeReq(q)
	_, _ = c.Write(p)
	ch <- true
}

func getDataSet(fileDir string) {
	fileList, _ := util.ReadDir(fileDir)
	for _, file := range fileList {
		if filedata, err := util.ReadFile(fileDir + "/" + file.Name()); err == nil {
			key := sha256.Sum256(filedata)
			keys = append(keys, key)
			dataSet = append(dataSet, filedata)
		}
	}
}

func applyArgs() {
	flag.IntVar(&clientNum, "client-num", 2, "an int")
	flag.IntVar(&clientParallel, "client-parallel", 2, "an int")
	flag.StringVar(&fileDir, "file", "./benchmark/files", "a string")
	flag.Parse()
}

var Padding = []byte{}

func genData(index int) ([]byte, [proto.HashSize]byte) {
	for index > 0 {
		Padding = append(Padding, byte(index&0xff))
		index >>= 8
	}
	data := append(dataSet[0], Padding...)
	return data, sha256.Sum256(data)
}

func do() {
	clientTotal := 0
	clientSending := 0
	ch := make(chan bool, 1)
	for i := 0; clientNum > i; i += 1 {
		clientTotal += 1
		for clientMax <= clientSending {
			_ = <-ch
			clientSending -= 1
		}
		data, key := genData(i)
		go client(ch, data, key)
		clientSending += 1
	}
	for clientSending > 0 {
		_ = <-ch
		clientSending -= 1
	}
}

func main() {
	applyArgs()
	getDataSet(fileDir)

	startTime := time.Now()
	do()
	elapsed := time.Since(startTime)
	fmt.Printf("client-num: %d, repeats: %d, elapsed: %s\n", clientNum, repeats, elapsed.String())

}
