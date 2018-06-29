package middlewares

import (
	"github.com/ONSBR/Plataforma-EventManager/processor"
	"github.com/ONSBR/Plataforma-EventManager/sdk"
)

//EnrichEvent get event bindings on event store
func EnrichEvent(c *processor.Context) (err error) {
	c.Event.ApplyDefaultFields()
	c.Event.Bindings, err = sdk.EventBindings(c.Event.Name)
	if err == nil && len(c.Event.Bindings) > 0 {
		c.Event.SystemID = c.Event.Bindings[0].SystemID
	}
	return
}
