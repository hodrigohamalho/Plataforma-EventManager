package actions

import (
	"github.com/ONSBR/Plataforma-EventManager/bus"
	"github.com/ONSBR/Plataforma-EventManager/domain"
)

var broker *bus.Broker

func SetBroker(_broker *bus.Broker) {
	broker = _broker
}

func PushEventToExecutor(event domain.Event) error {
	if event.Owner == "" {
		event.Owner = "anonymous"
	}
	if event.AppOrigin == "" {
		event.AppOrigin = "anonymous"
	}
	/*
		go executor.PushEvent(event)
		log.Info("Saving event to EventStore")
		return eventstore.Push(event)*/
	return broker.Publish("#", event)
}
