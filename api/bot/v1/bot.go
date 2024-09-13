package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type HelloReq struct {
	g.Meta `path:"/hello" tags:"OnText" method:"get" summary:"You first bot api"`
}
type HelloRes struct {
	g.Meta `mime:"text/html" example:"string"`
}
