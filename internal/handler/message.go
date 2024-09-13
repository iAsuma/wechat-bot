package handler

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/gogf/gf/v2/frame/g"
	"wechatbot/internal/logic"
)

type MessageHandler struct{}

func Message() *MessageHandler {
	return &MessageHandler{}
}

func (c *MessageHandler) OnText(msgCtx *openwechat.MessageContext) {
	msg := msgCtx.Message
	g.Log().Info(msgCtx.Context(), "OnText", msg.Content)

	return
}

func (c *MessageHandler) OnImage(msgCtx *openwechat.MessageContext) {
	msg := msgCtx.Message
	fmt.Println("OnImage: ", msg.Content)
	return
}

func (c *MessageHandler) OnEmoticon(msgCtx *openwechat.MessageContext) {
	msg := msgCtx.Message
	fmt.Println("OnEmoticon: ", msg.RawContent)
	return
}

func (c *MessageHandler) OnVoice(msgCtx *openwechat.MessageContext) {
	msg := msgCtx.Message
	fmt.Println("OnVoice: ", msg.Content)
	return
}

func (c *MessageHandler) OnFriend(msgCtx *openwechat.MessageContext) {
	msg := msgCtx.Message
	fmt.Println("OnFriend: ", msg.Content)
	msg.AsRead()

	if msg.IsText() {
		logic.Message.FriendText(msgCtx)
	}

	return
}

func (c *MessageHandler) OnGroup(msgCtx *openwechat.MessageContext) {
	msg := msgCtx.Message

	if msg.StatusNotify() {
		return
	}
	msg.AsRead()

	if msg.IsText() {
		logic.Message.GroupText(msgCtx)
	}

	return
}
