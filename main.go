package main

import (
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/os/gctx"
	"wechatbot/internal/cmd"
	_ "wechatbot/internal/packed"
)

func main() {
	cmd.OneBot.Run(gctx.GetInitCtx())
}
