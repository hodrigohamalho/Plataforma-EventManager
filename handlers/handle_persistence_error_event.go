package handlers

import "github.com/ONSBR/Plataforma-EventManager/processor"
import log "github.com/sirupsen/logrus"

//HandlePersistenceErrorEvent handle persist error events to avoid blocking platform
func HandlePersistenceErrorEvent(c *processor.Context) error {
	log.Info("HanlePersistenceErrorEvent %s", c.Event.Name)
	return nil
}
