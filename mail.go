package tools

import (
	"net/smtp"
	"time"
)

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
