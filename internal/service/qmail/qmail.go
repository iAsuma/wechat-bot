package qmail

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/text/gstr"
	"net/smtp"
	"strings"
	"time"
)

type Email struct {
	Host     string
	Port     string
	Auth     smtp.Auth
	Password string
	From     Sender
	To       []Receiver
	Msg      Paper
}

type Paper struct {
	Subject     string `json:"subject"`
	Body        string `json:"body"`
	ContentType string `json:"content_type"`
}

type Sender struct {
	Email string
	Name  string
}

type Receiver struct {
	Email string
	Name  string
}

const (
	MailMsgTpl = `From: "{fromName}" <{fromEmail}>
To: {toEmail}
Subject: {subject}
Date: {date}
{contentType}

{body}`
)

func NewEmail(host, port, user, password string) *Email {
	return &Email{
		Host: host,
		Port: port,
		Auth: smtp.PlainAuth("", user, password, host),
	}
}

func New() *Email {
	return &Email{}
}

func (e *Email) SendMail() error {
	var to []string
	for _, t := range e.To {
		to = append(to, t.Email)
	}

	toEmail := strings.Join(to, ";")

	if e.Msg.ContentType == "" {
		e.Msg.ContentType = "Content-Type: text/plain; charset=UTF-8"
	}

	content := gstr.ReplaceByMap(MailMsgTpl, g.MapStrStr{
		"{fromName}":    e.From.Name,
		"{fromEmail}":   e.From.Email,
		"{toEmail}":     toEmail,
		"{subject}":     e.Msg.Subject,
		"{date}":        time.Now().Format(time.RFC1123Z),
		"{contentType}": e.Msg.ContentType,
		"{body}":        e.Msg.Body,
	})
	//fmt.Println("msg", content)
	g.Log().Stdout(false).Info(gctx.GetInitCtx(), "msg", content)
	msg := []byte(content)

	err := smtp.SendMail(e.Host+":"+e.Port, e.Auth, e.From.Email, to, msg)
	return err
}

func (e *Email) Send(subject, body, mailType string) error {

	var contentType string
	if mailType == "html" {
		contentType = "Content-Type: text/" + mailType + "; charset=UTF-8"
	} else {
		contentType = "Content-Type: text/plain; charset=UTF-8"
	}

	e.Msg = Paper{
		Subject:     subject,
		Body:        body,
		ContentType: contentType,
	}

	err := e.SendMail()
	return err
}
