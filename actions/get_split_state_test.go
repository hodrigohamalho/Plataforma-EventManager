package actions

import (
	"fmt"
	"testing"

	"github.com/ONSBR/Plataforma-EventManager/domain"

	"github.com/PMoneda/http"
	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldGetSplitState(t *testing.T) {
	Convey("should retrieve current execution state", t, func() {

		mock := http.ReponseMock{
			Method: "GET",
			URL:    "*",
			ReponseBody: `
[{
	"tag" : "f56cfca7-6282-11e8-8808-0242ac12000c",
	"events" : [
		{
			"name" : "create.client.request",
			"eventOut" : "cadastro-cliente.done",
			"version" : "135b9ce9-b73d-4e86-b807-8df8fff73fb9",
			"tag" : "f56cfca7-6282-11e8-8808-0242ac12000c",
			"status" : "success",
			"branch" : "master",
			"scope" : "execution"
		}
	]
}]
			`,
		}

		http.With(t, func(ctx *http.MockContext) {
			ctx.RegisterMock(&mock)
			evt := domain.Event{
				Tag: "f56cfca7-6282-11e8-8808-0242ac12000c",
			}
			splitState, err := GetSplitState(&evt)
			if err != nil {
				t.Fail()
			}

			if !splitState.IsComplete() {
				t.Fail()
			}
		})

	})
}

func TestShouldReturnErrorWhenContractIsNotValid(t *testing.T) {
	Convey("should return error when contract is not valid", t, func() {

		mock := http.ReponseMock{
			Method: "GET",
			URL:    "*",
			ReponseBody: `
{
	"tag" : "f56cfca7-6282-11e8-8808-0242ac12000c",
	"events" : [
		{
			"name" : "create.client.request",
			"eventOut" : "cadastro-cliente.done",
			"version" : "135b9ce9-b73d-4e86-b807-8df8fff73fb9",
			"tag" : "f56cfca7-6282-11e8-8808-0242ac12000c",
			"status" : "success",
			"branch" : "master",
			"scope" : "execution"
		}
	]
}
			`,
		}

		http.With(t, func(ctx *http.MockContext) {
			ctx.RegisterMock(&mock)
			evt := domain.Event{
				Tag: "f56cfca7-6282-11e8-8808-0242ac12000c",
			}
			_, err := GetSplitState(&evt)
			if err == nil {
				t.Fail()
			}
		})

	})
}

func TestShouldNotFindSplitState(t *testing.T) {
	Convey("should not find split state", t, func() {
		mock := http.ReponseMock{
			Method:      "GET",
			ReponseBody: `[]`,
			URL:         "*",
		}
		http.With(t, func(ctx *http.MockContext) {
			ctx.RegisterMock(&mock)
			evt := domain.Event{
				Tag: "f56cfca7-6282-11e8-8808-0242ac12000c",
			}
			_, err := GetSplitState(&evt)
			if err == nil {
				ctx.Fail()
			}
		})
	})
}

func TestShouldReturnError(t *testing.T) {
	Convey("should return error on request", t, func() {
		mock := http.ReponseMock{
			Method:        "GET",
			URL:           "*",
			ResponseError: fmt.Errorf("error"),
		}
		http.With(t, func(ctx *http.MockContext) {
			ctx.RegisterMock(&mock)
			evt := domain.Event{
				Tag: "f56cfca7-6282-11e8-8808-0242ac12000c",
			}
			_, err := GetSplitState(&evt)
			if err == nil {
				ctx.Fail()
			}
		})
	})
}
