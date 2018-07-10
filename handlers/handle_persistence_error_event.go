package handlers

import (
	"fmt"

	"github.com/ONSBR/Plataforma-EventManager/processor"
	log "github.com/sirupsen/logrus"
)

//HandlePersistenceErrorEvent handle persist error events to avoid blocking platform
func HandlePersistenceErrorEvent(c *processor.Context) error {
	log.Debug(fmt.Sprintf("HanlePersistenceErrorEvent %s on branch %s", c.Event.Name, c.Event.Branch))
	return c.Publish("store.executor", c.Event)
}
