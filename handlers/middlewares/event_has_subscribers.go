package middlewares

import (
	"github.com/ONSBR/Plataforma-EventManager/processor"
	log "github.com/sirupsen/logrus"
)

//EventHasSubscribers checks if an event has at least one operation binding
func EventHasSubscribers(c *processor.Context) error {
	log.Info("Event has subscribers %s", c.Event.Name)
	/*if refused, err := actions.RefuseEvent(c.Event); err != nil {
		return infra.NewSystemException(err.Error())
	} else if refused {
		return infra.NewSubscriberNotFoundException(fmt.Sprintf("Event %s has no subscribers", c.Event.Name))
	}*/
	return nil
}
