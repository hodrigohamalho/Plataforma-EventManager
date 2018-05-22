package flow

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ONSBR/Plataforma-EventManager/bus"
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/infra"
	"github.com/ONSBR/Plataforma-EventManager/lock"
	"github.com/ONSBR/Plataforma-EventManager/processor"
	"github.com/ONSBR/Plataforma-EventManager/sdk"
	log "github.com/sirupsen/logrus"
)

var broker bus.Dispatcher
var systemEvents = []string{"system.reprocessing.error", "system.executor.enable.debug", "system.executor.disable.debug", "system.process.persist.error", "system.events.reprocessing.request", "system.events.reproduction.request"}

//GetStoreEventFlow for just storing event on event store
func GetStoreEventFlow(dispatcher bus.Dispatcher) *processor.Processor {
	p := processor.NewProcessor(dispatcher)

	p.CutOff(eventNotRegistered)

	p.Where("system.*").Execute(checkSystemEvent).Dispatch("executor.store")

	p.Where("*").Dispatch("store")
	return p
}

//GetEventFlow returns a processor with events flow applied
func GetEventFlow(dispatcher bus.Dispatcher) *processor.Processor {
	broker = dispatcher

	if err := swapPersistEventToExecutorQueue(); err != nil {
		log.Error(err)
	}
	p := processor.NewProcessor(dispatcher)

	p.CutOff(eventNotRegistered)

	p.Where("*.persist.request").Execute(handlePersistEvent).Dispatch("persist.inexecution")

	p.Where("*.exception").Dispatch("exception.executor")

	p.Where("*.done").Execute(handleReprocessing).Dispatch("executor.finished")

	p.Where("system.process.persist.error").Execute(handlePersistenceDone).Dispatch("persist_error")

	p.Where("system.*").Execute(checkSystemEvent).Dispatch("executor")

	p.Where("*").Execute(checkPlatformAvailability).Dispatch("executor")

	return p
}

/*
Event Handlers
*/

func eventNotRegistered(event *domain.Event) (err error) {
	if isSystemEvent(event) {
		return nil
	}
	if len(event.Bindings) == 0 {
		err = infra.NewSubscriberNotFoundException(fmt.Sprintf("Event %s has no subscribers", event.Name))
	}
	return
}

func swapPersistEventToExecutorQueue() error {
	log.Debug("Swapping persist event to executor queue")
	err := broker.Swap(bus.EVENT_PERSIST_QUEUE, "executor.store")
	if err != nil {
		return err
	}
	_, err = broker.Pop(bus.EVENT_PERSIST_REQUEST_QUEUE)
	return err
}

func handlePersistenceDone(event *domain.Event) error {
	return swapPersistEventToExecutorQueue()
}

func handleReprocessing(event *domain.Event) error {
	err := swapPersistEventToExecutorQueue()
	if event.Scope != "reprocessing" {
		return err
	} else if err != nil {
		return err
	}
	reprocessingOrigin := new(domain.ReprocessingData)
	reprocessingOrigin.ParseEvent(event)
	if response, err := sdk.GetDocument("reprocessing", map[string]string{"id": reprocessingOrigin.ID}); err != nil {
		return infra.NewComponentException(err.Error())
	} else {
		reprocessingList := make([]domain.Reprocessing, 0, 1)
		if err := json.Unmarshal([]byte(response), &reprocessingList); err != nil {
			return infra.NewComponentException(err.Error())
		}
		reprocessing := reprocessingList[0]
		doneReprocessing := true
		for _, evt := range reprocessing.Events {
			reprocessingData := new(domain.ReprocessingData)
			reprocessingData.ParseEvent(evt)
			if reprocessingData.Tag == reprocessingOrigin.Tag && reprocessingOrigin.Branch == reprocessingData.Branch {
				reprocessingData.Executed = true
				evt.Reprocessing["executed"] = true
			}
			doneReprocessing = doneReprocessing && reprocessingData.Executed
		}
		reprocessing.Done = doneReprocessing
		sdk.ReplaceDocument("reprocessing", map[string]string{"id": reprocessingOrigin.ID}, reprocessing)
		if doneReprocessing {
			log.Info(fmt.Sprintf("Reprocessing %s already done", reprocessingOrigin.ID))
			event.Branch = "master"
			return nil
		} else {
			log.Info(fmt.Sprintf("Suppressing event %s in scope %s in branch %s", event.Name, event.Scope, event.Branch))
		}
	}
	return infra.NewRunningReprocessingException("reprocessing still running")
}

func handlePersistEvent(event *domain.Event) error {
	_, err := broker.First(bus.EVENT_PERSIST_REQUEST_QUEUE)
	if err != nil && err.Error() == infra.PersistEventQueueEmpty {
		broker.Publish("executor.store.inexecution", event.ToCeleryMessage())
		return err
	} else if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func handleUnlock(event *domain.Event) error {
	eventNameParts := strings.Split(event.Name, ".")
	solutionID := eventNameParts[0]
	return lock.Unlock(solutionID)
}

func checkSystemEvent(e *domain.Event) error {
	for _, sysEvt := range systemEvents {
		if sysEvt == e.Name {
			return nil
		}
	}
	return infra.NewArgumentException(fmt.Sprintf("Event %s is not a valid platform event", e.Name))
}

func isSystemEvent(e *domain.Event) bool {
	for _, sysEvt := range systemEvents {
		if sysEvt == e.Name {
			return true
		}
	}
	return false
}

func checkPlatformAvailability(event *domain.Event) error {
	if solutionID, err := sdk.GetSolutionIDByEventName(event.Name); err == nil {
		if locked, err := lock.SolutionIsLocked(solutionID); err != nil {
			return infra.NewComponentException(err.Error())
		} else if locked {
			return infra.NewPlatformLockedException(fmt.Sprintf("solution %s is locked by reprocessing", solutionID))
		}
	} else {
		return err
	}
	return nil
}
