package main

import (
	"../server/server"
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

const endpoint = "/tmp/tmp.sock"

var (
	fileDir   string
	repeats   int
	clientNum int
	dataSet   [][]byte
	keys      [][util.HashSize]byte
)

func client(start chan bool, wg *sync.WaitGroup) {
	c, err := net.Dial("unix", endpoint)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	_ = <-start
	for i := 0; i < repeats; i++ {
		for j, data := range dataSet {
			q := tkvsProtocol.Protocol{tkvsProtocol.SET, keys[j], data}
			p := tkvsProtocol.Serialize(q)
			if _, err := c.Write(p); err != nil {
				break
			}
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
	println(clientNum, repeats, fileDir)
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
	go server.Serve("unix", endpoint)
	//wait for staring server
	time.Sleep(time.Second)

	applyArgs()
	getDataSet(fileDir)

	startTime := time.Now()
	do()
	elapsed := time.Since(startTime)
	fmt.Printf("client-num: %d, repeats: %d, elapsed: %s", clientNum, repeats, elapsed.String())
}
