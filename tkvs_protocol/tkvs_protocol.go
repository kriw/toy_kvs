package tkvs_protocol

type RequestMethod byte

const (
	GET RequestMethod = iota
	SET
	OK
	CLOSE
	ERROR
)

type Protocol struct {
	Method RequestMethod
	Data   string
}

func Serialize(data Protocol) []byte {
	return append([]byte(data.Data), byte(data.Method))
}

func Deserialize(data []byte) Protocol {
	return Protocol{RequestMethod(data[len(data)-1]), string(data[0 : len(data)-1])}
}
