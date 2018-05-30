package actions

import (
	"testing"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldUpdateSplitState(t *testing.T) {
	Convey("should update split state", t, func() {
		evt := domain.NewEvent()
		evt.Branch = "master"
		evt.Name = "name"
		spState := domain.NewSplitState()
		state := evt.ToEventState()
		state.EventOut = "name"
		spState.AddEventState(state)
		err := UpdateSplitState(evt, spState, "success")
		So(err, ShouldBeNil)
		So(spState.Events[0].Status, ShouldEqual, "success")
	})
}
