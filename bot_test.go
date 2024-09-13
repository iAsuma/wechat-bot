package main

import (
	"context"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/sashabaranov/go-openai"
	"os"
	"strings"
	"testing"
)

func TestAny(t *testing.T) {
	f, err := os.Open("storage.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	fmt.Println(f)
}

func TestMainTest(t *testing.T) {
	reloadStorage := openwechat.NewFileHotReloadStorage("storage.json")
	defer reloadStorage.Close()

	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式

	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 登陆
	//if err := bot.Login(); err != nil {
	//	fmt.Println(err)
	//	return
	//}

	if err := bot.HotLogin(reloadStorage, openwechat.NewRetryLoginOption()); err != nil {
		fmt.Println("login err", err)
		return
	}

	// 获取登陆的用户
	self, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}

	dispatcher := openwechat.NewMessageMatchDispatcher()

	// 只处理消息类型为文本类型的消息
	dispatcher.OnText(func(ctx *openwechat.MessageContext) {
		msg := ctx.Message
		//fmt.Println("my", self.ID(), self.NickName, self.UserName)
		//fmt.Println("msg", msg, gjson.MustEncodeString(msg))
		fmt.Println("OnText: ", msg.Content)
		send, _ := msg.Sender()
		fmt.Println("send: ", send.UserName, send.NickName, send)

		rec, _ := msg.Receiver()
		fmt.Println("rec: ", rec.UserName, rec.NickName, rec)
		//msg.ReplyText("bot aaa")
	})

	dispatcher.OnGroup(func(ctx *openwechat.MessageContext) {
		msg := ctx.Message

		if msg.MsgType == 51 {
			return
		}

		fmt.Println("# Group In #")
		fmt.Println("SendByGroup: ", msg.IsSendByGroup())

		//sg, _ := msg.SenderInGroup()
		//fmt.Println("send group:", sg.UserName, sg.NickName, sg)

		AtNikeName := "@" + self.NickName
		if strings.Contains(msg.Content, AtNikeName) {
			fmt.Println("正在思考")
			inputText := gstr.Replace(msg.Content, AtNikeName, "")
			inputText = gstr.Trim(inputText)
			//replyContent := openAi(ctx.Context(), inputContet)
			//msg.ReplyText(replyContent)
			go replyText(ctx, inputText)
			fmt.Println("我是写在协程后的程序")
		}
	})

	// 注册消息回调函数
	bot.MessageHandler = dispatcher.AsMessageHandler()

	bot.LogoutCallBack = func(bot *openwechat.Bot) {
		err = bot.CrashReason()
		fmt.Println("logout-err", err)
	}

	// 注册消息处理函数
	//bot.MessageHandler = func(msg *openwechat.Message) {
	//	if msg.IsText() && msg.Content == "ping" {
	//		msg.ReplyText("pong")
	//	}
	//}

	// 获取所有的好友
	friends, err := self.Friends()
	fmt.Println(friends, err)

	friends.First()

	// 获取所有的群组
	groups, err := self.Groups()
	fmt.Println(groups, err)

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}

func replyText(ctx *openwechat.MessageContext, inputText string) {
	msg := ctx.Message
	replyContent := openAi(ctx.Context(), inputText)
	fmt.Println("我在协程里...")
	msg.ReplyText(replyContent)
}

func openAi(ctx context.Context, inputText string) string {
	config := openai.DefaultConfig("sk-C9KR1Bh0xuPoYlQAFd57C84c71F242FdA4D8639811A438Db")
	config.BaseURL = "http://127.0.0.1:3000/v1"

	client := openai.NewClientWithConfig(config)

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "你是万能小助手",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: inputText,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return ""
	}

	fmt.Println(resp.Choices[0].Message.Content)

	return resp.Choices[0].Message.Content
}
