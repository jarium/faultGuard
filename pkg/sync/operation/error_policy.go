package operation

type ErrorPolicy interface {
	Apply(o *Operation, h Handler) error
}
