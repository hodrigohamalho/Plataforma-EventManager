package processor

import (
	"github.com/ONSBR/Plataforma-EventManager/bus"
	"github.com/ONSBR/Plataforma-EventManager/domain"
)

//Context is structure that have the current state
type Context struct {
	Event      *domain.Event
	dispatcher bus.Dispatcher
	Session    map[string]string
}

//Publish sends message to broker
func (c *Context) Publish(routingKey string, event *domain.Event) error {
	return c.dispatcher.Publish("store."+routingKey, event.ToCeleryMessage())
}