package tkvsProtocol

import (
	"../util"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
)

type RequestMethod byte
type ResponseCode byte

const (
	SIZEOF_REQ   = 1 //[byte]
	SIZEOF_RES   = 1 //[byte]
	SIZEOF_INT64 = 8
)

//For Request
const (
	GET RequestMethod = iota
	SET
	ERROR_INPUT
)

//For Response
const (
	SUCCESS ResponseCode = iota
	NOTFOUND
	FILEEXIST
	TIMEOUT
	ERROR
)

type ResponseParam struct {
	Response ResponseCode
	DataSize uint64
	Data     []byte
}

type RequestParam struct {
	Method RequestMethod
	Size   uint64
	Hash   [util.HashSize]byte
	Data   []byte
}

func GetHeader(header []byte) (byte, uint64) {
	method := header[:SIZEOF_REQ][0]
	sizeBytes := header[SIZEOF_REQ : SIZEOF_REQ+SIZEOF_INT64]

	size := binary.BigEndian.Uint64(sizeBytes)
	return method, size
}

func SerializeReq(data RequestParam) []byte {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	if err := e.Encode(data); err != nil {
		fmt.Println(`failed gob Encode`, err)
	}
	return b.Bytes()
}

//FIXME
func SerializeRes(data ResponseParam) []byte {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	if err := e.Encode(data); err != nil {
		fmt.Println(`failed gob Encode`, err)
	}
	return b.Bytes()
}

func DeserializeReq(data []byte) RequestParam {
	m := RequestParam{}
	b := bytes.Buffer{}
	b.Write(data)
	dec := gob.NewDecoder(&b)
	if err := dec.Decode(&m); err != nil {
		fmt.Println(`failed gob Decode`, err)
		return RequestParam{GET, 0, [util.HashSize]byte{}, make([]byte, 0)}
	} else {
		return m
	}
}

func DeserializeRes(data []byte) ResponseParam {
	m := ResponseParam{}
	b := bytes.Buffer{}
	b.Write(data)
	dec := gob.NewDecoder(&b)
	if err := dec.Decode(&m); err != nil {
		fmt.Println(`failed gob Decode`, err)
		return ResponseParam{ERROR, 0, make([]byte, 0)}
	} else {
		return m
	}
}
