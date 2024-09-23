package handler

import (
	"context"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/gogf/gf/v2/frame/g"
	"log"
	"wechatbot/internal/service"
	"wechatbot/utility/qutil"
)

type BotHandler struct{}

func Bot() *BotHandler {
	return &BotHandler{}
}

func (h *BotHandler) HealCheckCallback(ctx context.Context, resp openwechat.SyncCheckResponse) {
	msg := fmt.Sprintf("HealCheck WARNING RetCode:%s Selector:%s", resp.RetCode, resp.Selector)
	logPrint := false
	if !resp.NorMal() && !resp.HasNewMessage() {
		logPrint = true
	}

	if !resp.Success() {
		if resp.RetCode == "1102" {
			msg += "\n 机器人可能已退出登录"
		}
		_ = service.NewEmail().Send(ctx, "WX-Bot 状态异常", msg)
		logPrint = true
	}

	if logPrint {
		log.Println(msg)
	}
}

func (h *BotHandler) LoginQrcodeUrl(ctx context.Context, uuid string) {
	config, _ := g.Cfg().GetWithEnv(ctx, "app.env")
	appEnv := config.String()

	// 微信登录二维码
	qrcodeUrl := openwechat.GetQrcodeUrl(uuid)

	if appEnv == "dev" {
		println("访问下面网址扫描二维码登录")
		println(qrcodeUrl)

		// browser open the login url
		_ = qutil.Open(qrcodeUrl)
	} else {
		msg := "登录二维码：\n " + qrcodeUrl
		g.Log().Info(ctx, msg)
		_ = service.NewEmail().Send(ctx, "WX-Bot 登录", msg)
	}

}

func (h *BotHandler) LoginCallBack(ctx context.Context, bot *openwechat.Bot, resp openwechat.CheckLoginResponse) {
	config, _ := g.Cfg().GetWithEnv(ctx, "app.env")
	appEnv := config.String()

	self, _ := bot.GetCurrentUser()
	msg := self.NickName + "-登录成功"

	if appEnv == "dev" {
		println(msg)
	} else {
		friends, _ := self.Friends()
		f := friends.SearchByNickName(1, "落").First()
		if f != nil {
			f.SendText(msg + "\n欢迎使用微信机器人，请勿乱用，如涉及违法，后果自负。")
		}

		_ = service.NewEmail().Send(ctx, "WX-Bot 登录", msg)
	}
}
