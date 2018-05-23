package actions

import (
	"strings"
	"testing"

	"github.com/ONSBR/Plataforma-EventManager/domain"

	"github.com/ONSBR/Plataforma-EventManager/clients/http"
	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldFinalizeProcessInstance(t *testing.T) {
	Convey("Should finalize process on apicore", t, func() {
		mock := http.ReponseMock{
			Method: "POST",
			URL:    "*",
		}
		http.RegisterMock(&mock)
		event := domain.NewEvent()
		event.Name = "test"
		event.Payload["instance_id"] = "1"
		FinalizeProcessInstance(event)
		if !strings.Contains(mock.RequestBody(), `"status":"finished"`) {
			t.Fail()
		}

		if mock.CalledTimes() == 0 {
			t.Fail()
		}

	})
}
