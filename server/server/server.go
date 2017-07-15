package server

import (
	"../../tkvs_protocol"
	"log"
	"net"
	"strings"
	"time"
)

var database = make(map[string]string)

func get(key string) string {
	return database[key]
}

func set(key string, value string) {
	database[key] = value
}

func backgroundRead(conn net.Conn, c chan []byte, connClosed chan bool) {
	for {
		buf := make([]byte, 512)
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
			time.Sleep(10 * time.Second)
			timeout <- true
		}()

		select {
		case rawReq := <-rx:
			req := tkvs_protocol.Deserialize(rawReq)
			response := handleReq(req)
			send(tkvs_protocol.Serialize(response))
		case <-timeout:
			//send timeout message
			sendData := tkvs_protocol.Protocol{tkvs_protocol.CLOSE, ""}
			send(tkvs_protocol.Serialize(sendData))
			println("timeout")
			return
		case <-connClosed:
			return
		}
	}
}

func handleReq(req tkvs_protocol.Protocol) tkvs_protocol.Protocol {
	method := req.Method
	data := req.Data
	switch method {
	case tkvs_protocol.GET:
		res := get(data)
		return tkvs_protocol.Protocol{tkvs_protocol.OK, res}
	case tkvs_protocol.SET:
		if ds := strings.Split(data, ","); len(ds) == 2 {
			key, value := ds[0], ds[1]
			set(key, value)
			return tkvs_protocol.Protocol{tkvs_protocol.OK, ""}
		}
	case tkvs_protocol.CLOSE:
		return tkvs_protocol.Protocol{tkvs_protocol.CLOSE, ""}
	}
	return tkvs_protocol.Protocol{tkvs_protocol.ERROR, ""}
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
