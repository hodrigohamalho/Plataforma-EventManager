package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ONSBR/Plataforma-EventManager/infra/factories"

	"github.com/ONSBR/Plataforma-EventManager/actions"
	"github.com/ONSBR/Plataforma-EventManager/api"
	"github.com/ONSBR/Plataforma-EventManager/bus"
	"github.com/ONSBR/Plataforma-EventManager/eventstore"
	log "github.com/sirupsen/logrus"
)

var local bool

func init() {
	flag.BoolVar(&local, "local", false, "to run service with local rabbitmq and services")
	os.Setenv("DATABASE", "event_manager")
	os.Setenv("RETENTION_POLICY", "platform_events")
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

}

func registerActionsToRabbitMq() bus.Dispatcher {
	broker := factories.GetBroker()
	broker.RegisterWorker(1, bus.EventstoreQueue, actions.PushEventToEventStore)
	broker.RegisterWorker(1, bus.EventProcessFinishedQueue, actions.FinalizeProcessInstance)
	broker.RegisterWorker(1, bus.EventExceptionQueue, actions.SetFailureProcess)
	return broker
}

func main() {
	flag.Parse()
	if local {
		os.Setenv("RABBITMQ_HOST", "localhost")
		os.Setenv("RABBITMQ_USERNAME", "guest")
		os.Setenv("RABBITMQ_PASSWORD", "guest")
		os.Setenv("PORT", "8089")
	}
	bus.Init()
	fmt.Println(logo())
	log.Info("Starting Event Manager")
	log.Info("Installing Bus")
	registerActionsToRabbitMq()
	log.Info("Starting Mongo")
	eventstore.Install()
	log.Info("Starting API")
	api.Build()
}

func logo() (logo string) {
	logo = `
 	 ______               _     __  __
	|  ____|             | |   |  \/  |
	| |____   _____ _ __ | |_  | \  / | __ _ _ __   __ _  __ _  ___ _ __
	|  __\ \ / / _ \ '_ \| __| | |\/| |/ _' | '_ \ / _' |/ _' |/ _ \ '__|
	| |___\ V /  __/ | | | |_  | |  | | (_| | | | | (_| | (_| |  __/ |
	|______\_/ \___|_| |_|\__| |_|  |_|\__,_|_| |_|\__,_|\__, |\___|_|
	                                                      __/ |
	                                                     |___/
	`
	return
}
