package operation

import (
	"fmt"
	"github.com/pborman/uuid"
	"time"
)

type Operation struct {
	Id                   string
	Name                 string
	CreatedAt            time.Time
	Params               map[string]interface{}
	errorChainByHandlers map[string]error
}

func New(name string, params map[string]interface{}) *Operation {
	return &Operation{
		Id:                   uuid.NewUUID().String(),
		Name:                 name,
		CreatedAt:            time.Now(),
		Params:               params,
		errorChainByHandlers: map[string]error{},
	}
}

func (o *Operation) GetParam(name string) (interface{}, bool) {
	val, ok := o.Params[name]
	return val, ok
}

func (o *Operation) GetErrorChainByHandlers() map[string]error {
	return o.errorChainByHandlers
}

func (o *Operation) AddError(handlerId string, err error) {
	if o.errorChainByHandlers[handlerId] == nil {
		o.errorChainByHandlers[handlerId] = fmt.Errorf("handlerId:%s err:%w", handlerId, err)
		return
	}

	o.errorChainByHandlers[handlerId] = fmt.Errorf("%w, handlerId:%s err:%w", o.errorChainByHandlers[handlerId], handlerId, err)
}

func (o *Operation) AddErrorWithPolicy(handlerId string, policyName string, err error) {
	if o.errorChainByHandlers[handlerId] == nil {
		o.errorChainByHandlers[handlerId] = fmt.Errorf("handlerId:%s policy:%s err:%w", handlerId, policyName, err)
		return
	}

	o.errorChainByHandlers[handlerId] = fmt.Errorf("%w, handlerId:%s policy:%s err:%w", o.errorChainByHandlers[handlerId], handlerId, policyName, err)
}
