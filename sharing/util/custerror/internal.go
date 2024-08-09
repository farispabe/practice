package custerror

type Internal struct {
	message string
}

func NewInternal(message string) *Internal {
	return &Internal{message: message}
}

func (i *Internal) Error() string {
	return i.message
}
