package actions

import (
	"errors"
	"testing"

	"github.com/ONSBR/Plataforma-EventManager/domain"

	"github.com/ONSBR/Plataforma-EventManager/bus"
	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldSwapEvents(t *testing.T) {
	Convey("should swap persistence events", t, func() {
		dispatcher := new(bus.DispatcherMock)

		dispatcher.OnSwap = func(from, to string) error {
			So(from, ShouldEqual, bus.EventPersistQueue)
			So(to, ShouldEqual, "executor.store")
			return nil
		}
		dispatcher.OnPop = func(queue string) (*domain.Event, error) {
			So(bus.EventPersistRequestQueue, ShouldEqual, queue)
			return domain.NewEvent(), nil
		}
		SwapPersistEventToExecutorQueue(dispatcher)
	})

	Convey("should return error when swap crashes", t, func() {
		dispatcher := new(bus.DispatcherMock)

		dispatcher.OnSwap = func(from, to string) error {
			return errors.New("error")
		}
		err := SwapPersistEventToExecutorQueue(dispatcher)
		So(err, ShouldNotBeNil)
	})
}
