package handlers

import "github.com/ONSBR/Plataforma-EventManager/processor"
import log "github.com/sirupsen/logrus"

//HandleDoneEvent handle done event to control execution flow
func HandleDoneEvent(c *processor.Context) error {
	//TODO
	log.Info("HandleDoneEvent %s", c.Event.Name)
	return nil
}
