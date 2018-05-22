package middlewares

import (
	"github.com/ONSBR/Plataforma-EventManager/processor"
	"github.com/ONSBR/Plataforma-EventManager/sdk"
	log "github.com/sirupsen/logrus"
)

//EnrichEvent get event bindings on event store
func EnrichEvent(c *processor.Context) (err error) {
	log.Debug("Enriching event ", c.Event.Name)
	c.Event.ApplyDefaultFields()
	c.Event.Bindings, err = sdk.EventBindings(c.Event.Name)
	return
}
