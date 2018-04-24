package base32

import (
	b "encoding/base32"
	"io/ioutil"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDecode(t *testing.T) {
	Convey("TestDecode", t, func() {
		Convey("Equel with encoding/base32", func() {
			data, e := ioutil.ReadFile("base32decode.go")
			So(e, ShouldBeNil)
			source := b.StdEncoding.EncodeToString(data)
			dst1, e1 := b.StdEncoding.DecodeString(source)
			dst2 := Decode(source)

			So(e1, ShouldBeNil)
			So(reflect.DeepEqual(dst1, dst2), ShouldBeTrue)
		})
	})
}
