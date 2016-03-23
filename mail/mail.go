// Package mail 邮件相关功能
//

package mail

import (
	"net/smtp"
	"time"
)

// SendMail 发送邮件
//
// body 发送内容
//
// to 接收者列表
//
func SendMail(body string, to []string) error {
	// Set up authentication information.
	f := "13164955841@163.com"
	host := "smtp.163.com"
	auth := smtp.PlainAuth("", f, "plck965xlm", host)
	err := smtp.SendMail(host+":25", auth, f, to, []byte("To: "+to[0]+"\r\n"+
		"From: "+f+"\r\n"+
		"Subject: 行情播报,每日一暴"+time.Now().Format("2006-01-02\r\n")+
		"\r\n"+body))
	return err
}
