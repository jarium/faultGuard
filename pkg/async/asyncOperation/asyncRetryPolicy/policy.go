package asyncRetryPolicy

import (
	"github.com/jarium/faultGuard/pkg/async/asyncOperation"
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

func (p *Policy) Apply(o *asyncOperation.Operation, h asyncOperation.Handler, c chan error) error {
	var err error

	for i := 0; i < p.maxRetries; i++ {
		time.Sleep(p.delayBetweenRetries)

		go h.Func(o, c)

		err = <-c

		if err == nil {
			return nil
		}

		o.AddErrorWithPolicy(h.Id, "retry", err)
	}

	return err
}
