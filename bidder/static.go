package bidder

import (
	"github.com/valyala/fasthttp"
)

var (
	staticHandler = fasthttp.FSHandler("./public", 0)
)

func StaticHandle(ctx *fasthttp.RequestCtx) {

	staticHandler(ctx)
}
