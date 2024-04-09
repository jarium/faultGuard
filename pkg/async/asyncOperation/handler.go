package asyncOperation

type Handler struct {
	Id       string
	Func     func(o *Operation, c chan error)
	Policies []ErrorPolicy
}

func NewHandler(id string, f func(o *Operation, c chan error), p ...ErrorPolicy) Handler {
	return Handler{
		Id:       id,
		Func:     f,
		Policies: p,
	}
}
