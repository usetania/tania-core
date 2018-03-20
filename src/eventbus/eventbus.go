package eventbus

import (
	"github.com/asaskevich/EventBus"
)

type TaniaEventBus interface {
	Publish(eventName string, event interface{})
	Subscribe(eventName string, handlerFunc interface{})
}

type SimpleEventBus struct {
	bus EventBus.Bus
}

func NewSimpleEventBus(bus EventBus.Bus) *SimpleEventBus {
	return &SimpleEventBus{bus: bus}
}

func (e *SimpleEventBus) Publish(eventName string, event interface{}) {
	e.bus.Publish(eventName, event)
}

func (e *SimpleEventBus) Subscribe(eventName string, handler interface{}) {
	e.bus.Subscribe(eventName, handler)
}
