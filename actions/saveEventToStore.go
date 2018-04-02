package actions


package actions

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/eventstore"
	"github.com/ONSBR/Plataforma-EventManager/executor"
	log "github.com/sirupsen/logrus"
)

func SaveEventToStore(event domain.Event) error {
	if event.Owner == "" {
		event.Owner = "anonymous"
	}
	if event.AppOrigin == "" {
		event.AppOrigin = "anonymous"
	}
	log.Info("Saving event to EventStore")
	return eventstore.Push(event)
}
