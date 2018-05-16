package domain

import (
	"time"

	"github.com/ONSBR/Plataforma-EventManager/infra"
)

//Event define a basic platform event contract
type Event struct {
	Timestamp     string                 `json:"timestamp"`
	Branch        string                 `json:"branch"`
	Name          string                 `json:"name,omitempty"`
	AppOrigin     string                 `json:"appOrigin,omitempty"`
	Owner         string                 `json:"owner,omitempty"`
	InstanceID    string                 `json:"instanceId,omitempty"`
	Scope         string                 `json:"scope,omitempty"`
	ReferenceDate string                 `json:"referenceDate,omitempty"`
	Payload       map[string]interface{} `json:"payload,omitempty"`
	Reproduction  map[string]interface{} `json:"reproduction,omitempty"`
	Reprocessing  map[string]interface{} `json:"reprocessing,omitempty"`
	Bindings      []*Operation
	Commands      []*Command
}

type Command struct {
	Event
}

func NewEvent() *Event {
	event := new(Event)
	event.Commands = make([]*Command, 0, 0)
	event.Bindings = make([]*Operation, 0, 0)
	return event
}

func (e *Event) IsValid() error {

	_, err := time.Parse(time.RFC3339, e.ReferenceDate)
	if err != nil && e.ReferenceDate != "" {
		return infra.NewArgumentException(err.Error())
	}
	return nil
}

func (e *Event) GetReferenceDate() (time.Time, error) {
	return time.Parse(time.RFC3339, e.ReferenceDate)
}

func (e *Event) WillDispatchReprocessing() bool {
	t, _ := e.GetReferenceDate()
	return !time.Time(t).IsZero() && time.Now().UTC().After(time.Time(t).UTC())
}

//ToCeleryMessage transform event to a celery compatible message
func (e *Event) ToCeleryMessage() *CeleryMessage {
	return getCeleryMessage(e)
}

//GetCommand returns a command from instance event
func (e *Event) GetCommand() *Command {
	cmd := new(Command)
	cmd.AppOrigin = e.AppOrigin
	cmd.Branch = e.Branch
	cmd.InstanceID = e.InstanceID
	cmd.Name = e.Name
	cmd.Owner = e.Owner
	cmd.Payload = e.Payload
	cmd.Reprocessing = e.Reprocessing
	cmd.Reproduction = e.Reproduction
	cmd.Scope = e.Scope
	return cmd
}

//AppendCommand adds a new command to a commands list
func (e *Event) AppendCommand(command *Command) {
	e.Commands = append(e.Commands, command)
}

//HasCommands returns true if this event has at least one command
func (e *Event) HasCommands() bool {
	return len(e.Commands) > 0
}
