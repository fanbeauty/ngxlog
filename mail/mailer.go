package mail

import (
	"bytes"
	"html/template"
	"fmt"
	"net/smtp"
	"log"
	"strings"
)

type Config struct {
	Server   string
	Port     int
	Email    string
	Code     string
	Username string
}

type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

const (
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

//define mail setting
var config = Config{Server: "smtp.qq.com", Port: 587, Email: "1783590642@qq.com", Code: "adoqdlclwatshahb", Username: "梅超凡(Major)"}
//var config = Config{Server: "smtp.gmail.com", Port: 465, Email: "meichaofan0921@gmail.com", Code: "huanhuan0921", Username: "梅超凡(Major)"}

func NewRequest(to []string, subject string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		from:    config.Username,
	}
}

func (r *Request) parseTemplate(fileName string, data interface{}) error {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}
	r.body = buffer.String()
	return nil
}

func (r *Request) sendMail() bool {
	msg := []byte("To: " + strings.Join(r.to, ",") + "\r\nFrom: " + r.from +
		"<" + config.Email + ">\r\nSubject: " + r.subject + "\r\n" + MIME + "\r\n\r\n" + r.body)
	//body := "To: " + r.to[0] + "\r\nSubject: " + r.subject + "\r\n" + MIME + "\r\n" + r.body
	SMTP := fmt.Sprintf("%s:%d", config.Server, config.Port)
	if err := smtp.SendMail(SMTP, smtp.PlainAuth("", config.Email, config.Code, config.Server), config.Email, r.to, msg); err != nil {
		return false
	}
	return true
}

func (r *Request) Send(templateName string, items interface{}) {
	err := r.parseTemplate(templateName, items)
	if err != nil {
		log.Fatal(err)
	}
	if ok := r.sendMail(); ok {
		log.Printf("Email has been sent to %s\n", r.to)
	} else {
		log.Printf("Failed to send the email to %s\n", r.to)
	}
}
