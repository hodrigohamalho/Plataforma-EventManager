package eventstore

import "github.com/ONSBR/Plataforma-EventManager/domain"

func Push(event domain.Event) error {
	return influxPush(event)
}
