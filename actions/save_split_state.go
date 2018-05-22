package actions

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/sdk"
)

//SaveSplitState persist events executing state on process memory to track execution flow
func SaveSplitState(events []*domain.Event) error {
	return sdk.SaveDocument("bla", events)
}
