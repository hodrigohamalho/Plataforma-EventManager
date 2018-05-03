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

func (e *Event) ToCeleryMessage() *CeleryMessage {
	return getCeleryMessage(e)
}
