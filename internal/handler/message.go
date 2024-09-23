package handler

import (
	"context"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtimer"
	"github.com/golang-module/carbon/v2"
	"time"
	"wechatbot/internal/consts"
	"wechatbot/internal/logic"
)

type MessageHandler struct{}

func Message() *MessageHandler {
	return &MessageHandler{}
}

func (h *MessageHandler) Listen() openwechat.MessageHandler {
	dispatcher := openwechat.NewMessageMatchDispatcher()

	// 创建消息处理中心
	dispatcher.RegisterHandler(func(message *openwechat.Message) bool {
		return true
	}, h.OnAny) // 该handler放在第一个位置，保证所有消息都先被处理
	dispatcher.OnText(h.OnText)
	dispatcher.OnImage(h.OnImage)
	dispatcher.OnEmoticon(h.OnEmoticon)
	dispatcher.OnVoice(h.OnVoice)
	dispatcher.OnFriend(h.OnFriend)
	dispatcher.OnGroup(h.OnGroup)
	dispatcher.OnFriendAdd(h.OnFriendAdd)

	return dispatcher.AsMessageHandler()
}

func (h *MessageHandler) OnAny(msgCtx *openwechat.MessageContext) {
	msgCtx.Message.Set(consts.MsgGetStartTimeKey, carbon.Now())
}

func (h *MessageHandler) OnText(msgCtx *openwechat.MessageContext) {
	msg := msgCtx.Message
	g.Log().Info(msgCtx.Context(), "OnText", msg.Content)
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

	gtimer.SetTimeout(msg.Context(), time.Second*8, func(ctx context.Context) {
		msg.Agree("现在可以开始跟我AI对话了")
	})
}
