package tkvs_protocol

//HACKME better name of this package

type RequestKind byte

const (
	GET RequestKind = iota
	SET
	OK
	CLOSE
)

type Protocol struct {
	DataKind RequestKind
	Data     string
}

func Serialize(data Protocol) []byte {
	return append([]byte(data.Data), byte(data.DataKind))
}

func Deserialize(data []byte) Protocol {
	return Protocol{RequestKind(data[len(data)-1]), string(data[0 : len(data)-1])}
}
