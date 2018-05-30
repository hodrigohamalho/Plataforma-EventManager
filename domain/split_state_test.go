package domain

import "testing"
import . "github.com/smartystreets/goconvey/convey"

func TestShouldAppendEventStateOnSplitState(t *testing.T) {
	Convey("should add event state to split state", t, func() {
		spState := NewSplitState()
		evtState := new(EventState)
		spState.AddEventState(evtState)
		evtState.Status = "success"
		So(len(spState.Events), ShouldEqual, 1)
		So(spState.IsComplete(), ShouldBeTrue)
	})
}
