package handler

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/gogf/gf/v2/frame/g"
	"wechatbot/internal/logic"
)

type messageHandler struct{}

func Message() *messageHandler {
	return &messageHandler{}
}

func (h *messageHandler) Listen() openwechat.MessageHandler {
	dispatcher := openwechat.NewMessageMatchDispatcher()

	// 创建消息处理中心
	dispatcher.OnText(h.OnText)
	dispatcher.OnImage(h.OnImage)
	dispatcher.OnEmoticon(h.OnEmoticon)
	dispatcher.OnVoice(h.OnVoice)
	dispatcher.OnFriend(h.OnFriend)
	dispatcher.OnGroup(h.OnGroup)

	return dispatcher.AsMessageHandler()
}

func (h *messageHandler) OnText(msgCtx *openwechat.MessageContext) {
	msg := msgCtx.Message
	g.Log().Info(msgCtx.Context(), "OnText", msg.Content)

	return
}

func (h *messageHandler) OnImage(msgCtx *openwechat.MessageContext) {
	msg := msgCtx.Message
	fmt.Println("OnImage: ", msg.Content)
	return
}

func (h *messageHandler) OnEmoticon(msgCtx *openwechat.MessageContext) {
	msg := msgCtx.Message
	fmt.Println("OnEmoticon: ", msg.RawContent)
	return
}

func (h *messageHandler) OnVoice(msgCtx *openwechat.MessageContext) {
	msg := msgCtx.Message
	fmt.Println("OnVoice: ", msg.Content)
	return
}

func (h *messageHandler) OnFriend(msgCtx *openwechat.MessageContext) {
	msg := msgCtx.Message
	fmt.Println("OnFriend: ", msg.Content)
	msg.AsRead()

	if msg.IsText() {
		logic.Message.FriendText(msgCtx)
	}

	return
}

func (h *messageHandler) OnGroup(msgCtx *openwechat.MessageContext) {
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
