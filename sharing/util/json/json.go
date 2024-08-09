package json

import (
	"encoding/json"
	"sync"
)

var instance Handler
var syncOnce sync.Once

func init() {
	if instance == nil {
		syncOnce.Do(func() {
			instance = Default()
		})
	}
}

type defaultJsonHandler struct{}

func Default() Handler {
	return &defaultJsonHandler{}
}

func SetHandler(handler Handler) {
	instance = handler
}

func (djh *defaultJsonHandler) Marshal(data interface{}) (bytesData []byte, err error) {
	return json.Marshal(data)
}

func (djh *defaultJsonHandler) Unmarshal(bytesData []byte, result interface{}) (err error) {
	return json.Unmarshal(bytesData, result)
}

func Marshal(data interface{}) (bytesData []byte, err error) {
	return instance.Marshal(data)
}

func Unmarshal(bytesData []byte, result interface{}) (err error) {
	return instance.Unmarshal(bytesData, result)
}
