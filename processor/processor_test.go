package processor

import (
	"regexp"
	"testing"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	. "github.com/smartystreets/goconvey/convey"
)

type dispatchMock struct {
	test  *testing.T
	event *domain.Event
}

func (f dispatchMock) Publish(exchange, routingKey string, message interface{}) error {
	f.event = message.(*domain.Event)
	return nil
}

func (f *dispatchMock) Assert(fnc func(*domain.Event) error) {
	if fnc(f.event) != nil {
		f.test.Fail()
	}
}

func TestShouldAssertRegex(t *testing.T) {
	Convey("Should test correct pattern matching", t, func() {
		matched, err := regexp.MatchString(convertPattern("*.persist.request"), "ec498841-59e5-47fd-8075-136d79155705.persist.request")
		if err != nil {
			t.Fail()
		}
		if !matched {
			t.Fail()
		}
	})
}
