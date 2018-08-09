package handlers

import (
	"fmt"

	"github.com/ONSBR/Plataforma-EventManager/actions"
	"github.com/ONSBR/Plataforma-EventManager/processor"
	log "github.com/sirupsen/logrus"
)

//HandleGeneralEvent handle general event
func HandleGeneralEvent(c *processor.Context) error {
	log.Debug(fmt.Sprintf("HandleGeneralEvent %s on branch %s with scope %s", c.Event.Name, c.Event.Branch, c.Event.Scope))
	if c.Event.IsReproduction() {
		return c.Publish("store.executor", c.Event)
	}
	if c.Event.IsReprocessing() {
		return handleReprocessingGeneralEvent(c)
	}
	if c.Event.IsExecution() {
		return handleExecutionGeneralEvent(c)
	}
	return nil
}

func handleExecutionGeneralEvent(c *processor.Context) error {
	if events, err := actions.SplitEvent(c.Event); err != nil {
		log.Error(err)
		return err
	} else if err := actions.SaveSplitState(events); err != nil {
		log.Error(err)
		return err
	} else {
		for _, event := range events {
			if err := c.Publish("store.executor", event); err != nil {
				log.Error(err)
				return err
			}
		}
	}
	return nil
}

func handleReprocessingGeneralEvent(c *processor.Context) error {
	return c.Publish("store.executor", c.Event)
}
