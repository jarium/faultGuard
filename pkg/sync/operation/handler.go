package operation

type Handler struct {
	Id       string
	Func     func(o *Operation) error
	Policies []ErrorPolicy
}

func NewHandler(id string, f func(o *Operation) error, p ...ErrorPolicy) Handler {
	return Handler{
		Id:       id,
		Func:     f,
		Policies: p,
	}
}
