package asyncOperation

import (
	"context"
	"fmt"
	"github.com/pborman/uuid"
	"sync"
	"time"
)

type Operation struct {
	Id                   string
	Name                 string
	CreatedAt            time.Time
	Params               map[string]interface{}
	errorChainByHandlers map[string]error
	mu                   sync.RWMutex
	ctx                  context.Context
	cancelFunc           context.CancelFunc
}

func New(name string, params map[string]interface{}) *Operation {
	ctx, cancel := context.WithCancel(context.Background())

	return &Operation{
		Id:                   uuid.NewUUID().String(),
		Name:                 name,
		CreatedAt:            time.Now(),
		Params:               params,
		errorChainByHandlers: map[string]error{},
		ctx:                  ctx,
		cancelFunc:           cancel,
	}
}

func (o *Operation) Done() {
	o.cancelFunc()
}

func (o *Operation) GetParam(name string) (interface{}, bool) {
	val, ok := o.Params[name]
	return val, ok
}

// GetErrorChainByHandlers after all handlers finished executing. Waits until all the handlers are finished.
func (o *Operation) GetErrorChainByHandlers() map[string]error {
	select {
	case <-o.ctx.Done():
		return o.errorChainByHandlers
	}
}

func (o *Operation) AddError(handlerId string, err error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.errorChainByHandlers[handlerId] == nil {
		o.errorChainByHandlers[handlerId] = fmt.Errorf("handlerId:%s err:%w", handlerId, err)
		return
	}

	o.errorChainByHandlers[handlerId] = fmt.Errorf("%w, handlerId:%s err:%w", o.errorChainByHandlers[handlerId], handlerId, err)
}

func (o *Operation) AddErrorWithPolicy(handlerId string, policyName string, err error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.errorChainByHandlers[handlerId] == nil {
		o.errorChainByHandlers[handlerId] = fmt.Errorf("handlerId:%s policy:%s err:%w", handlerId, policyName, err)
		return
	}

	o.errorChainByHandlers[handlerId] = fmt.Errorf("%w, handlerId:%s policy:%s err:%w", o.errorChainByHandlers[handlerId], handlerId, policyName, err)
}
