package logic

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

type lMessageLogic struct{}

func Message() *lMessageLogic {
	return &lMessageLogic{}
}

func (l *lMessageLogic) FriendText(msgCtx *openwechat.MessageContext) {
	msg := msgCtx.Message
	inputText := gstr.Trim(msg.Content)

	if gstr.LenRune(inputText) < 3 {
		Chat().Reply(msgCtx, "少于3个字不会进行AI智能回答")
		return
	}

	go Chat().AiReply(msgCtx, inputText)
	return
}

func (l *lMessageLogic) GroupText(msgCtx *openwechat.MessageContext) {
	msg := msgCtx.Message

	if msg.IsSendBySelf() {
		return
	}

	atUsers, err := gregex.MatchAllString("@([^\u2005]+)[\u2005]{1}", msg.Content)
	if err != nil {
		return
	}

	if len(atUsers) > 1 {
		fmt.Println("@了 1 个以上人")
		return
	}

	// users[0][1] == self.NickName || users[0][1] == selfInGroup.DisplayName
	if msg.IsAt() {
		AtNikeName := atUsers[0][0]
		inputText := gstr.Replace(msg.Content, AtNikeName, "")
		inputText = gstr.Trim(inputText)
		msg.Content = inputText // 替换内容

		if gstr.LenRune(inputText) < 3 {
			sendInGroup, _ := msg.SenderInGroup()
			senderName := sendInGroup.DisplayName
			if senderName == "" {
				senderName = sendInGroup.NickName
			}
			replyText := fmt.Sprintf("@%s 少于3个字不会进行AI智能回答", senderName)
			Chat().Reply(msgCtx, replyText)
			return
		}

		go func() {
			Chat().ClearStorageData(msgCtx)
			Chat().AiReply(msgCtx, inputText)
		}()
	} else if len(atUsers) == 0 && Chat().HasHistoryData(msgCtx) {
		go Chat().AiReply(msgCtx, gstr.Trim(msg.Content))
	}

	return
}
