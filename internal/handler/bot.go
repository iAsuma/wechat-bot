package handler

import (
	"context"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"log"
	"wechatbot/internal/service"
)

type BotHandler struct{}

func Bot() *BotHandler {
	return &BotHandler{}
}

func (h *BotHandler) HealCheckCallback(resp openwechat.SyncCheckResponse) {
	msg := fmt.Sprintf("HealCheck WARNING RetCode:%s Selector:%s", resp.RetCode, resp.Selector)
	logPrint := false
	if !resp.NorMal() && !resp.HasNewMessage() {
		logPrint = true
	}

	if !resp.Success() {
		if resp.RetCode == "1102" {
			msg += "\n 机器人可能已退出登录"
		}
		_ = service.NewEmail().Send(context.Background(), "WX-Bot 状态异常", msg)
		logPrint = true
	}

	if logPrint {
		log.Println(msg)
	}
}
