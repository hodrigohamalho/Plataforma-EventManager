package actions

import (
	"testing"

	"github.com/ONSBR/Plataforma-EventManager/clients/http"
	"github.com/ONSBR/Plataforma-EventManager/domain"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldPushEventToEventStore(t *testing.T) {
	Convey("Should push event to event store", t, func() {
		mock := http.ReponseMock{
			Method: "POST",
			URL:    "*",
		}
		http.RegisterMock(&mock)
		event := domain.NewEvent()
		PushEventToEventStore(event)

		if mock.CalledTimes() == 0 {
			t.Fail()
		}

	})
}
