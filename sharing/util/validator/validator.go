package validator

import "sync"

var instance Handler
var syncOnce sync.Once

func init() {
	if instance == nil {
		syncOnce.Do(func() {
			instance = Default()
		})
	}
}

type Request interface {
	Validate() error
}

type defaultValidator struct{}

func Default() Handler {
	return &defaultValidator{}
}

func SetHandler(handler Handler) {
	instance = handler
}

func (dv *defaultValidator) Validate(request Request) error {
	return request.Validate()
}

func Validate(request Request) error {
	return instance.Validate(request)
}
