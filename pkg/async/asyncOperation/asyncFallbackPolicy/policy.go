package asyncFallbackPolicy

import (
	"github.com/jarium/faultGuard/pkg/async/asyncOperation"
)

type Policy struct {
	fallbackFunc func(o *asyncOperation.Operation, h asyncOperation.Handler, c chan error)
}

func New(fallbackFunc func(o *asyncOperation.Operation, h asyncOperation.Handler, c chan error)) *Policy {
	return &Policy{
		fallbackFunc: fallbackFunc,
	}
}

func (p *Policy) Apply(o *asyncOperation.Operation, h asyncOperation.Handler, c chan error) error {
	go p.fallbackFunc(o, h, c)
	err := <-c

	if err == nil {
		return nil
	}

	o.AddErrorWithPolicy(h.Id, "fallback", err)
	return err
}
