package asyncExecutor

import (
	"faultGuard/pkg/async/asyncOperation"
	"sync"
)

type Executor struct {
	handlers map[string][]asyncOperation.Handler
	mu       sync.RWMutex
}

func New() *Executor {
	return &Executor{
		handlers: map[string][]asyncOperation.Handler{},
	}
}

// Attach a handler to an operation
func (e *Executor) Attach(operationName string, h asyncOperation.Handler) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.handlers[operationName] = append(e.handlers[operationName], h)
}

// Run handlers of an operation, call policies of it upon failure
func (e *Executor) Run(o *asyncOperation.Operation) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	handlers, ok := e.handlers[o.Name]

	if !ok {
		return
	}

	wg := &sync.WaitGroup{}

	for _, h := range handlers {
		wg.Add(1)

		go func(h asyncOperation.Handler, wg *sync.WaitGroup) {
			defer wg.Done()

			c := make(chan error)

			go h.Func(o, c)

			e.applyPoliciesOnError(o, h, c)
		}(h, wg)
	}

	wg.Wait()
	o.Done()
}

func (e *Executor) applyPoliciesOnError(o *asyncOperation.Operation, h asyncOperation.Handler, c chan error) {
	o.AddError(h.Id, <-c)

	for _, p := range h.Policies {
		p.Apply(o, h, c)
	}
}