package infra

import "testing"
import . "github.com/smartystreets/goconvey/convey"

func TestShouldCloneObjects(t *testing.T) {

	type Test struct {
		ID string
	}
	from := new(Test)
	from.ID = "A"

	to := new(Test)
	Convey("should clone objects", t, func() {
		err := Clone(from, to)
		So(err, ShouldBeNil)
		So(to.ID, ShouldEqual, from.ID)
	})

}
