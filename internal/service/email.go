package service

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"wechatbot/internal/model"
	"wechatbot/internal/service/qmail"
	"wechatbot/internal/wechatbot"
)

type lEmailService struct{}

func NewEmail() *lEmailService {
	return &lEmailService{}
}

func (l *lEmailService) Send(ctx context.Context, subject, body string) error {
	noticeBody := qmail.NoticeContent(body)

	config, err := g.Cfg().GetWithEnv(ctx, "email")
	if config.IsNil() || err != nil {
		return nil
	}
	emailConfig := model.EmailConfig{}
	_ = config.Struct(&emailConfig)

	userName := wechatbot.GetBotNickName()
	if userName == "" {
		userName = "WX机器人"
	}

	email := qmail.NewEmail(emailConfig.Host, emailConfig.Port, emailConfig.From, emailConfig.Password)
	email.From = qmail.Sender{
		Email: emailConfig.From,
		Name:  userName,
	}
	email.To = []qmail.Receiver{
		{
			Email: emailConfig.To,
		},
	}
	email.Msg = qmail.Paper{
		Body:        noticeBody,
		Subject:     subject,
		ContentType: "Content-Type: text/html; charset=UTF-8",
	}

	err = email.SendMail()
	if err != nil {
		g.Log().Info(ctx, "Send mail error!", err)
	} else {
		g.Log().Info(ctx, "Send mail success!")
	}

	return err
}
