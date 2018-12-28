package actions

import (
	"encoding/json"
	"fmt"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/sdk"
	log "github.com/sirupsen/logrus"
)

//FinalizeProcessInstance on ApiCore
func FinalizeProcessInstance(event *domain.Event) (err error) {
	id := event.InstanceID
	log.Debug(fmt.Sprintf("update process instance %s to finished from event %s", id, event.Name))
	err = sdk.UpdateProcessInstance(id, "finished")
	if err != nil {
		if buf, err := json.Marshal(event); err != nil {
			log.Error(err)
		} else {
			log.Error(fmt.Sprintf("raw event: %s", string(buf)))
		}
		log.Error(err)
	}
	return
}
