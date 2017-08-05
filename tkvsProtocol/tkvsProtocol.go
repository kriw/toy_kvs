package tkvsProtocol

import (
	"../util"
	"bytes"
	"encoding/gob"
	"fmt"
)

type RequestMethod byte

const (
	GET RequestMethod = iota
	SET
	OK
	SAVE
	CLOSE
	ERROR
)

type Protocol struct {
	Method RequestMethod
	Key    [util.HashSize]byte
	Data   []byte
}

func Serialize(data Protocol) []byte {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	if err := e.Encode(data); err != nil {
		fmt.Println(`failed gob Encode`, err)
	}
	return b.Bytes()
}

func Deserialize(data []byte) Protocol {
	m := Protocol{}
	b := bytes.Buffer{}
	b.Write(data)
	dec := gob.NewDecoder(&b)
	if err := dec.Decode(&m); err != nil {
		fmt.Println(`failed gob Decode`, err)
		return Protocol{ERROR, [util.HashSize]byte{}, make([]byte, 0)}
	} else {
		return m
	}
}
