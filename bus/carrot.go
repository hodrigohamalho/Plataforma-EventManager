package bus

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/PMoneda/carrot"
	log "github.com/sirupsen/logrus"
)

type CarrotBroker struct {
	mux        sync.Mutex
	builder    *carrot.Builder
	subscriber *carrot.Subscriber
	publisher  *carrot.Publisher
	picker     *carrot.Picker
	config     carrot.ConnectionConfig
}

//PublishIn publish message to a exchange with a routing key
func (broker *CarrotBroker) PublishIn(exchange, routingKey string, message interface{}) error {
	if body, err := parseMessage(message); err != nil {
		return err
	} else {
		if err := broker.publisher.Publish(exchange, routingKey, carrot.Message{
			ContentType: "application/json",
			Encoding:    "utf-8",
			Data:        body,
		}); err != nil {
			return err
		}
	}
	return nil
}

//Get a message from queue
func (broker *CarrotBroker) Get(queue string, action func(*domain.Event) error) error {
	broker.mux.Lock()
	defer broker.mux.Unlock()
	context, ok, err := broker.picker.Pick(queue)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("empty_queue")
	}
	celeryMessage := new(domain.CeleryMessage)
	if err := json.Unmarshal(context.Message.Data, celeryMessage); err != nil {
		return err
	} else {
		if len(celeryMessage.Args) > 0 {
			if err := action(&celeryMessage.Args[0]); err != nil {
				return context.Nack(true)
			} else {
				return context.Ack()
			}
		} else {
			return context.Nack(true)
		}
	}
	return nil
}

//Pop a message from queue
func (broker *CarrotBroker) Pop(queue string) (event *domain.Event, err error) {
	context, ok, err := broker.picker.Pick(queue)
	if err != nil {
		return
	}
	if !ok {
		return nil, fmt.Errorf("empty_queue")
	}
	defer context.Ack()
	event = new(domain.Event)
	if err := json.Unmarshal(context.Message.Data, event); err != nil {
		return nil, err
	}
	return event, nil
}

//First a message from queue
func (broker *CarrotBroker) First(queue string) (event *domain.Event, err error) {
	context, ok, err := broker.picker.Pick(queue)
	if err != nil {
		return
	}
	if !ok {
		return nil, fmt.Errorf("empty_queue")
	}
	defer context.Nack(true)
	event = new(domain.Event)
	if err := json.Unmarshal(context.Message.Data, event); err != nil {
		return nil, err
	}
	return event, nil
}

//Swap swaps first message from some queue to another
func (broker *CarrotBroker) Swap(fromQueue, routingKey string) error {
	broker.mux.Lock()
	defer broker.mux.Unlock()
	context, ok, err := broker.picker.Pick(fromQueue)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("empty_queue")
	}
	if err := broker.Publish(routingKey, context.Message.Data); err != nil {
		context.Nack(true)
		return err
	} else {
		return context.Ack()
	}
}

//RegisterWorker to consume messages from queue
func (broker *CarrotBroker) RegisterWorker(qtd int, qname string, callback func(event *domain.Event) error) error {
	broker.subscriber.Subscribe(carrot.SubscribeWorker{
		Scale: uint(qtd),
		Queue: qname,
		Handler: func(context *carrot.MessageContext) error {
			celery := new(domain.CeleryMessage)
			if err := json.Unmarshal(context.Message.Data, celery); err != nil {
				if err := context.RedirectTo(exchangeName, "#.exception.#"); err == nil {
					context.Ack()
				} else {
					context.Nack(true)
				}
				return err
			}
			if len(celery.Args) == 0 {
				context.RedirectTo(exchangeName, "#.exception.#")
				return fmt.Errorf("invalid data")
			}
			if err := callback(&celery.Args[0]); err != nil {
				log.Error(err)
				err := context.RedirectTo(exchangeName+"_error", "#.exception.#")
				if err != nil {
					return context.Nack(true)
				}
				return context.Ack()
			} else {
				context.Ack()
			}
			return nil
		},
	})
	return nil
}

//Publish a message to specific routing key
func (broker *CarrotBroker) Publish(routingKey string, message interface{}) error {
	broker.mux.Lock()
	defer broker.mux.Unlock()
	return broker.PublishIn(exchangeName, routingKey, message)
}

var builder *carrot.Builder

var subscriber *carrot.Subscriber

var publisher *carrot.Publisher

var picker *carrot.Picker

//Init broker
func Init() {
	config := carrot.ConnectionConfig{
		Host:     os.Getenv("RABBITMQ_HOST"),
		Username: os.Getenv("RABBITMQ_USERNAME"),
		Password: os.Getenv("RABBITMQ_PASSWORD"),
		VHost:    "/",
	}
	errC := fmt.Errorf("error")
	var conn *carrot.BrokerClient
	for errC != nil {
		conn, errC = carrot.NewBrokerClient(&config)
		if errC == nil {
			break
		}
		log.Error(errC)
		time.Sleep(5 * time.Second)
	}
	builder = carrot.NewBuilder(conn)
	builder.UseVHost("plataforma_v1.0")
	builder.DeclareTopicExchange(exchangeName)
	builder.DeclareTopicExchange(exchangeName + "_error")
	builder.UpdateTopicPermission(config.Username, exchangeName)
	builder.UpdateTopicPermission(config.Username, exchangeName+"_error")
	DeclareQueue(exchangeName, EventProcessFinishedQueue, "#.finished.#")
	DeclareQueue(exchangeName, eventExecutorQueue+"_backup", "#.executor.#")
	DeclareQueue(exchangeName, eventExecutorQueue, "#.executor.#")
	DeclareQueue(exchangeName, EventPersistQueue, "#.persist.#")
	DeclareQueue(exchangeName, EventExceptionQueue, "#.exception.#")
	DeclareQueue(exchangeName, eventPersistErrorQueue, "#.persist_error.#")
	DeclareQueue(exchangeName, EventProcessFinishedQueue, "#.finished.#")
	DeclareQueue(exchangeName, EventstoreQueue, "#.store.#")
	subConn, _ := carrot.NewBrokerClient(&config)
	subscriber = carrot.NewSubscriber(subConn)
	subscriber.SetMaxRetries(30)
	pubConn, _ := carrot.NewBrokerClient(&config)
	publisher = carrot.NewPublisher(pubConn)
	pickerConn, _ := carrot.NewBrokerClient(&config)
	picker = carrot.NewPicker(pickerConn)
	carrotBroker = new(CarrotBroker)
	carrotBroker.builder = builder
	carrotBroker.config = config
	carrotBroker.picker = picker
	carrotBroker.subscriber = subscriber
	carrotBroker.publisher = publisher
}

var carrotBroker *CarrotBroker

func createQueue(exchange, queue, routingKey string) error {
	if err := builder.DeclareTopicExchange(exchange); err != nil {
		return err
	}
	if err := builder.DeclareQueue(queue); err != nil {
		return err
	}

	if err := builder.BindQueueToExchange(queue, exchange, routingKey); err != nil {
		return err
	}

	if err := builder.UpdateTopicPermission(os.Getenv("RABBITMQ_USERNAME"), exchange); err != nil {
		return err
	}
	return nil
}

func DeclareQueue(exchange, queue, routingKey string) error {
	if err := createQueue(exchange, queue, routingKey); err != nil {
		return err
	}
	if err := createQueue(exchange+"_error", queue+".error", routingKey); err != nil {
		return err
	}
	return nil
}
