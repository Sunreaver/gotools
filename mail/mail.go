// Package mail 邮件相关功能
//

package mail

import (
	"encoding/json"
	"io/ioutil"
	"net/smtp"
	"time"
)

// AuthConfig 默认读取
// ./auth.json => AuthConfig
type AuthConfig struct {
	Mail string `json:"username"`
	Pwd  string `json:"password"`
	Smtp string `json:"smtp"`
}

// SendMail 发送邮件
//
// body 发送内容
//
// to 接收者列表
//
func SendMail(body string, to []string) error {

	d, e := ioutil.ReadFile("auth.json")
	if e != nil {
		return e
	}
	var cfg AuthConfig
	e = json.Unmarshal(d, &cfg)
	if e != nil {
		return e
	}
	// Set up authentication information.
	auth := smtp.PlainAuth("", cfg.Mail, cfg.Pwd, cfg.Smtp)
	err := smtp.SendMail(cfg.Smtp+":25", auth, cfg.Mail, to,
		[]byte("To: "+to[0]+"\r\n"+
			"From: "+"每日一报<"+cfg.Mail+">\r\n"+
			"Subject: 行情播报,每日一暴"+time.Now().Format("2006-01-02\r\n")+
			"\r\n"+body))
	return err
}
