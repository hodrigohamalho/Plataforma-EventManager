package processor

import (
	"fmt"
	"testing"

	"github.com/ONSBR/Plataforma-EventManager/domain"
)

type dispatchMock struct {
	test  *testing.T
	event *domain.Event
}

func (f dispatchMock) Publish(exchange, routingKey string, message interface{}) error {
	f.event = message.(*domain.Event)
	return nil
}

func (f *dispatchMock) Assert(fnc func(*domain.Event) error) {
	if fnc(f.event) != nil {
		f.test.Fail()
	}
}

func TestShouldExecuteProcessorFlow(t *testing.T) {

	event := new(domain.Event)
	event.Name = "123121324234.domain.persist"
	dispatcher := new(dispatchMock)
	dispatcher.test = t
	p := NewProcessor(*dispatcher)
	p.Where("*.persist").Dispatch("exchange_persistencia", "#")
	if err := p.Push(event); err != nil {
		t.Fail()
	}
}

func TestShouldNotExecuteProcessorFlow(t *testing.T) {

	event := new(domain.Event)
	event.Name = "123121324234.domain.persist"
	dispatcher := new(dispatchMock)
	dispatcher.test = t
	p := NewProcessor(*dispatcher)
	p.Where("*.bla").Dispatch("exchange_persistencia", "#")
	if err := p.Push(event); err == nil {
		t.Fail()
	}
}

func TestShouldChangeEventName(t *testing.T) {

	event := new(domain.Event)
	event.Name = "123121324234.domain.persist"
	dispatcher := new(dispatchMock)
	dispatcher.test = t
	p := NewProcessor(*dispatcher)
	p.Where("*.persist").Execute(func(evt *domain.Event) error {
		event.Name = "changed.name"
		return nil
	}).Dispatch("exchange_persistencia", "#")

	if err := p.Push(event); err != nil {
		t.Fail()
	}

	if event.Name != "changed.name" {
		t.Fail()
	}
}

func TestShouldStopFlow(t *testing.T) {

	event := new(domain.Event)
	event.Name = "123121324234.domain.persist"
	dispatcher := new(dispatchMock)
	dispatcher.test = t
	p := NewProcessor(*dispatcher)
	p.Where("*.persist").Execute(func(evt *domain.Event) error {
		return fmt.Errorf("Abort execution")
	}).Dispatch("exchange_persistencia", "#")

	if err := p.Push(event); err == nil {
		t.Fail()
	}

	dispatcher.Assert(func(e *domain.Event) error {
		fmt.Println(e)
		return fmt.Errorf("abort")
	})
}
