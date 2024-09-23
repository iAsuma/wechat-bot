package main

import (
	"context"
	"testing"
	"wechatbot/internal/service"
	"wechatbot/internal/service/qmail"
)

func TestAny(t *testing.T) {
	service.NewEmail().Send(context.TODO(), "WX-Bot 系统故障", qmail.ErrorContent("bot error: HEALTH CHECK WARING Retcode 1102 selecte:0 "))
}

func TestMainTest(t *testing.T) {

}
