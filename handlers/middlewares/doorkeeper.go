package middlewares

import (
	"fmt"

	"github.com/ONSBR/Plataforma-EventManager/processor"
	"github.com/ONSBR/Plataforma-EventManager/sdk"
	log "github.com/sirupsen/logrus"
)

//Doorkeeper look if event can proceed
func Doorkeeper(c *processor.Context) (err error) {
	log.Debug("check if event can pass throug platform gateway")
	if c.Event.IsReprocessing() {
		return
	}

	canGo := sdk.EventCanProceed(c.Event)
	if !canGo {
		return fmt.Errorf(fmt.Sprintf("plataform is locked to this event %s from system %s", c.Event.Name, c.Event.SystemID))
	}
	return
}
