package asyncFallbackPolicy

import (
	"faultGuard/pkg/async/asyncOperation"
)

type Policy struct {
	fallbackFunc func(o *asyncOperation.Operation, h asyncOperation.Handler, c chan error)
}

func New(fallbackFunc func(o *asyncOperation.Operation, h asyncOperation.Handler, c chan error)) *Policy {
	return &Policy{
		fallbackFunc: fallbackFunc,
	}
}

func (p *Policy) Apply(o *asyncOperation.Operation, h asyncOperation.Handler, c chan error) {
	go p.fallbackFunc(o, h, c)
	err := <-c

	if err != nil {
		o.AddErrorWithPolicy(h.Id, "fallback", err)
	}
}
