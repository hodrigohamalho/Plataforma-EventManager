package processor

import (
	"regexp"
	"testing"

	"github.com/ONSBR/Plataforma-EventManager/bus"
	"github.com/ONSBR/Plataforma-EventManager/domain"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldAssertRegex(t *testing.T) {
	Convey("should test correct pattern matching", t, func() {
		matched, err := regexp.MatchString(convertPattern("*.persist.request"), "ec498841-59e5-47fd-8075-136d79155705.persist.request")
		So(err, ShouldBeNil)
		So(matched, ShouldBeTrue)
	})
}

func TestShouldBuildProcessor(t *testing.T) {
	mockDispatcher := bus.DispatcherMock{}

	evt := domain.NewEvent()
	Convey("should push an event to processor", t, func() {
		proc := NewProcessor(&mockDispatcher)
		So(proc.Push(evt), ShouldBeNil)
	})

	Convey("should route an event to some listener", t, func() {
		evt.Name = "event"
		routed := false
		proc := NewProcessor(&mockDispatcher)
		proc.When("event", func(ctx *Context) error {
			So(ctx.Event.Name, ShouldEqual, evt.Name)
			routed = true
			return nil
		})
		proc.When("*", func(ctx *Context) error {
			t.Fail()
			return nil
		})
		proc.Push(evt)
		So(routed, ShouldBeTrue)
	})

	Convey("should use middleware before route", t, func() {
		evt.Name = "event"
		routed := false
		proc := NewProcessor(&mockDispatcher)
		proc.Use("event", func(ctx *Context) error {
			ctx.Event.Tag = "my_tag"
			return nil
		})
		proc.When("event", func(ctx *Context) error {
			routed = true
			So(ctx.Event.Tag, ShouldEqual, "my_tag")
			return nil
		})
		proc.When("*", func(ctx *Context) error {
			routed = false
			return nil
		})
		proc.Push(evt)
		So(routed, ShouldBeTrue)
	})
}
