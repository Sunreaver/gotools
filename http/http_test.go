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
			resp, e := Get("http://www.baidu.com/robots.txt")
			So(e, ShouldBeNil)
			So(string(resp.GetContent()), ShouldContainSubstring, "User-agent: Googlebot")
			So(resp.GetStatusCode(), ShouldEqual, http.StatusOK)
		})
	})
}

func TestPost(t *testing.T) {
	Convey("TestPost", t, func() {
		Convey("Post Baidu", func() {
			resp, e := Post("http://www.baidu.com/robots.txt",
				map[string]string{"Content-Type": "application/json"},
				strings.NewReader("abc"))

			So(e, ShouldBeNil)
			So(string(resp.GetContent()), ShouldContainSubstring, "User-agent: Googlebot")
			So(resp.GetStatusCode(), ShouldEqual, http.StatusOK)
		})
	})
}
