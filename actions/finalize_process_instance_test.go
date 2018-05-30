package actions

import (
	"strings"
	"testing"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/PMoneda/http"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldFinalizeProcessInstance(t *testing.T) {
	Convey("should finalize process on apicore", t, func() {
		mock := http.ReponseMock{
			Method: "POST",
			URL:    "*",
		}
		http.With(t, func(ctx *http.MockContext) {
			ctx.RegisterMock(&mock)
			event := domain.NewEvent()
			event.Name = "test"
			event.Payload["instance_id"] = "1"
			FinalizeProcessInstance(event)

			So(strings.Contains(mock.RequestBody(), `"status":"finished"`), ShouldBeTrue)
			So(mock.CalledTimes(), ShouldBeGreaterThan, 0)

		})
	})
}
