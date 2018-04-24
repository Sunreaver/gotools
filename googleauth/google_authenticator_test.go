package googleauth

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var key = "MFRGGZDFMZTWQ2LKNNWG23Q="

func TestMakeGoogleAuthenticator(t *testing.T) {
	Convey("TestMakeGoogleAuthenticator", t, func() {
		Convey("BeEqueal", func() {
			data, e := MakeGoogleAuthenticator(key, 1564561)
			So(e, ShouldBeNil)
			So(data, ShouldEqual, "966012")
		})
	})
}
