package bus

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ONSBR/Plataforma-EventManager/domain"

	"github.com/ONSBR/Plataforma-EventManager/infra"
	rab "github.com/michaelklishin/rabbit-hole"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var rmqc *rab.Client
var connection *amqp.Connection

const vhostName = "plataforma_v1.0"
const EXCHANGE_NAME = "events.publish"
const EVENTSTORE_QUEUE = "event.store.queue"
const EVENT_PERSIST_QUEUE = "event.persist.queue"
const EVENT_EXCEPTION_QUEUE = "event.exception.queue"
const EVENT_PERSIST_ERROR_QUEUE = "event.persist.error.queue"

const EVENT_EXECUTOR_QUEUE = "event.executor.queue"

//const EVENT_EXECUTOR_QUEUE = "celery"

type Broker struct {
	mux        *sync.Mutex
	connection *amqp.Connection
	channel    *amqp.Channel
	workers    []worker
}

type worker struct {
	qtd      int
	q_name   string
	callback func([]byte) error
}

func GetBroker() *Broker {
	maxRetries := 10
	delay := 8 * time.Second
	for {
		if broker, err := Install(); err != nil {
			log.Error(err)
			log.Error(fmt.Sprintf("Trying to reconnect in %d seconds, Remaining Retries %d", delay, maxRetries))
			maxRetries--
			if maxRetries < 0 {
				panic("Cannot connect to RabbitMq, exiting")
			}
			time.Sleep(delay)
		} else {
			return broker
		}
	}

}

func Install() (*Broker, error) {

	host := infra.GetEnv("RABBITMQ_HOST", "127.0.0.1")
	amqpPort := infra.GetEnv("RABBITMQ_AMQP_PORT", "5672")
	apiPort := infra.GetEnv("RABBITMQ_API_PORT", "15672")
	user := infra.GetEnv("RABBITMQ_USERNAME", "guest")
	password := infra.GetEnv("RABBITMQ_PASSWORD", "guest")
	amqpURI := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", user, password, host, amqpPort, vhostName)
	rmqc, _ = rab.NewClient(fmt.Sprintf("http://%s:%s/", host, apiPort), user, password)

	if err := createVhost(host, apiPort, user, password); err != nil {
		return &Broker{}, err
	}

	if err := setTopicPermission(host, apiPort, user, password, vhostName, EXCHANGE_NAME); err != nil {
		return &Broker{}, err
	}

	log.Info(amqpURI)
	connection, err := amqp.Dial(amqpURI)
	if err != nil {
		return &Broker{}, fmt.Errorf("Dial: %s", err)
	}
	log.Printf("got Connection, getting Channel")
	channel, err := connection.Channel()
	if err != nil {
		return &Broker{}, fmt.Errorf("Channel: %s", err)
	}

	err = declareExchange(channel, EXCHANGE_NAME, "topic")
	if err != nil {
		return &Broker{}, err
	}
	err = declareQueues(channel, []string{EVENT_EXECUTOR_QUEUE, EVENTSTORE_QUEUE, EVENT_PERSIST_QUEUE, EVENT_EXCEPTION_QUEUE, EVENT_PERSIST_ERROR_QUEUE})
	if err != nil {
		return &Broker{}, err
	}
	//store.executor.persist.exception
	err = bindQueueToExchange(vhostName, EVENT_EXECUTOR_QUEUE, "#.executor.#")
	if err != nil {
		return &Broker{}, err
	}
	err = bindQueueToExchange(vhostName, EVENTSTORE_QUEUE, "#.store.#")
	if err != nil {
		return &Broker{}, err
	}
	err = bindQueueToExchange(vhostName, EVENT_PERSIST_QUEUE, "#.persist.#")
	if err != nil {
		return &Broker{}, err
	}
	err = bindQueueToExchange(vhostName, EVENT_EXCEPTION_QUEUE, "#.exception.#")
	if err != nil {
		return &Broker{}, err
	}

	err = bindQueueToExchange(vhostName, EVENT_PERSIST_ERROR_QUEUE, "#.persist_error.#")
	if err != nil {
		return &Broker{}, err
	}
	broker := Broker{
		mux:        new(sync.Mutex),
		connection: connection,
		channel:    channel,
		workers:    make([]worker, 0, 0),
	}
	return &broker, nil
}

func declareQueues(channel *amqp.Channel, names []string) error {
	for _, name := range names {
		if err := declareQueue(channel, name); err != nil {
			return err
		}

	}
	return nil
}

func declareQueue(channel *amqp.Channel, name string) error {
	_, err := channel.QueueDeclare(name, true, false, false, false, nil)
	_, err = channel.QueueDeclare(errorQueue(name), true, false, false, false, nil)
	return err
}

func errorQueue(q string) string {
	return q + ".error"
}

func bindQueueToExchange(vhost, queue, routingKey string) (err error) {
	_, err = rmqc.DeclareBinding(vhost, rab.BindingInfo{
		Destination:     queue,
		DestinationType: "q",
		RoutingKey:      routingKey,
		Vhost:           vhost,
		Source:          EXCHANGE_NAME,
	})

	_, err = rmqc.DeclareBinding(vhost, rab.BindingInfo{
		Destination:     errorQueue(queue),
		DestinationType: "q",
		RoutingKey:      errorQueue(queue),
		Vhost:           vhost,
		Source:          EXCHANGE_NAME + ".error",
	})
	return
}

