package flow

import (
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

func GetStoreEventFlow(dispatcher bus.Dispatcher) *processor.Processor {
	p := processor.NewProcessor(dispatcher)

	p.CutOff(eventNotRegistered)

	p.Where("system.*").Execute(checkSystemEvent).Dispatch("executor.store")

	p.Where("*").Execute(checkPlatformAvailability).Dispatch("store")
	return p
}

var systemEvents = []string{"system.executor.enable.debug", "system.executor.disable.debug", "system.process.persist.error"}

//GetEventFlow returns a processor with events flow applied
func GetEventFlow(dispatcher bus.Dispatcher) *processor.Processor {
	broker = dispatcher

	if err := swapPersistEventToExecutorQueue(); err != nil {
		log.Error(err)
	}
	p := processor.NewProcessor(dispatcher)

	p.CutOff(eventNotRegistered)

	p.Where("*.persist.request").Execute(handlePersistEvent).Dispatch("persist")

	p.Where("*.persist.unlock").Execute(handleUnlock).Dispatch("store")

	p.Where("*.exception").Dispatch("exception.store")

	p.Where("*.done").Execute(handlePersistenceDone).Dispatch("store.executor")

	p.Where("system.process.persist.error").Execute(handlePersistenceDone).Dispatch("store.persist_error")

	p.Where("system.*").Execute(checkSystemEvent).Dispatch("executor.store")

	p.Where("*").Execute(checkPlatformAvailability).Dispatch("executor.store")

	return p
}

/*
Event Handlers
*/

func eventNotRegistered(event *domain.Event) error {
	if isSystemEvent(event) {
		return nil
	}
	_, err := sdk.EventHasBindings(event.Name)
	return err
}

func swapPersistEventToExecutorQueue() error {
	log.Info("Swapping persist event to executor queue")
	return broker.Swap(bus.EVENT_PERSIST_QUEUE, "executor.store")
}

func handlePersistenceDone(event *domain.Event) error {
	return swapPersistEventToExecutorQueue()
}

func handlePersistEvent(event *domain.Event) error {
	eventNameParts := strings.Split(event.Name, ".")
	solutionID := eventNameParts[0]
	if locked, err := lock.SolutionIsLocked(solutionID); err != nil {
		return infra.NewComponentException(err.Error())
	} else if locked {
		return infra.NewPlatformLockedException(fmt.Sprintf("solution %s is locked by reprocessing", solutionID))
	} else if event.WillDispatchReprocessing() {
		log.Info("Locking solution", event)
		return lock.Lock(solutionID, event)
	}
	err := broker.Get(bus.EVENT_PERSIST_QUEUE, func(queuedEvent *domain.Event) error {
		return fmt.Errorf("ignore")
	})
	if err != nil && err.Error() == infra.PersistEventQueueEmpty {
		broker.Publish("persist.executor.store", event.ToCeleryMessage())
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
