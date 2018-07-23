package middlewares

import (
	"github.com/ONSBR/Plataforma-EventManager/processor"
	"github.com/ONSBR/Plataforma-EventManager/sdk"
	"github.com/ONSBR/Plataforma-EventManager/util"
)

//EnrichEvent get event bindings on event store
func EnrichEvent(c *processor.Context) (err error) {
	c.Event.ApplyDefaultFields()
	c.Event.Bindings, err = sdk.EventBindings(c.Event.Name)
	if err == nil && len(c.Event.Bindings) > 0 {
		c.Event.SystemID = c.Event.Bindings[0].SystemID
		c.Event.Version = c.Event.Bindings[0].Version
	}
	if c.Event.IdempotencyKey != "" {
		return
	}
	attrs := map[string]interface{}{
		"name":    c.Event.Name,
		"version": c.Event.Version,
		"system":  c.Event.SystemID,
		"payload": c.Event.Payload,
	}
	key, err := util.SHA1(attrs)
	if err != nil {
		return err
	}
	c.Event.IdempotencyKey = key
	return
}
