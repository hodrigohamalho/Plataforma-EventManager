package actions

import (
	"testing"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/PMoneda/http"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldPushEventToEventStore(t *testing.T) {
	Convey("should push event to event store", t, func() {
		mock := http.ReponseMock{
			Method: "POST",
			URL:    "*",
		}
		http.With(t, func(ctx *http.MockContext) {
			ctx.RegisterMock(&mock)
			event := domain.NewEvent()
			PushEventToEventStore(event)
			if mock.CalledTimes() == 0 {
				t.Fail()
			}
		})
	})
}