func setTopicPermission(host, apiPort, user, password, vhostName, mainExchangeName string) (err error) {
	_, err = rmqc.UpdateTopicPermissionsIn(vhostName, user, rab.TopicPermission{
		Exchange: mainExchangeName,
		Read:     ".*",
		Write:    ".*",
	})
	return
}

func declareExchange(channel *amqp.Channel, name, exchangeType string) error {
	if err := channel.ExchangeDeclare(
		name,         // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	if err := channel.ExchangeDeclare(
		name+".error", // name
		exchangeType,  // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // noWait
		nil,           // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}
	return nil
}

func createVhost(host, port, user, password string) error {
	// creates or updates individual vhost
	if _, err := rmqc.PutVhost(vhostName, rab.VhostSettings{Tracing: false}); err != nil {
		return err
	}
	return nil
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

func (broker *Broker) PublishIn(exchange, routingKey string, message interface{}) error {

	if body, err := parseMessage(message); err != nil {
		return err
	} else {
		if err := broker.channel.Publish(
			EXCHANGE_NAME,
			routingKey,
			false,
			false,
			amqp.Publishing{
				Headers:         amqp.Table{},
				ContentType:     "application/json",
				ContentEncoding: "utf-8",
				Body:            body,
				DeliveryMode:    amqp.Persistent,
				Priority:        0,
			},
		); err != nil {
			return err
		}
	}
	return nil
}

//Publish a message to specific routing key
func (broker *Broker) Publish(routingKey string, message interface{}) error {
	broker.mux.Lock()
	defer broker.mux.Unlock()
	var err error
	err = broker.PublishIn(EXCHANGE_NAME, routingKey, message)
	recovery := false
	for err != nil {
		recovery = true
		log.Error(err)
		log.Info("Trying to reconnect to Rabbitmq")
		time.Sleep(5 * time.Second)
		if errR := broker.reconnect(); errR != nil {
			log.Error(errR)
			err = errR
		} else {
			err = broker.PublishIn(EXCHANGE_NAME, routingKey, message)
		}
	}
	if recovery {
		log.Info("Reconnected to Rabbitmq")
	}
	return nil
}

func (broker *Broker) reconnect() error {
	host := infra.GetEnv("RABBITMQ_HOST", "127.0.0.1")
	amqpPort := infra.GetEnv("RABBITMQ_AMQP_PORT", "5672")
	user := infra.GetEnv("RABBITMQ_USERNAME", "guest")
	password := infra.GetEnv("RABBITMQ_PASSWORD", "guest")
	amqpURI := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", user, password, host, amqpPort, vhostName)
	connection, err := amqp.Dial(amqpURI)
	if err != nil {
		return fmt.Errorf("Dial: %s", err)
	}
	log.Printf("got Connection, getting Channel")
	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}
	broker.channel = channel
	broker.connection = connection
	return broker.Listen()
}

func (broker *Broker) Get(queue string, action func(*domain.Event) error) error {
	broker.mux.Lock()
	defer broker.mux.Unlock()
	if msg, ok, err := broker.channel.Get(queue, false); err != nil {
		return err
	} else if !ok {
		return infra.NewEmptyQueueException("empty_queue")
	} else {
		celeryMessage := new(domain.CeleryMessage)
		if err := json.Unmarshal(msg.Body, celeryMessage); err != nil {
			return err
		} else {
			if len(celeryMessage.Args) > 0 {
				if err := action(&celeryMessage.Args[0]); err != nil {
					return msg.Nack(false, true)
				} else {
					return msg.Ack(false)
				}
			} else {
				return msg.Nack(false, true)
			}
		}

	}
	return nil
}

func (broker *Broker) Swap(fromQueue, routingKey string) error {
	broker.mux.Lock()
	defer broker.mux.Unlock()
	if msg, ok, err := broker.channel.Get(fromQueue, false); err != nil {
		return err
	} else if !ok {
		return nil
	} else if err := broker.PublishIn(EXCHANGE_NAME, routingKey, msg.Body); err != nil {
		log.Error(err)
		return msg.Nack(false, true)
	} else {
		return msg.Ack(false)
	}
	return nil
}

//RegisterWorker to consume messages from queue
func (broker *Broker) RegisterWorker(qtd int, q_name string, callback func([]byte) error) error {
	wk := worker{
		qtd:      qtd,
		q_name:   q_name,
		callback: callback,
	}
	broker.workers = append(broker.workers, wk)
	return nil
}

func (broker *Broker) Listen() error {
	return broker.runWorkers()
}

//RegisterWorker to consume messages from queue
func (broker *Broker) runWorkers() error {
	for _, wk := range broker.workers {
		for i := 0; i < wk.qtd; i++ {
			ch, err := broker.connection.Channel()
			if err != nil {
				return err
			}
			msgs, err := ch.Consume(
				wk.q_name, // queue
				"",        // consumer
				false,     // auto-ack
				false,     // exclusive
				false,     // no-local
				false,     // no-wait
				nil,       // args
			)
			go func() {
				for event := range msgs {
					if err := wk.callback(event.Body); err != nil {
						log.Error(err)
						if err := broker.PublishIn(EXCHANGE_NAME+".error", errorQueue(wk.q_name), event.Body); err != nil {
							//TODO what is the best approach?
							log.Error(err)
							event.Nack(false, true)
							return
						} else {
							event.Nack(false, false)
						}
						log.Error(err)
					} else {
						event.Ack(false)
					}
				}
			}()
		}
	}

	return nil
}
