package actions

import (
	"encoding/json"
	"fmt"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/sdk"
	log "github.com/sirupsen/logrus"
)

//FinalizeProcessInstance on ApiCore
func FinalizeProcessInstance(payload []byte) (err error) {
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
	id := event.Payload["instance_id"].(string)
	err = sdk.UpdateProcessInstance(id, "finished")
	return
}
