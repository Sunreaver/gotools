// Package mail 邮件相关功能
//
package mail

import (
	"encoding/json"
	"io/ioutil"
	"net/smtp"
)

// AuthConfig 默认读取
// ./auth.json => AuthConfig
type AuthConfig struct {
	Mail string `json:"username"`
	Pwd  string `json:"password"`
	SMTP string `json:"smtp"`
}

// SendMail 发送邮件
//
// body 发送内容
//
// to 接收者列表
//
func SendMail(authFile, body, from, subject string, to []string) error {

	d, e := ioutil.ReadFile(authFile)
	if e != nil {
		return e
	}
	var cfg AuthConfig
	e = json.Unmarshal(d, &cfg)
	if e != nil {
		return e
	}
	// Set up authentication information.
	auth := smtp.PlainAuth("", cfg.Mail, cfg.Pwd, cfg.SMTP)
	err := smtp.SendMail(cfg.SMTP+":25", auth, cfg.Mail, to,
		[]byte("To: "+to[0]+"\r\n"+
			"From: "+from+"<"+cfg.Mail+">\r\n"+
			"Subject: "+subject+
			"\r\n"+body))
	return err
}
