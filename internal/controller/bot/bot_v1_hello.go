package bot

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"

	"wechatbot/api/bot/v1"
)

func (c *ControllerV1) Hello(ctx context.Context, req *v1.HelloReq) (res *v1.HelloRes, err error) {
	g.RequestFromCtx(ctx).Response.Writeln("OnText! Im a bot!")
	return
}
