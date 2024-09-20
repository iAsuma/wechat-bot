package main

import (
	"context"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/sashabaranov/go-openai"
	"net/smtp"
	"strings"
	"testing"
	"wechatbot/internal/service"
)

func MergeSlice(s1 []string, s2 []string) []string {
	slice := make([]string, len(s1)+len(s2))
	copy(slice, s1)
	copy(slice[len(s1):], s2)
	return slice
}

func SendToMail(user, password, host, subject, date, body, mailtype, replyToAddress string, to []string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var contentType string
	if mailtype == "html" {
		contentType = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		contentType = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	//ccAddress := strings.Join(cc, ";")
	//bccAddress := strings.Join(bcc, ";")
	toAddress := strings.Join(to, ";")

	content := `From: "{username}" <{user}>
To: "{toname}" <{toAddress}>
Subject: {subject}
Date: {date}
{contentType}

{body}`

	content = gstr.ReplaceByMap(content, g.MapStrStr{
		"{user}":        user,
		"{username}":    "WechatBot",
		"{toAddress}":   toAddress,
		"{toname}":      "Asuma",
		"{subject}":     subject,
		"{date}":        date,
		"{contentType}": contentType,
		"{body}":        body,
	})

	//fmt.Println("content", content)

	//msg := []byte("To: " + toAddress + "\r\nFrom: " + user + "\r\nSubject: " + subject + "\r\nDate: " + date + "\r\nReply-To: " + replyToAddress + "\r\n" + contentType + "\r\n\r\n" + body)

	//fmt.Println("msg1", "To: "+toAddress+"\r\nFrom: "+user+"\r\nSubject: "+subject+"\r\nDate: "+date+"\r\nReply-To: "+replyToAddress+"\r\n"+contentType+"\r\n\r\n"+body)

	msg := []byte(content)
	//fmt.Println("msg2", content)
	//sendTo := MergeSlice(to, cc)
	//sendTo = MergeSlice(sendTo, bcc)
	err := smtp.SendMail(host, auth, user, to, msg)
	return err
}

func TestAny(t *testing.T) {
	//user := "770878450@qq.com"
	//password := "ghyqosxpinkbbdbi"
	//host := "smtp.qq.com:587"
	//to := []string{"sqiu_li@163.com"}
	////var cc []string
	////var bcc []string
	subject := "WX-bot系统故障"
	//date := fmt.Sprintf("%s", time.Now().Format(time.RFC1123Z))
	//mailtype := "html"
	//replyToAddress := "251025241@qq.com"
	body := "<html><body><h3>请勿回复520132</h3></body></html>"
	//fmt.Println("send email")
	//err := SendToMail(user, password, host, subject, date, body, mailtype, replyToAddress, to)
	//if err != nil {
	//	fmt.Println("Send mail error!", err)
	//} else {
	//	fmt.Println("Send mail success!")
	//}

	service.NewEmail().Send(gctx.New(), subject, body)
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

func sendEmail(addrSMTP string, from string, to []string, subject string, body string) error {
	auth := smtp.PlainAuth("", "770878450@qq.com", "your_password", "smtp.gmail.com")
	email := "your_email@qq.com" // sender email

	// Create the message
	msg := "From: " + from + "\r\nTo: " + to[0] + "\r\nSubject: " + subject + "\r\n\r\n" + body

	// send email
	err := smtp.SendMail(addrSMTP, auth, email, to, []byte(msg))
	if err != nil {
		return err
	}
	fmt.Printf("Email success sent...\n")

	return nil
}
