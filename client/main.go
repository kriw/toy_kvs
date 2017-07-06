package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func reader(r io.Reader) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			return
		}
		fmt.Print(string(buf[0:n]))
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	c, err := net.Dial("unix", "/tmp/echo.sock")
	if err != nil {
		panic(err)
	}
	defer c.Close()

	go reader(c)
	for scanner.Scan() {
		if _, err := c.Write([]byte(scanner.Text())); err != nil {
			log.Fatal("write error:", err)
			break
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
	}
}
