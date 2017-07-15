package server

import (
	"../../tkvs_protocol"
	"log"
	"net"
	"time"
)

var database = make(map[string][]byte)

func get(key string) []byte {
	return database[key]
}

func set(key string, value []byte) {
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
			b := make([]byte, 0)
			sendData := tkvs_protocol.Protocol{tkvs_protocol.CLOSE, b, b}
			send(tkvs_protocol.Serialize(sendData))
			println("timeout")
			return
		case <-connClosed:
			return
		}
	}
}

func handleReq(req tkvs_protocol.Protocol) tkvs_protocol.Protocol {
	empty := make([]byte, 0)
	method := req.Method
	key := string(req.Key)
	data := req.Data
	switch method {
	case tkvs_protocol.GET:
		res := get(string(key))
		return tkvs_protocol.Protocol{tkvs_protocol.OK, empty, res}
	case tkvs_protocol.SET:
		value := data
		set(key, value)
		return tkvs_protocol.Protocol{tkvs_protocol.OK, empty, empty}
	case tkvs_protocol.CLOSE:
		return tkvs_protocol.Protocol{tkvs_protocol.CLOSE, empty, empty}
	}
	return tkvs_protocol.Protocol{tkvs_protocol.ERROR, empty, empty}
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
