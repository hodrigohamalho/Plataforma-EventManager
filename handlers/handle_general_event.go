package handlers

import "github.com/ONSBR/Plataforma-EventManager/processor"
import log "github.com/sirupsen/logrus"

//HandleGeneralEvent handle general event
func HandleGeneralEvent(c *processor.Context) error {
	log.Info("HandleGeneralEvent %s", c.Event.Name)
	return nil
}
