package domain

import (
	"strings"

	"github.com/google/uuid"
)

//TODO maybe will be better put this events on apicore
var systemEvents = []string{
	"system.reprocessing.error",
	"system.executor.enable.debug",
	"system.executor.disable.debug",
	"system.process.persist.error",
	"system.events.reprocessing.request",
	"system.events.reproduction.request",
}

//Event define a basic platform event contract
type Event struct {
	Timestamp    string                 `json:"timestamp"`
	Branch       string                 `json:"branch"`
	Name         string                 `json:"name,omitempty"`
	Tag          string                 `json:"tag"`
	AppOrigin    string                 `json:"appOrigin,omitempty"`
	Owner        string                 `json:"owner,omitempty"`
	InstanceID   string                 `json:"instanceId,omitempty"`
	Scope        string                 `json:"scope,omitempty"`
	Payload      map[string]interface{} `json:"payload,omitempty"`
	Reproduction map[string]interface{} `json:"reproduction,omitempty"`
	Reprocessing map[string]interface{} `json:"reprocessing,omitempty"`
	Bindings     []*Operation
	Commands     []*Command
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

func (e *Event) IsEndingEvent() bool {
	return strings.HasSuffix(e.Name, ".done") || strings.HasSuffix(e.Name, ".error") || strings.HasSuffix(e.Name, ".exception")
}

//ToCeleryMessage transform event to a celery compatible message
func (e *Event) ToCeleryMessage() *CeleryMessage {
	return getCeleryMessage(e)
}

//ApplyDefaultFields apply default fields for branch, scope and tag
func (e *Event) ApplyDefaultFields() {
	if e.Branch == "" {
		e.Branch = "master"
	}
	if e.Scope == "" {
		e.Scope = "execution"
	}
	if e.Tag == "" {
		uuid, _ := uuid.NewUUID()
		e.Tag = uuid.String()
	}
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

//IsSystemEvent returns true if this event is a internal platform event
func (e *Event) IsSystemEvent() bool {
	for _, sysEvt := range systemEvents {
		if sysEvt == e.Name {
			return true
		}
	}
	return false
}
