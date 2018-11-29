package bus

import (
	"encoding/json"

	rab "github.com/michaelklishin/rabbit-hole"
	"github.com/streadway/amqp"
)

var rmqc *rab.Client
var connection *amqp.Connection

const vhostName = "plataforma_v1.0"
const exchangeName = "events.publish"

//EventstoreQueue is a queue that will receive all events to be stored in influxdb
const EventstoreQueue = "event.store.queue"

//EventsReplayQueue is the queue that will receive events if solution is recording events
const EventsReplayQueue = "event.replay.%s.queue"

//EventPersistQueue is a queue that will receive all persist event
const EventPersistQueue = "event.persist.queue"

//EventExceptionQueue is a queue that will receive all exception event
const EventExceptionQueue = "event.exception.queue"

const eventPersistErrorQueue = "event.persist.error.queue"

//EventPersistRequestQueue is a queue that will receive all persist event
const EventPersistRequestQueue = "event.persist.request.queue"

//EventProcessFinishedQueue is a queue that will receive all persist event
const EventProcessFinishedQueue = "event.process.finished.queue"

const eventExecutorQueue = "event.executor.queue"

//GetBroker returns a new broker
func GetBroker() Dispatcher {
	return carrotBroker
}

func parseMessage(message interface{}) (body []byte, err error) {
	switch t := message.(type) {
	case []byte:
		body = t
	default:
		body, err = json.Marshal(t)
		if err != nil {
			return nil, err
		}
	}
	return
}
