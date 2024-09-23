package qmail

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"wechatbot/internal/consts"
)

func ErrorContent(msg string) string {
	body := consts.ErrorEmail

	bodyStr := gstr.ReplaceByMap(body, g.MapStrStr{
		"{GEN_TIME}":     gtime.Now().String(),
		"{SOME_MESSAGE}": msg,
	})

	return bodyStr
}
