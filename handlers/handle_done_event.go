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
	err := actions.SwapPersistEventToExecutorQueue(c.Dispatcher())
	if err != nil {
		log.Error(err)
		return err
	}
	log.Debug(fmt.Sprintf("HandleDoneEvent %s on branch %s", c.Event.Name, c.Event.Branch))
	splitState, err := actions.GetSplitState(c.Event)
	if err != nil {
		log.Error(err)
		return err
	}
	if err := actions.UpdateSplitState(c.Event, splitState, domain.Success); err != nil {
		log.Error(err)
		return err
	}
	if splitState.IsComplete() {
		log.Info("Dispatching done event")
		return c.Publish("store.executor.finished", c.Event)
	}
	log.Info(fmt.Sprintf("Supressing event %s on branch %s", c.Event.Name, c.Event.Branch))
	return c.Publish("store", c.Event)
}
