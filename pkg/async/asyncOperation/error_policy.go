package asyncOperation

type ErrorPolicy interface {
	Apply(o *Operation, h Handler, c chan error) error
}
