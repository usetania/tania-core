package eventbus

import (
	"github.com/Tanibox/tania-server/src/helper/structhelper"
	"github.com/asaskevich/EventBus"
)

type TaniaEventBus interface {
	Publish(event interface{})
	Subscribe(eventName string, handlerFunc interface{})
}

type SimpleEventBus struct {
	bus EventBus.Bus
}

func NewSimpleEventBus(bus EventBus.Bus) *SimpleEventBus {
	return &SimpleEventBus{bus: bus}
}

func (e *SimpleEventBus) Publish(event interface{}) {
	name := structhelper.GetName(event)
	e.bus.Publish(name, event)
}

func (e *SimpleEventBus) Subscribe(eventName string, handler interface{}) {
	e.bus.Subscribe(eventName, handler)
}
