package util

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldHashSha1(t *testing.T) {
	Convey("should hash any object", t, func() {
		s, err := SHA1("hello")
		So(err, ShouldBeNil)

		s1, err := SHA1("hello")
		So(err, ShouldBeNil)
		So(s, ShouldEqual, s1)
	})

}
