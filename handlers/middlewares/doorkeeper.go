package middlewares

import (
	"fmt"
	"strings"

	"github.com/ONSBR/Plataforma-EventManager/processor"
	"github.com/ONSBR/Plataforma-EventManager/sdk"
)

//Doorkeeper look if event can proceed
func Doorkeeper(c *processor.Context) (err error) {
	if c.Event.IsReprocessing() {
		return
	}
	if strings.HasSuffix(c.Event.Name, ".done") {
		//ending events can pass
		return
	}
	canGo := sdk.EventCanProceed(c.Event)
	if !canGo {
		return fmt.Errorf(fmt.Sprintf("plataform is locked to this event %s from system %s", c.Event.Name, c.Event.SystemID))
	}
	return
}
