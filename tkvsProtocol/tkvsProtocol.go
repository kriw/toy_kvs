package tkvsProtocol

import (
	"../util"
	"bytes"
	"encoding/binary"
)

type RequestMethod byte
type ResponseCode byte

const (
	SIZEOF_METHOD   = 1 //[byte]
	SIZEOF_RES      = 1 //[byte]
	SIZEOF_INT64    = 8
	HEADER_REQ_SIZE = SIZEOF_METHOD + SIZEOF_INT64 + util.HashSize
)

//For Request
const (
	GET RequestMethod = iota
	SET
	CLOSE_CLI
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
	Method   RequestMethod
	DataSize uint64
	Hash     [util.HashSize]byte
	Data     []byte
}

func GetHeader(header []byte) (byte, uint64) {
	method := header[:SIZEOF_METHOD][0]
	sizeBytes := header[SIZEOF_METHOD : SIZEOF_METHOD+SIZEOF_INT64]

	size := binary.LittleEndian.Uint64(sizeBytes)
	return method, size
}

func RequestToStr(req RequestParam) string {
	switch res {
	case GET:
		return "Get"
	case SET:
		return "Set"
	case CLOSE_CLI:
		return "Close Client"
	case ERROR_INPUT:
		return "Error Input"
	}
	return "Unknow"
}
func ResponseToStr(res ResponseCode) string {
	switch res {
	case SUCCESS:
		return "Success"
	case NOTFOUND:
		return "Not Found"
	case FILEEXIST:
		return "File Exist"
	case TIMEOUT:
		return "Time out"
	case ERROR:
		return "Error"
	}
	return "Unknow"
}

func encodeInt64ToBytes(n uint64) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, n)
	return buf.Bytes()
}

//FIXME
func SerializeReq(data RequestParam) []byte {
	b := make([]byte, 0)
	b = append(b, byte(data.Method))
	b = append(b, encodeInt64ToBytes(data.DataSize)...)
	b = append(b, data.Hash[:]...)
	b = append(b, data.Data...)
	return b
}

func SerializeRes(data ResponseParam) []byte {
	b := make([]byte, 0)
	b = append(b, byte(data.Response))
	b = append(b, encodeInt64ToBytes(data.DataSize)...)
	b = append(b, data.Data...)
	return b
}

func DeserializeReq(data []byte) RequestParam {
	ret := RequestParam{}
	ret.Method = RequestMethod(data[0])
	ret.DataSize = binary.LittleEndian.Uint64(data[1 : 1+SIZEOF_INT64])
	tmp := 1 + SIZEOF_INT64
	copy(ret.Hash[:], data[tmp:tmp+util.HashSize])
	tmp = tmp + util.HashSize
	ret.Data = data[tmp : tmp+int(ret.DataSize)]
	return ret
}

func DeserializeRes(data []byte) ResponseParam {
	ret := ResponseParam{}
	ret.Response = ResponseCode(data[0])
	ret.DataSize = binary.LittleEndian.Uint64(data[1 : 1+SIZEOF_INT64])
	tmp := 1 + SIZEOF_INT64
	ret.Data = data[tmp : tmp+int(ret.DataSize)]
	return ret
}
