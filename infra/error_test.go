package infra

import "testing"
import . "github.com/smartystreets/goconvey/convey"

func TestShouldCreateNewException(t *testing.T) {
	Convey("should create a new exception object", t, func() {
		ex := NewException("tag", "message")
		So(ex.Error(), ShouldEqual, "message")

	})

	Convey("should create a new subscribe not found exception object", t, func() {
		ex := NewSubscriberNotFoundException("message")
		So(ex.HTTPStatus(), ShouldEqual, 200)
		So(ex.Error(), ShouldEqual, "message")

	})

	Convey("should create a new platform locked exception object", t, func() {
		ex := NewPlatformLockedException("message")
		So(ex.HTTPStatus(), ShouldEqual, mapErrorTagHTTPStatus[PlatformLocked])
		So(ex.Error(), ShouldEqual, "message")
	})

	Convey("should create a new empty queue exception object", t, func() {
		ex := NewEmptyQueueException("message")
		So(ex.HTTPStatus(), ShouldEqual, mapErrorTagHTTPStatus[PersistEventQueueEmpty])
		So(ex.Error(), ShouldEqual, "message")
	})

	Convey("should create a new running reprocessing exception object", t, func() {
		ex := NewRunningReprocessingException("message")
		So(ex.HTTPStatus(), ShouldEqual, mapErrorTagHTTPStatus[RunningReprocessing])
		So(ex.Error(), ShouldEqual, "message")
	})

	Convey("should create a new system exception object", t, func() {
		ex := NewSystemException("message")
		So(ex.HTTPStatus(), ShouldEqual, mapErrorTagHTTPStatus[SystemError])
		So(ex.Error(), ShouldEqual, "message")
	})

	Convey("should create a new argument exception object", t, func() {
		ex := NewArgumentException("message")
		So(ex.HTTPStatus(), ShouldEqual, mapErrorTagHTTPStatus[InvalidArguments])
		So(ex.Error(), ShouldEqual, "message")
	})

	Convey("should create a new component exception object", t, func() {
		ex := NewComponentException("message")
		So(ex.HTTPStatus(), ShouldEqual, mapErrorTagHTTPStatus[PlatformComponentError])
		So(ex.Error(), ShouldEqual, "message")
	})

	Convey("should bind a tag to http status", t, func() {
		BindHTTPStatusToErrorTag("my_tag", 700)
		ex := NewException("my_tag", "message")
		So(ex.HTTPStatus(), ShouldEqual, 700)
		if ex.HTTPStatus() != 700 {
			t.Fail()
		}
	})

	Convey("should return error 500 when tag is not mapped", t, func() {
		ex := NewException("new_tag", "message")
		So(ex.HTTPStatus(), ShouldEqual, 500)
		So(ex.Error(), ShouldEqual, "message")

	})
}
