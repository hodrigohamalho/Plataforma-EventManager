package actions

import (
	"strings"
	"testing"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/PMoneda/http"
	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldSetFailureProcess(t *testing.T) {
	Convey("should update process instance status to failure", t, func() {
		mock := http.ReponseMock{
			Method: "POST",
			URL:    "*",
		}
		http.With(t, func(ctx *http.MockContext) {
			ctx.RegisterMock(&mock)
			event := domain.NewEvent()
			event.Name = "test"
			event.InstanceID = "1"
			event.Payload["instance_id"] = "1"
			SetFailureProcess(event)
			So(strings.Contains(mock.RequestBody(), `"status":"failure"`), ShouldBeTrue)
			So(mock.CalledTimes(), ShouldBeGreaterThan, 0)
		})
	})
}
