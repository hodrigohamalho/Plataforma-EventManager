package actions

import (
	"fmt"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/sdk"
	log "github.com/sirupsen/logrus"
)

//FinalizeProcessInstance on ApiCore
func FinalizeProcessInstance(event *domain.Event) (err error) {
	if event.Payload == nil {
		err = fmt.Errorf("event payload is nil")
		return
	}
	id, ok := event.Payload["instance_id"]
	if !ok {
		err = fmt.Errorf("payload has no instance_id")
		return
	}
	switch i := id.(type) {
	case string:
		log.Debug(fmt.Sprintf("update process instance %s to finished from event %s", id, event.Name))
		err = sdk.UpdateProcessInstance(i, "finished")
		if err != nil {
			log.Error(err)
		}
	default:
		log.Error("instance_id is not valid")
	}
	return
}
