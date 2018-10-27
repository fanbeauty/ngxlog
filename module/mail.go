package module

import (
	"net/smtp"
	"strings"
	"fmt"
)

func Send(subject string, content string, to []string) {
	//邮箱地址
	userEmail := "1783590642@qq.com"
	smtpPort := ":587"
	emailAuthCode := "adoqdlclwatshahb" //授权码
	smtpHost := "smtp.qq.com"
	auth := smtp.PlainAuth("", userEmail, emailAuthCode, smtpHost)
	nickname := "梅超凡"
	user := userEmail
	contentType := "Content-Type:text/plain;charset=UTF-8"
	body := content
	msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + nickname +
		"<" + user + ">\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	err := smtp.SendMail(smtpHost+smtpPort, auth, user, to, msg)
	if err != nil {
		fmt.Print("send mail error:%v", err)
	}
}
