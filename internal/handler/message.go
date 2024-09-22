package handler

import (
	"context"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtimer"
	"time"
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
	dispatcher.OnFriendAdd(h.OnFriendAdd)

	return dispatcher.AsMessageHandler()
}

func (h *MessageHandler) OnText(msgCtx *openwechat.MessageContext) {
	msg := msgCtx.Message
	g.Log().Info(msgCtx.Context(), "OnText", msg.Content)

	send, _ := msg.Sender()
	fmt.Println("send", send.NickName, send.UserName, send.DisplayName, send.ID(), send.Uin, send.RemarkName, send.RemarkPYQuanPin, send.RemarkPYInitial, "#")

	rec, _ := msg.Receiver()
	fmt.Println("rec", rec.NickName, rec.UserName, rec.DisplayName, rec.ID(), rec.Uin, rec.RemarkName, rec.RemarkPYQuanPin, rec.RemarkPYInitial, "#")

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

func (h *MessageHandler) OnFriendAdd(msgCtx *openwechat.MessageContext) {
	msg := msgCtx.Message
	fmt.Println("OnFriendAdd: ", msg.Content)

	gtimer.SetTimeout(msg.Context(), time.Second*5, func(ctx context.Context) {
		msg.Agree("现在可以开始跟我AI对话了")
	})
}
