package main

import (
	"../formData"
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func reader(r io.Reader, ch chan bool) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			return
		}
		data := formData.Deserialize(buf[0:n])
		switch data.DataKind {
		case formData.CLOSE:
			ch <- false
		default:
			print(data.Data)
		}
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	c, err := net.Dial("unix", "/tmp/echo.sock")
	isConnecting := make(chan bool)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	go reader(c, isConnecting)
	for scanner.Scan() {
		select {
		case t := <-isConnecting:
			if !t {
				return
			}
		default:
			if _, err := c.Write([]byte(scanner.Text())); err != nil {
				log.Fatal("write error:", err)
				break
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading standard input:", err)
			}
		}
	}
}
