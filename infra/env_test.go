package infra

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldGetEnvOrDefault(t *testing.T) {
	Convey("should get env var or default", t, func() {
		os.Setenv("TESTE", "CURRENT")
		So(GetEnv("TESTE", "DEFAULT"), ShouldEqual, "CURRENT")
		So(GetEnv("NOT_EXIST", "DEFAULT"), ShouldEqual, "DEFAULT")
	})
}
