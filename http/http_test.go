package http

import (
	"net/http"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGet(t *testing.T) {
	Convey("TestGet", t, func() {
		Convey("Get Baidu", func() {
			data, code, e := Get("http://www.baidu.com/robots.txt")
			So(string(data), ShouldContainSubstring, "User-agent: Googlebot")
			So(e, ShouldBeNil)
			So(code, ShouldEqual, http.StatusOK)
		})
	})
}

func TestPost(t *testing.T) {
	Convey("TestPost", t, func() {
		Convey("Post Baidu", func() {
			data, code, e := Post("http://www.baidu.com/robots.txt", "application/json", strings.NewReader("abc"))
			So(string(data), ShouldContainSubstring, "User-agent: Googlebot")
			So(e, ShouldBeNil)
			So(code, ShouldEqual, http.StatusOK)
		})
	})
}
