package proto

type RequestMethod byte
type ResponseCode byte

const (
	HashSize        = 32
	SIZEOF_METHOD   = 1 //[byte]
	SIZEOF_RES      = 1 //[byte]
	SIZEOF_INT64    = 8
	HEADER_REQ_SIZE = SIZEOF_METHOD + SIZEOF_INT64 + HashSize
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
	Hash     [HashSize]byte
	Data     []byte
}

func MethodToStr(method RequestMethod) string {
	switch method {
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
