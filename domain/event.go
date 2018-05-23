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
}

//NewEvent creates a new Event Instance
func NewEvent() *Event {
	event := new(Event)
	event.Bindings = make([]*Operation, 0, 0)
	return event
}

//IsEndingEvent checkis event's name ends with .done, .error or .exception
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

//IsSystemEvent returns true if this event is a internal platform event
func (e *Event) IsSystemEvent() bool {
	for _, sysEvt := range systemEvents {
		if sysEvt == e.Name {
			return true
		}
	}
	return false
}

//ToEventState converts an event to state event
func (e *Event) ToEventState() EventState {
	return EventState{
		Branch: e.Branch,
		Name:   e.Name,
		Scope:  e.Scope,
		Status: Pending,
	}
}
