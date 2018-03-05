package eventstore

import "github.com/ONSBR/Plataforma-EventManager/domain"

type InfluxDb struct {
}

//Push event to Influx
func (s InfluxDb) Push(event domain.Event) error {
	return nil
}

//NewInfluxStorage creates a new instance of influx storage
func NewInfluxStorage() InfluxDb {
	influx := new(InfluxDb)
	return *influx
}
