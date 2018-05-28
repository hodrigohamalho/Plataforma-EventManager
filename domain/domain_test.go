package domain

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func pickEvent() *Event {
	event := NewEvent()
	event.Name = "a"
	event.Branch = "a"
	event.Scope = "execution"
	return event
}

func TestShouldConvertEventToEventState(t *testing.T) {
	event := pickEvent()
	Convey("Should get event state from domain event", t, func() {
		state := event.ToEventState()
		if state.Branch != event.Branch {
			t.Fail()
		}
		if state.Name != event.Name {
			t.Fail()
		}
		if state.Scope != event.Scope {
			t.Fail()
		}
		if state.Status != Pending {
			t.Fail()
		}

	})
}

func TestShouldTransformEventToCeleryMessage(t *testing.T) {
	Convey("Should transform domain evento to a celery task", t, func() {
		event := pickEvent()
		celeryMessage := event.ToCeleryMessage()
		if celeryMessage.Task != "tasks.process" {
			t.Fail()
		}
		if len(celeryMessage.Args) == 0 {
			t.Fail()
		}
	})
}

func TestShouldVerifyIfEventIsAnEndingEvent(t *testing.T) {
	Convey("Should assert that events ending with .done, .error and .exception is an ending event", t, func() {
		event := pickEvent()
		if event.IsEndingEvent() {
			t.Fail()
		}
		event.Name = "a.done"
		if !event.IsEndingEvent() {
			t.Fail()
		}
		event.Name = "a.error"
		if !event.IsEndingEvent() {
			t.Fail()
		}
		event.Name = "a.exception"
		if !event.IsEndingEvent() {
			t.Fail()
		}
	})
}

func TestShouldTestEventScope(t *testing.T) {
	Convey("should verify if an event is a reprocessing", t, func() {
		event := pickEvent()
		event.Scope = "reprocessing"
		if !event.IsReprocessing() {
			t.Fail()
		}
	})

	Convey("should verify if an event is a execution", t, func() {
		event := pickEvent()
		if !event.IsExecution() {
			t.Fail()
		}

		event.Scope = "execution"

		if !event.IsExecution() {
			t.Fail()
		}

	})

	Convey("should verify if an event is a reproduction", t, func() {
		event := pickEvent()
		event.Scope = "reproduction"
		if !event.IsReproduction() {
			t.Fail()
		}
	})
}

func TestShouldApplyDefaultFields(t *testing.T) {
	Convey("should apply default fields in event", t, func() {
		event := NewEvent()
		event.ApplyDefaultFields()

		if event.Scope != "execution" {
			t.Fail()
		}

		if event.Branch != "master" {
			t.Fail()
		}

		if event.Tag == "" {
			t.Fail()
		}
	})
}

func TestShouldCheckIfEventIsSystemEvent(t *testing.T) {
	Convey("should check if an event is a system event", t, func() {
		for _, e := range SystemEvents {
			evt := NewEvent()
			evt.Name = e
			if !evt.IsSystemEvent() {
				t.Fail()
			}
		}

		evt := NewEvent()
		evt.Name = "test"
		if evt.IsSystemEvent() {
			t.Fail()
		}
	})
}
