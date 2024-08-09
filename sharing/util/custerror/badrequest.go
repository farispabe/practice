package custerror

type BadRequest struct {
	message string
}

func NewBadRequest(message string) *BadRequest {
	return &BadRequest{message: message}
}

func (br *BadRequest) Error() string {
	return br.message
}
