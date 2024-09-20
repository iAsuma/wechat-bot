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
	if !resp.NorMal() && !resp.HasNewMessage() {
		log.Println(msg)
	}

	if !resp.Success() {
		_ = service.NewEmail().Send(context.Background(), "WX-Bot 状态异常", msg)
		log.Println(msg)
	}
}
