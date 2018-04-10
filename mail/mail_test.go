package mail

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/smtp"
	"testing"

	. "github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSendMail(t *testing.T) {
	Convey("test send mail", t, func() {
		Convey("for succ", func() {
			guardError := errors.New("guard_OK")

			// 为ioutil.ReadFile打桩
			Patch(ioutil.ReadFile, func(_ string) ([]byte, error) {
				config := AuthConfig{
					Mail: "x",
					Pwd:  "c",
					SMTP: "d",
				}
				return json.Marshal(&config)
			})

			defer UnpatchAll()

			// 为smtp.SendMail打桩
			Patch(smtp.SendMail, func(_ string, _ smtp.Auth, _ string, _ []string, _ []byte) error {
				return guardError
			})

			result := SendMail("test.txt", "hello", "a@x.com", "hi", []string{"b@x.com"})

			So(result, ShouldEqual, guardError)
		})
	})
}
