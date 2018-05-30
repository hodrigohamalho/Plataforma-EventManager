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
	Convey("should get event state from domain event", t, func() {
		state := event.ToEventState()
		So(state.Branch, ShouldEqual, event.Branch)
		So(state.Name, ShouldEqual, event.Name)
		So(state.Scope, ShouldEqual, event.Scope)
		So(state.Status, ShouldEqual, Pending)
	})
}

func TestShouldTransformEventToCeleryMessage(t *testing.T) {
	Convey("should transform domain evento to a celery task", t, func() {
		event := pickEvent()
		celeryMessage := event.ToCeleryMessage()
		So(celeryMessage.Task, ShouldEqual, "tasks.process")
		So(len(celeryMessage.Args), ShouldEqual, 1)
	})
}

func TestShouldVerifyIfEventIsAnEndingEvent(t *testing.T) {
	Convey("should assert that events ending with .done, .error and .exception is an ending event", t, func() {
		event := pickEvent()
		So(event.IsEndingEvent(), ShouldBeFalse)
		event.Name = "a.done"
		So(event.IsEndingEvent(), ShouldBeTrue)
		event.Name = "a.error"
		So(event.IsEndingEvent(), ShouldBeTrue)
		event.Name = "a.exception"
		So(event.IsEndingEvent(), ShouldBeTrue)
	})
}

func TestShouldTestEventScope(t *testing.T) {
	Convey("should verify if an event is a reprocessing", t, func() {
		event := pickEvent()
		event.Scope = "reprocessing"
		So(event.IsReprocessing(), ShouldBeTrue)
	})

	Convey("should verify if an event is a execution", t, func() {
		event := pickEvent()
		So(event.IsExecution(), ShouldBeTrue)
		event.Scope = "execution"
		So(event.IsExecution(), ShouldBeTrue)
	})

	Convey("should verify if an event is a reproduction", t, func() {
		event := pickEvent()
		event.Scope = "reproduction"
		So(event.IsReproduction(), ShouldBeTrue)
	})
}

func TestShouldApplyDefaultFields(t *testing.T) {
	Convey("should apply default fields in event", t, func() {
		event := NewEvent()
		event.ApplyDefaultFields()
		So(event.Scope, ShouldEqual, "execution")
		So(event.Branch, ShouldEqual, "master")
		So(event.Tag, ShouldNotEqual, "")
	})
}

func TestShouldCheckIfEventIsSystemEvent(t *testing.T) {
	Convey("should check if an event is a system event", t, func() {
		for _, e := range SystemEvents {
			evt := NewEvent()
			evt.Name = e
			So(evt.IsSystemEvent(), ShouldBeTrue)
		}

		evt := NewEvent()
		evt.Name = "test"
		So(evt.IsSystemEvent(), ShouldBeFalse)
	})
}
