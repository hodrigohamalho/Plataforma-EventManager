package actions

import (
	"fmt"
	"testing"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/PMoneda/http"
	. "github.com/smartystreets/goconvey/convey"
)

const openBranches = `
[{ "name" : "branch-01" }, { "name" : "branch-02" }]
`

func TestShouldSplitEvents(t *testing.T) {
	Convey("should split events into commands", t, func() {
		http.With(t, func(ctx *http.MockContext) {
			ctx.RegisterMock(&http.ReponseMock{
				Method:      "GET",
				URL:         "*",
				ReponseBody: openBranches,
			})
			event := domain.NewEvent()
			event.Name = "test.event"
			event.Scope = "execution"
			event.Bindings = append(event.Bindings, &domain.Operation{Reprocessable: true})
			evts, err := SplitEvent(event)
			So(err, ShouldBeNil)
			So(len(evts), ShouldEqual, 3)
		})

	})

	Convey("should return single command when event is not splitable", t, func() {
		http.With(t, func(ctx *http.MockContext) {
			ctx.RegisterMock(&http.ReponseMock{
				Method:      "GET",
				URL:         "*",
				ReponseBody: openBranches,
			})
			event := domain.NewEvent()
			evts, err := SplitEvent(event)
			So(err, ShouldBeNil)
			So(len(evts), ShouldEqual, 1)
		})

	})

	Convey("should return error on error request", t, func() {
		http.With(t, func(ctx *http.MockContext) {
			ctx.RegisterMock(&http.ReponseMock{
				Method:        "GET",
				URL:           "*",
				ResponseError: fmt.Errorf("error"),
			})
			event := domain.NewEvent()
			event.Name = "test.event"
			event.Scope = "execution"
			event.Bindings = append(event.Bindings, &domain.Operation{Reprocessable: true})
			evts, err := SplitEvent(event)
			So(err, ShouldNotBeNil)
			So(len(evts), ShouldEqual, 0)
		})

	})

}
