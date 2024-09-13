package logic

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"wechatbot/internal/service"
)

type lMessageLogic struct{}

var (
	Message = &lMessageLogic{}
)

func (l *lMessageLogic) FriendText(msgCtx *openwechat.MessageContext) {
	msg := msgCtx.Message
	inputText := gstr.Trim(msg.Content)

	if gstr.LenRune(inputText) < 3 {
		msg.ReplyText("少于3个字不会进行AI智能回答")
		return
	}

	go func() {
		replyContent := service.NewOpenAiService().Chat(msgCtx.Context(), inputText)
		msg.ReplyText(replyContent)
	}()
}

func (l *lMessageLogic) GroupText(msgCtx *openwechat.MessageContext) {
	//self, err := msgCtx.Bot().GetCurrentUser()
	msg := msgCtx.Message

	if msg.IsSendBySelf() {
		return
	}

	users, err := gregex.MatchAllString("@([^ ]+)[ ]{1}", msg.Content)
	if err != nil {
		return
	}

	if len(users) != 1 {
		fmt.Println("需要并且只能@1个用户", users, len(users))
		return
	}

	// users[0][1] == self.NickName || users[0][1] == selfInGroup.DisplayName
	if msg.IsAt() {
		AtNikeName := users[0][0]
		inputText := gstr.Replace(msg.Content, AtNikeName, "")
		inputText = gstr.Trim(inputText)

		if gstr.LenRune(inputText) < 3 {
			sendInGroup, _ := msg.SenderInGroup()
			senderName := sendInGroup.DisplayName
			if senderName == "" {
				senderName = sendInGroup.NickName
			}
			replyText := fmt.Sprintf("@%s 少于3个字不会进行AI智能回答", senderName)
			msg.ReplyText(replyText)
			return
		}

		go func() {
			replyContent := service.NewOpenAiService().Chat(msgCtx.Context(), inputText)
			msg.ReplyText(replyContent)
		}()
		return
	}
}
