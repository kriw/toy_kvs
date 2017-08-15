package main

import (
	"../tkvsProtocol"
	"../util"
	"crypto/sha256"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"sync"
	"time"
)

const (
	endpoint = "127.0.0.1:8000"
	sock     = "tcp"
	BUF_SIZE = 256 * 1024 * 1024
)

var (
	fileDir   string
	repeats   int
	clientNum int
	dataSet   [][]byte
	keys      [][util.HashSize]byte
)

func client(start chan bool, wg *sync.WaitGroup) {
	c, err := net.Dial(sock, endpoint)
	buf := [BUF_SIZE]byte{}
	if err != nil {
		panic(err)
	}
	defer c.Close()

	_ = <-start
	for i := 0; i < repeats; i++ {
		for j, data := range dataSet {
			q := tkvsProtocol.RequestParam{tkvsProtocol.SET, uint64(len(data)), keys[j], data}
			p := tkvsProtocol.SerializeReq(q)
			if _, err := c.Write(p); err != nil {
				break
			}
			_, _ = c.Read(buf[:])
		}
	}
	wg.Done()
}

func getDataSet(fileDir string) {
	fileList, _ := ioutil.ReadDir(fileDir)
	for _, file := range fileList {
		if filedata, err := ioutil.ReadFile(file.Name()); err == nil {
			key := sha256.Sum256(filedata)
			keys = append(keys, key)
			dataSet = append(dataSet, filedata)
		}
	}
}

func applyArgs() {
	flag.IntVar(&clientNum, "client-num", 2, "an int")
	flag.IntVar(&repeats, "repeats", 5, "an int")
	flag.StringVar(&fileDir, "file", "./benchmark/files", "a string")
	flag.Parse()
}

func do() {
	var wg sync.WaitGroup
	start := make(chan bool)
	for i := 0; i < clientNum; i++ {
		wg.Add(1)
		go client(start, &wg)
	}
	for i := 0; i < clientNum; i++ {
		start <- true
	}
	wg.Wait()
}

func main() {
	applyArgs()
	getDataSet(fileDir)

	startTime := time.Now()
	do()
	elapsed := time.Since(startTime)
	fmt.Printf("client-num: %d, repeats: %d, elapsed: %s", clientNum, repeats, elapsed.String())
}
