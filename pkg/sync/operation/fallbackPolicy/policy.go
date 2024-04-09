package fallbackPolicy

import (
	"faultGuard/pkg/sync/operation"
)

type Policy struct {
	fallbackFunc func(o *operation.Operation, h operation.Handler) error
}

func New(fallbackFunc func(o *operation.Operation, h operation.Handler) error) *Policy {
	return &Policy{
		fallbackFunc: fallbackFunc,
	}
}

func (p *Policy) Apply(o *operation.Operation, h operation.Handler) error {
	err := p.fallbackFunc(o, h)

	if err == nil {
		return nil
	}

	o.AddErrorWithPolicy(h.Id, "fallback", err)
	return err
}
