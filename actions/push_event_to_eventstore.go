package actions

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	log "github.com/sirupsen/logrus"
)

//PushEventToEventStore from broker
func PushEventToEventStore(payload []byte) (err error) {
	celeryMessage := new(domain.CeleryMessage)
	err = json.Unmarshal(payload, celeryMessage)
	if len(celeryMessage.Args) == 0 {
		return fmt.Errorf("celery message should have an event")
	}
	event := celeryMessage.Args[0]
	if err != nil {
		log.Error(err.Error())
		return
	}
	if err = SaveEventToStore(event); err != nil {
		log.Error(err.Error())
		time.Sleep(10 * time.Second)
	}
	return
}
