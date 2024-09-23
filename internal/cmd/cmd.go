package cmd

import (
	"context"
	"github.com/eatmoreapple/openwechat"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"wechatbot/internal/handler"
	"wechatbot/internal/service"
	"wechatbot/internal/wechatbot"

	"wechatbot/internal/controller/bot"
)

var (
	Main = gcmd.Command{
		Name:        "main",
		Usage:       "main",
		Brief:       "start http server",
		Description: "start http",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					bot.NewV1(),
				)
			})
			s.Run()
			return nil
		},
	}
	OneBot = gcmd.Command{
		Name:        "one",
		Usage:       "one",
		Brief:       "wechat bot one",
		Description: "start a wechat bot for yourself",
		Arguments:   nil,
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			reloadStorage := openwechat.NewFileHotReloadStorage("storage.json")
			defer reloadStorage.Close()

			bot := openwechat.DefaultBot(openwechat.WithContextOption(ctx), openwechat.Desktop) // 桌面模式
			ctx = bot.Context()

			// 生成登录二维码
			bot.UUIDCallback = func(uuid string) {
				handler.Bot().LoginQrcodeUrl(ctx, uuid)
			}

			// 登录回调
			bot.LoginCallBack = func(resp openwechat.CheckLoginResponse) {
				handler.Bot().LoginCallBack(ctx, bot, resp)
			}

			// 心跳回调函数
			bot.SyncCheckCallback = func(resp openwechat.SyncCheckResponse) {
				handler.Bot().HealCheckCallback(ctx, resp)
			}

			// 登录
			if err = bot.HotLogin(reloadStorage, openwechat.NewRetryLoginOption()); err != nil {
				g.Log().Error(ctx, err)
				return
			}

			// 注册消息回调函数
			bot.MessageHandler = handler.Message().Listen()

			// 设置登录用户到全局
			self, _ := bot.GetCurrentUser()
			wechatbot.SetBotNickName(self.NickName)

			// 阻塞主goroutine, 直到发生异常或者用户主动退出
			err = bot.Block()
			if err != nil {
				_ = service.NewEmail().Send(ctx, "WX-Bot 系统故障", "bot error: "+err.Error())
			}
			return nil
		},
	}
)
