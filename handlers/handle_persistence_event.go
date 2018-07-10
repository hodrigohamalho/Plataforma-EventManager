package handlers

import (
	"fmt"
	"sync"

	"github.com/ONSBR/Plataforma-EventManager/bus"
	"github.com/ONSBR/Plataforma-EventManager/processor"
	"github.com/ONSBR/Plataforma-EventManager/sdk"
	log "github.com/sirupsen/logrus"
)

var _hub map[string]bool
var once sync.Once
var mut sync.Mutex

//HandlePersistenceEvent handle persistence event to control execution flow
func HandlePersistenceEvent(c *processor.Context) error {
	log.Debug(fmt.Sprintf("HandlePersistenceEvent %s on branch %s", c.Event.Name, c.Event.Branch))
	once.Do(func() {
		_hub = make(map[string]bool)
	})
	_, ok := _hub[c.Event.SystemID]
	if !ok {
		mut.Lock()
		defer mut.Unlock()
		_hub[c.Event.SystemID] = true
		bus.DeclareQueue("persist", fmt.Sprintf("persist.%s.queue", c.Event.SystemID), c.Event.SystemID)
		ok := sdk.InitPersistHandler(c.Event)
		if !ok {
			log.Error("cannot start handling on maestro, persist events will not be processed by system: ", c.Event.SystemID)
		} else {
			//log.Info("wakeup maestro for system: ", c.Event.SystemID)
		}
	}
	return c.PublishIn("persist", c.Event.SystemID, c.Event)
}
