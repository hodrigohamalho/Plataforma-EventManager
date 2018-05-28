package actions

import (
	"testing"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldRefuseEvent(t *testing.T) {
	Convey("should refuse event when has no binding", t, func() {
		event := domain.NewEvent()
		refused, _ := RefuseEvent(event)
		if !refused {
			t.Fail()
		}
	})

	Convey("should not refuse when event is a system event", t, func() {
		for _, sysEventName := range domain.SystemEvents {
			evt := domain.NewEvent()
			evt.Name = sysEventName
			refused, _ := RefuseEvent(evt)
			if refused {
				t.Fail()
			}
		}
	})

	Convey("should not refuse when event has binding", t, func() {
		evt := domain.NewEvent()
		evt.Bindings = make([]*domain.Operation, 1)
		refused, _ := RefuseEvent(evt)
		if refused {
			t.Fail()
		}
	})
}
