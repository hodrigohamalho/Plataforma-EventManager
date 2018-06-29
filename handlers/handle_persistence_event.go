package handlers

import (
	"fmt"

	"github.com/ONSBR/Plataforma-EventManager/processor"
	log "github.com/sirupsen/logrus"
)

//HandlePersistenceEvent handle persistence event to control execution flow
func HandlePersistenceEvent(c *processor.Context) error {
	log.Debug(fmt.Sprintf("HandlePersistenceEvent %s on branch %s", c.Event.Name, c.Event.Branch))
	return c.Publish("persist", c.Event)
}
