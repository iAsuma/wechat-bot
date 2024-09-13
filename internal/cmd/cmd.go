package cmd

import (
	"context"
	"github.com/eatmoreapple/openwechat"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"wechatbot/internal/handler"

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

			bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式

			if err = bot.HotLogin(reloadStorage, openwechat.NewRetryLoginOption()); err != nil {
				g.Log().Error(ctx, err)
				return
			}

			// 创建消息处理中心
			dispatcher := openwechat.NewMessageMatchDispatcher()
			dispatcher.OnText(handler.Message().OnText)
			//dispatcher.OnImage(handler.Message().OnImage)
			//dispatcher.OnEmoticon(handler.Message().OnEmoticon)
			dispatcher.OnVoice(handler.Message().OnVoice)
			dispatcher.OnFriend(handler.Message().OnFriend)
			dispatcher.OnGroup(handler.Message().OnGroup)

			// 注册消息回调函数
			bot.MessageHandler = dispatcher.AsMessageHandler()

			// 阻塞主goroutine, 直到发生异常或者用户主动退出
			bot.Block()
			return nil
		},
	}
)
