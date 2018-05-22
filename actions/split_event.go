package actions

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/infra"
	"github.com/ONSBR/Plataforma-EventManager/sdk"
)

//SplitEvent splits an event into a commands based on opened branches
func SplitEvent(event *domain.Event) ([]*domain.Event, error) {
	commands := []*domain.Event{event}
	if !isSplitable(event) {
		return commands, nil
	}
	if branches, err := sdk.GetOpenBranchesBySystem(event.Bindings[0].SystemID); err != nil {
		return nil, err
	} else {
		for _, branch := range branches {
			command := new(domain.Event)
			if err := infra.Clone(event, command); err != nil {
				return nil, err
			}
			command.Branch = branch.Name
			commands = append(commands, command)
		}
	}
	return commands, nil
}

func isSplitable(event *domain.Event) bool {
	return event.Scope == "execution"
}
