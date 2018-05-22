package handlers

import "github.com/ONSBR/Plataforma-EventManager/processor"
import log "github.com/sirupsen/logrus"

//HandlePersistenceEvent handle persistence event to control execution flow
func HandlePersistenceEvent(c *processor.Context) error {
	//TODO
	log.Info("HandlePersistenceEvent %s", c.Event.Name)
	return nil
}
