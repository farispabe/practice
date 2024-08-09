package validator

type Handler interface {
	Validate(request Request) error
}
