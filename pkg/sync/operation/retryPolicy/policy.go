package retryPolicy

import (
	"github.com/jarium/faultGuard/pkg/sync/operation"
	"time"
)

type Policy struct {
	maxRetries          int
	delayBetweenRetries time.Duration
}

func New(maxRetries int, delayBetweenRetries time.Duration) *Policy {
	return &Policy{
		maxRetries:          maxRetries,
		delayBetweenRetries: delayBetweenRetries,
	}
}

func (p *Policy) Apply(o *operation.Operation, h operation.Handler) error {
	var err error

	for i := 0; i < p.maxRetries; i++ {
		time.Sleep(p.delayBetweenRetries)

		err = h.Func(o)

		if err == nil {
			return nil
		}

		o.AddErrorWithPolicy(h.Id, "retry", err)
	}

	return err
}
