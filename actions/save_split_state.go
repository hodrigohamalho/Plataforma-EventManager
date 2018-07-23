package actions

import (
	"fmt"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/infra"
	"github.com/ONSBR/Plataforma-EventManager/sdk"
	log "github.com/sirupsen/logrus"
)

//SaveSplitState persist events executing state on process memory to track execution flow
func SaveSplitState(events []*domain.Event) error {
	splitState := domain.NewSplitState()
	if len(events) == 0 {
		return infra.NewArgumentException("cannot save split state document with no events")
	}
	splitState.Tag = events[0].Tag
	for _, event := range events {
		state := event.ToEventState()
		if len(event.Bindings) > 0 {
			operation := event.Bindings[0]
			state.EventOut = operation.EventOut
			state.Version = operation.Version
		} else {
			log.Error(fmt.Sprintf("Event %s on branch %s with scope %s has no binding", event.Name, event.Branch, event.Scope))
		}
		splitState.AddEventState(state)
	}
	return sdk.SaveDocument("splitted_events", splitState)
}
