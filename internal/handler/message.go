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

func (h *MessageHandler) Listen() openwechat.MessageHandler {
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

func (h *MessageHandler) OnText(msgCtx *openwechat.MessageContext) {
	msg := msgCtx.Message
	g.Log().Info(msgCtx.Context(), "OnText", msg.Content)

	send, _ := msg.Sender()
	fmt.Println("send", send.NickName, send.UserName, send.DisplayName, send.ID(), send.Uin, "#")

	sgr, err := msg.SenderInGroup()
	if err != nil {
		return
	}
	fmt.Println("sgr", sgr.NickName, sgr.UserName, sgr.DisplayName, sgr.ID(), sgr.Uin, "#")

	return
}

func (h *MessageHandler) OnImage(msgCtx *openwechat.MessageContext) {
	msg := msgCtx.Message
	fmt.Println("OnImage: ", msg.Content)
	return
}

func (h *MessageHandler) OnEmoticon(msgCtx *openwechat.MessageContext) {
	msg := msgCtx.Message
	fmt.Println("OnEmoticon: ", msg.RawContent)
	return
}

func (h *MessageHandler) OnVoice(msgCtx *openwechat.MessageContext) {
	msg := msgCtx.Message
	fmt.Println("OnVoice: ", msg.Content)
	return
}

func (h *MessageHandler) OnFriend(msgCtx *openwechat.MessageContext) {
	msg := msgCtx.Message
	fmt.Println("OnFriend: ", msg.Content)
	msg.AsRead()

	if msg.IsText() {
		logic.Message().FriendText(msgCtx)
	}

	return
}

func (h *MessageHandler) OnGroup(msgCtx *openwechat.MessageContext) {
	msg := msgCtx.Message

	if msg.StatusNotify() {
		return
	}
	msg.AsRead()

	if msg.IsText() {
		logic.Message().GroupText(msgCtx)
	}

	return
}
