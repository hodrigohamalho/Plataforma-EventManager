package actions

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/eventstore"
	"github.com/ONSBR/Plataforma-EventManager/executor"
)

func PushEventToExecutor(event domain.Event) error {
	if event.Owner == "" {
		event.Owner = "anonymous"
	}
	if event.AppOrigin == "" {
		event.AppOrigin = "anonymous"
	}
	if err := executor.PushEvent(event); err != nil {
		return err
	}
	return eventstore.Push(event)
}
