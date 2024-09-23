package logic

import (
	"context"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtimer"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/golang-module/carbon/v2"
	"github.com/sashabaranov/go-openai"
	"time"
	"wechatbot/internal/consts"
	"wechatbot/internal/model"
	"wechatbot/internal/service"
)

type lChatLogic struct{}

func Chat() *lChatLogic {
	return &lChatLogic{}
}

// Reply 延迟回复，防封号（可能起作用）
func (l *lChatLogic) Reply(msgCtx *openwechat.MessageContext, msgContent string) {
	msgGetStartTime, exist := msgCtx.Message.Get(consts.MsgGetStartTimeKey)
	var timeCan time.Duration = 3

	if exist {
		msgLen := gstr.LenRune(msgContent)
		diffSeconds := msgGetStartTime.(carbon.Carbon).DiffInSeconds()
		var maxT int64 = 1

		if msgLen < 10 && diffSeconds < 3 {
			maxT = int64(grand.N(3, 5))
		} else if msgLen >= 10 && msgLen < 100 && diffSeconds < 5 {
			maxT = int64(grand.N(6, 10))
		} else if msgLen >= 100 && diffSeconds < 10 {
			maxT = int64(grand.N(10, 15))
		} else {
			maxT = 1
		}

		timeCan = time.Duration(maxT - diffSeconds)
	}

	gtimer.SetTimeout(msgCtx.Context(), time.Second*timeCan, func(ctx context.Context) {
		msgCtx.Message.ReplyText(msgContent)
	})
}

// AiReply AIh回复
func (l *lChatLogic) AiReply(msgCtx *openwechat.MessageContext, msgContent ...string) {
	var chatData []model.ChatMessage

	if len(msgContent) == 0 {
		msgContent = []string{msgCtx.Message.Content}
	}

	chatData = l.GetStorageData(msgCtx)
	chatData = append(chatData, model.ChatMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: msgContent[0],
	})

	replyContent := service.NewOpenAi().Chat(msgCtx.Context(), chatData)
	if replyContent == "" {
		return
	}

	l.Reply(msgCtx, replyContent)

	chatData = append(chatData, model.ChatMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: replyContent,
	})

	l.SetStorageData(msgCtx, chatData)
}

func (l *lChatLogic) GetStorageData(msgCtx *openwechat.MessageContext) []model.ChatMessage {
	var chatData []model.ChatMessage
	var (
		ctx     = msgCtx.Context()
		keyName = l.GetStorageKey(msgCtx)
	)

	data, err := g.Redis().Do(ctx, "get", keyName)
	if err != nil {
		g.Log().Info(ctx, "redis get data error", err)
		return nil
	}

	if data.IsEmpty() {
		return nil
	}

	err = data.Structs(&chatData)
	if err != nil {
		return nil
	}

	return chatData
}

func (l *lChatLogic) SetStorageData(msgCtx *openwechat.MessageContext, chatData []model.ChatMessage) (err error) {
	ctx := msgCtx.Context()
	keyName := l.GetStorageKey(msgCtx)

	_, err = g.Redis().Do(ctx, "setex", keyName, time.Hour.Seconds()*6, gjson.MustEncodeString(chatData))

	if err != nil {
		g.Log().Info(ctx, "redis get data error", err)
		return err
	}

	return nil
}

func (l *lChatLogic) ClearStorageData(msgCtx *openwechat.MessageContext) (err error) {
	ctx := msgCtx.Context()

	keyName := l.GetStorageKey(msgCtx)
	_, err = g.Redis().Do(ctx, "del", keyName)
	return
}

func (l *lChatLogic) HasHistoryData(msgCtx *openwechat.MessageContext) bool {
	ctx := msgCtx.Context()

	keyName := l.GetStorageKey(msgCtx)
	has, err := g.Redis().Do(ctx, "EXISTS", keyName)
	if err != nil {
		return false
	}

	return has.Bool()
}

func (l *lChatLogic) GetStorageKey(msgCtx *openwechat.MessageContext) string {
	msg := msgCtx.Message
	sender, _ := msg.Sender()

	id := sender.NickName

	if msg.IsSendByGroup() {
		fromUser, _ := msg.SenderInGroup()
		senderName := sender.NickName
		if senderName == "" {
			senderName = sender.UserName
		}
		id = fmt.Sprintf("%s_%s", senderName, fromUser.NickName)
	}

	keyName := gmd5.MustEncryptString(id)
	return keyName
}
