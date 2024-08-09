package json

type Handler interface {
	Marshal(data interface{}) (bytesData []byte, err error)
	Unmarshal(bytesData []byte, result interface{}) (err error)
}
