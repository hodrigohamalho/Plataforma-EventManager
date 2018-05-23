package actions

import (
	"testing"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldRefuseEvent(t *testing.T) {
	Convey("Should refuse event when has no binding", t, func() {
		event := domain.NewEvent()
		refused, _ := RefuseEvent(event)
		if !refused {
			t.Fail()
		}
	})

	Convey("Should not refuse when event is a system event", t, func() {
		for _, sysEventName := range domain.SystemEvents {
			evt := domain.NewEvent()
			evt.Name = sysEventName
			refused, _ := RefuseEvent(evt)
			if refused {
				t.Fail()
			}
		}
	})
}
