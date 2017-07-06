package formData

//HACKME better name of this package

type RequestKind byte

const (
	GET RequestKind = iota
	SET
	OK
	CLOSE
)

type FormData struct {
	DataKind RequestKind
	Data     string
}

func Serialize(formData FormData) []byte {
	return append([]byte(formData.Data), byte(formData.DataKind))
}

func Deserialize(formData []byte) FormData {
	return FormData{RequestKind(formData[len(formData)-1]), string(formData[0 : len(formData)-1])}
}
