package server

import (
	"log"
	"net"
)

func KeyValueServer(c net.Conn, proc func(string) string) {
	for {
		buf := make([]byte, 512)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}
		data := buf[0:nr]
		response := proc(string(data))
		_, err = c.Write([]byte(response))
		if err != nil {
			log.Fatal("Write: ", err)
		}
		_, err = c.Write([]byte("> "))
		if err != nil {
			log.Fatal("Write: ", err)
		}
	}
}
