package logger

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"

	// . "github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/zap"
)

var (
	curPath = ""
)

func TestMain(t *testing.T) {
	Convey("TestMain", t, func() {
		Convey("test exists", func() {
			e := exists("/usr")
			So(e, ShouldBeNil)

			randDir := ""
			for i := 0; i < 10; i++ {
				randDir += "/" + strconv.Itoa(rand.Int()+1000000)
			}
			e = exists(randDir)
			So(e, ShouldBeError)

			// set curPath
			curPath = "/tmp"
		})
	})
}

func TestInitLogger(t *testing.T) {
	Convey("TestInitLogger", t, func() {
		Convey("File not exists", func() {
			randDir := ""
			for i := 0; i < 100; i++ {
				randDir += "/" + strconv.Itoa(rand.Int()+1000000)
			}
			e := InitLogger(randDir, DebugLevel, time.Local)
			So(e, ShouldBeError)
		})

		Convey("File exists", func() {
			e := InitLogger(curPath, DebugLevel, time.Local)
			So(e, ShouldBeNil)
		})
	})
}

func TestGetLogger(t *testing.T) {
	Convey("TestGetLogger", t, func() {
		Convey("write 4 line with debuglevel", func() {
			e := InitLogger(curPath, DebugLevel, time.Local)
			So(e, ShouldBeNil)

			filename := "test.logger_test.go"
			rand.Seed(int64(time.Now().Nanosecond()))
			for i := 0; i < 20; i++ {
				filename += "." + strconv.Itoa(rand.Int()%100)
			}

			fullFilename := filepath.Join(curPath, filename)
			logger := GetLogger(filename)
			var wg sync.WaitGroup
			for i := 0; i < 100; i++ {
				wg.Add(1)
				go func(index int) {
					logger.Debug("msg", zap.Int("test", index))
					logger.Info("msg", zap.Int("test", index))
					logger.Warn("msg", zap.Int("test", index))
					logger.Error("msg", zap.Int("test", index))
					wg.Done()
				}(i)
			}
			wg.Wait()

			e = FlushAndCloseLogger(filename)
			So(e, ShouldBeNil)

			var data []byte
			data, e = ioutil.ReadFile(fullFilename)
			So(e, ShouldBeNil)

			cmd := exec.Command("rm", "-f", fullFilename)
			e = cmd.Run()
			So(e, ShouldBeNil)

			for i := 10; i < 100; i++ {
				So(bytes.Count(data, []byte("\"test\":"+strconv.Itoa(i))), ShouldEqual, 4)
			}
		})

		Convey("write 4 line with WarnLevel", func() {
			e := InitLogger(curPath, WarnLevel, time.Local)
			So(e, ShouldBeNil)

			filename := "test.logger_test.go"
			rand.Seed(int64(time.Now().Nanosecond()))
			for i := 0; i < 20; i++ {
				filename += "." + strconv.Itoa(rand.Int()%100)
			}

			fullFilename := filepath.Join(curPath, filename)
			logger := GetLogger(filename)
			var wg sync.WaitGroup
			for i := 0; i < 100; i++ {
				wg.Add(1)
				go func(index int) {
					logger.Debug("msg", zap.Int("test", index))
					logger.Info("msg", zap.Int("test", index))
					logger.Warn("msg", zap.Int("test", index))
					logger.Error("msg", zap.Int("test", index))
					wg.Done()
				}(i)
			}
			wg.Wait()

			e = FlushAndCloseLogger(filename)
			So(e, ShouldBeNil)

			var data []byte
			data, e = ioutil.ReadFile(fullFilename)
			So(e, ShouldBeNil)

			cmd := exec.Command("rm", "-f", fullFilename)
			e = cmd.Run()
			So(e, ShouldBeNil)

			for i := 10; i < 100; i++ {
				So(bytes.Count(data, []byte("\"test\":"+strconv.Itoa(i))), ShouldEqual, 2)
			}
		})
	})
}

func TestGetSugarLogger(t *testing.T) {
	Convey("TestGetSugarLogger", t, func() {
		Convey("write 4 line", func() {
			e := InitLogger(curPath, DebugLevel, time.Local)
			So(e, ShouldBeNil)

			filename := "test.logger_test.go"
			rand.Seed(int64(time.Now().Nanosecond()))
			for i := 0; i < 20; i++ {
				filename += "." + strconv.Itoa(rand.Int()%100)
			}

			fullFilename := filepath.Join(curPath, filename)
			logger := GetSugarLogger(filename)
			var wg sync.WaitGroup
			for i := 0; i < 100; i++ {
				wg.Add(1)
				go func(index int) {
					logger.Debugw("msg", "test", index)
					logger.Infow("msg", "test", index)
					logger.Warnw("msg", "test", index)
					logger.Errorw("msg", "test", index)
					wg.Done()
				}(i)
			}
			wg.Wait()

			e = FlushAndCloseLogger(filename)
			So(e, ShouldBeNil)

			var data []byte
			data, e = ioutil.ReadFile(fullFilename)
			So(e, ShouldBeNil)

			cmd := exec.Command("rm", "-f", fullFilename)
			e = cmd.Run()
			So(e, ShouldBeNil)

			for i := 10; i < 100; i++ {
				So(bytes.Count(data, []byte("\"test\":"+strconv.Itoa(i))), ShouldEqual, 4)
			}
		})
	})
}
