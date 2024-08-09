package custerror

type NotFound struct {
	message string
}

func NewNotFound(message string) *NotFound {
	return &NotFound{message: message}
}

func (nf *NotFound) Error() string {
	return nf.message
}
