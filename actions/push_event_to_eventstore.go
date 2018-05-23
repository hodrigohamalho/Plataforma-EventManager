package actions

import (
	"time"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	log "github.com/sirupsen/logrus"
)

//PushEventToEventStore from broker
func PushEventToEventStore(event *domain.Event) (err error) {
	if err = SaveEventToStore(*event); err != nil {
		log.Error(err.Error())
		time.Sleep(10 * time.Second)
	}
	return
}
