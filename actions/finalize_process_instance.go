package actions

import (
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
		log.Error(err)
	}
	return
}
