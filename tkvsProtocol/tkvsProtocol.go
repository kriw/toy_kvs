package tkvsProtocol

import (
	"../proto"
	"bytes"
	"encoding/binary"
)

func GetHeader(header []byte) (byte, uint64) {
	method := header[:proto.SIZEOF_METHOD][0]
	sizeBytes := header[proto.SIZEOF_METHOD : proto.SIZEOF_METHOD+proto.SIZEOF_INT64]

	size := binary.LittleEndian.Uint64(sizeBytes)
	return method, size
}

func encodeInt64ToBytes(n uint64) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, n)
	return buf.Bytes()
}

//FIXME
func SerializeReq(data proto.RequestParam) []byte {
	b := make([]byte, 0)
	b = append(b, byte(data.Method))
	b = append(b, encodeInt64ToBytes(data.DataSize)...)
	b = append(b, data.Hash[:]...)
	b = append(b, data.Data...)
	return b
}

func SerializeRes(data proto.ResponseParam) []byte {
	b := make([]byte, 0)
	b = append(b, byte(data.Response))
	b = append(b, encodeInt64ToBytes(data.DataSize)...)
	b = append(b, data.Data...)
	return b
}

func DeserializeReq(data []byte) proto.RequestParam {
	ret := proto.RequestParam{}
	ret.Method = proto.RequestMethod(data[0])
	ret.DataSize = binary.LittleEndian.Uint64(data[1 : 1+proto.SIZEOF_INT64])
	tmp := 1 + proto.SIZEOF_INT64
	copy(ret.Hash[:], data[tmp:tmp+proto.HashSize])
	tmp = tmp + proto.HashSize
	ret.Data = data[tmp : tmp+int(ret.DataSize)]
	return ret
}

func DeserializeRes(data []byte) proto.ResponseParam {
	ret := proto.ResponseParam{}
	ret.Response = proto.ResponseCode(data[0])
	ret.DataSize = binary.LittleEndian.Uint64(data[1 : 1+proto.SIZEOF_INT64])
	tmp := 1 + proto.SIZEOF_INT64
	ret.Data = data[tmp : tmp+int(ret.DataSize)]
	return ret
}
