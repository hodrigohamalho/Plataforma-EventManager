package middlewares

import "github.com/ONSBR/Plataforma-EventManager/processor"
import log "github.com/sirupsen/logrus"

//EnrichEvent get event bindings on event store
func EnrichEvent(c *processor.Context) error {
	log.Info("Enriching event %s", c.Event.Name)
	return nil
}
