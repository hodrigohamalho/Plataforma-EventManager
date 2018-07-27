package actions

import (
	"fmt"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/sdk"
)

//SetFailureProcess on ApiCore
func SetFailureProcess(event *domain.Event) (err error) {
	payload := event.Payload
	if payload == nil {
		err = fmt.Errorf("payload is nil")
		return
	}
	id, ok := payload["instance_id"]
	if !ok {
		err = fmt.Errorf("instance_id not found on payload")
	} else {
		switch i := id.(type) {
		case string:
			err = sdk.UpdateProcessInstance(i, "failure")
		default:
			err = fmt.Errorf("instance id is not string")
		}
	}

	return
}
