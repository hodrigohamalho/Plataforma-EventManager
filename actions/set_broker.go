package actions

import (
	"github.com/ONSBR/Plataforma-EventManager/bus"
)

var broker *bus.Broker

func SetBroker(_broker *bus.Broker) {
	broker = _broker
}
