package bidder

import (
	"github.com/valyala/fasthttp"
)

var (
	textDot = []byte(".")
)

func HealthHandle(ctx *fasthttp.RequestCtx) {

	ctx.Write(textDot)
}
