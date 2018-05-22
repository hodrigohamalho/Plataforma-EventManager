package handlers

import "github.com/ONSBR/Plataforma-EventManager/processor"
import log "github.com/sirupsen/logrus"

//HandleExceptionEvent handle exception events
func HandleExceptionEvent(c *processor.Context) error {
	//TODO
	log.Info("HandleExceptionEvent %s", c.Event.Name)
	return nil
}
