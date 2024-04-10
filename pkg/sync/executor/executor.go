package executor

import (
	"github.com/jarium/faultGuard/pkg/sync/operation"
	"sync"
)

type Executor struct {
	handlers map[string][]operation.Handler
	mu       sync.RWMutex
}

func New() *Executor {
	return &Executor{
		handlers: map[string][]operation.Handler{},
	}
}

// Attach a handler to an operation
func (e *Executor) Attach(operationName string, h operation.Handler) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.handlers[operationName] = append(e.handlers[operationName], h)
}

// Run handlers of an operation, call policies of it upon failure
func (e *Executor) Run(o *operation.Operation) error {
	e.mu.RLock()
	defer e.mu.RUnlock()

	handlers, ok := e.handlers[o.Name]

	if !ok {
		return nil
	}

	for _, h := range handlers {
		err := h.Func(o)

		if err == nil { //handler succeeded, proceed to execute next
			continue
		}

		o.AddError(h.Id, err)

		err = e.applyPolicies(o, h)

		if err != nil { //all policies failed
			return err
		}
	}

	return nil
}

func (e *Executor) applyPolicies(o *operation.Operation, h operation.Handler) error {
	var err error

	for _, p := range h.Policies {
		err = p.Apply(o, h)

		if err == nil { //one policy succeeded
			return nil
		}
	}

	return err
}
