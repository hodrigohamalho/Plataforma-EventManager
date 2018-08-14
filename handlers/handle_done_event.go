package handlers

import (
	"fmt"

	"github.com/ONSBR/Plataforma-EventManager/actions"
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/processor"
	log "github.com/sirupsen/logrus"
)

//HandleDoneEvent handle done event to control execution flow
func HandleDoneEvent(c *processor.Context) error {
	log.Debug(fmt.Sprintf("HandleDoneEvent %s on branch %s", c.Event.Name, c.Event.Branch))
	if c.Event.IsExecution() {
		return handleExecutionDone(c)
	}
	return c.Publish("store.executor.finished", c.Event)
}

func handleExecutionDone(c *processor.Context) error {
	splitState, err := actions.GetSplitState(c.Event)
	if err != nil {
		log.Error(err)
		c.Publish("store.executor.finished", c.Event)
		return err
	}
	if err := actions.UpdateSplitState(c.Event, splitState, domain.Success); err != nil {
		log.Error(err)
		c.Publish("store.executor.finished", c.Event)
		return err
	}
	if splitState.IsComplete() {
		//log.Debug("Dispatching done event")
		return c.Publish("store.executor.finished", c.Event)
	}
	log.Debug(fmt.Sprintf("Supressing event %s on branch %s", c.Event.Name, c.Event.Branch))
	return c.Publish("store.finished", c.Event)
}
