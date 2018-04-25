package logger

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	. "github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMain(t *testing.T) {
	Convey("TestMain", t, func() {
		Convey("test exists", func() {
			e := exists("/usr")
			So(e, ShouldBeNil)

			randDir := ""
			for i := 0; i < 10; i++ {
				randDir += strconv.Itoa(rand.Int() + 1000000)
			}
			e = exists(randDir)
			So(e, ShouldBeError)
		})
	})

}

func TestToEarlyMorningTimeDuration(t *testing.T) {
	Convey("TestsleepTime", t, func() {
		Convey("10:10:10 ==> 13h49m50s", func() {
			guard := Patch(time.Now, func() time.Time {
				return time.Date(1970, 1, 1, 10, 10, 10, 0, time.Local)
			})

			defer guard.Unpatch()

			du := ToEarlyMorningTimeDuration(time.Now())
			So(du, ShouldEqual, (13*3600+49*60+50)*time.Second)
		})

		Convey("00:00:00 ==> 24h00m00s", func() {
			guard := Patch(time.Now, func() time.Time {
				return time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local)
			})

			defer guard.Unpatch()

			du := ToEarlyMorningTimeDuration(time.Now())
			So(du, ShouldEqual, (24*3600+0*60+0)*time.Second)
		})

		Convey("00:00:01 ==> 23h59m59s", func() {
			guard := Patch(time.Now, func() time.Time {
				return time.Date(1970, 1, 1, 0, 0, 1, 0, time.Local)
			})

			defer guard.Unpatch()

			du := ToEarlyMorningTimeDuration(time.Now())
			So(du, ShouldEqual, (23*3600+59*60+59)*time.Second)
		})
	})
}
