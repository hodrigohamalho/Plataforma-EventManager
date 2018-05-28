package actions

import (
	"testing"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/PMoneda/http"
	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldSaveSplitState(t *testing.T) {
	Convey("should save split state", t, func() {
		http.With(t, func(ctx *http.MockContext) {
			mock := http.ReponseMock{
				URL:    "*",
				Method: "POST",
			}
			ctx.RegisterMock(&mock)
			events := make([]*domain.Event, 1, 1)
			events[0] = domain.NewEvent()
			events[0].Bindings = append(events[0].Bindings, new(domain.Operation))

			err := SaveSplitState(events)
			if err != nil {
				t.Fail()
			}
		})
	})

	Convey("should return error when event list is empty", t, func() {

		err := SaveSplitState(make([]*domain.Event, 0))
		if err == nil {
			t.Fail()
		}
	})
}
