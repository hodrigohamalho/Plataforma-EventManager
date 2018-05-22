package handlers

import "github.com/ONSBR/Plataforma-EventManager/processor"
import log "github.com/sirupsen/logrus"

//HandleSystemEvent handle all system event received on event manager
func HandleSystemEvent(c *processor.Context) error {
	//TODO
	log.Info("HandleSystemEvent %s", c.Event.Name)
	return nil
}
