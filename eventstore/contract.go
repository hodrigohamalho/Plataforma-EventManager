package eventstore

import "github.com/ONSBR/Plataforma-EventManager/domain"

//EventStore defines a contract to store all events on platform
type EventStore interface {
	Push(domain.Event) error
}

func Push(event domain.Event) error {
	return NewInfluxStorage().Push(event)
}
