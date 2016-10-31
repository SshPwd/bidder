package bidder

import (
	"bufio"
	"fmt"
	"github.com/valyala/fasthttp"
	"golang_bidder/browscap"
	"net/http"
	"net/http/httptest"
	"net/http/pprof"
	"strings"
)

func DebugIp2locationHandle(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Set("Content-Type", "text/plain; charset=utf-8")

	remoteIpStr := ctx.RemoteIP().String()

	if paramIp := ctx.QueryArgs().PeekBytes([]byte("ip")); paramIp != nil {
		remoteIpStr = string(paramIp)
	}

	ctx.WriteString(remoteIpStr)
	ctx.WriteString("\n\n")

	info := Geoip.Find(remoteIpStr)

	ctx.WriteString(fmt.Sprintf("Country: %s\nCity: %s\nRegion: %s\nType: %s\nHosting: %v",
		info.Country,
		info.City,
		info.Region,
		info.Type,
		info.IsHosting))
}

func DebugBrowscapHandle(ctx *fasthttp.RequestCtx) {

	userAgent := string(ctx.UserAgent())

	if paramUa := ctx.QueryArgs().PeekBytes([]byte("ua")); paramUa != nil {
		userAgent = string(paramUa)
	}

	ctx.Response.Header.Set("Content-Type", "text/plain; charset=utf-8")

	ctx.WriteString(userAgent)
	ctx.WriteString("\n\n")

	info := Browsers.FindStr(userAgent)

	ctx.WriteString(fmt.Sprintf("Platform: %s\nOS: %s\nCrawler: %v",
		browscap.PlatformName[info.PlatformId],
		browscap.OsName[info.OsId],
		info.Crawler))
}

func DebugPprof(ctx *fasthttp.RequestCtx) {

	debug := false

	if headerDebug := ctx.Request.Header.Peek("X-Debug"); len(headerDebug) > 0 {
		if headerDebug[0] == '1' {
			debug = true
		}
	}

	if !debug {
		StatusNotFound(ctx)
		return
	}

	switch string(ctx.Path()) {
	case "/debug/pprof/":

		req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(fmt.Sprintf("%s", &ctx.Request))))
		if err != nil {
			return
		}

		rw := httptest.NewRecorder()
		pprof.Index(rw, req)

		ctx.Response.Header.Set("Content-Type", "text/html; charset=utf-8")
		ctx.SetStatusCode(http.StatusOK)
		ctx.Write(rw.Body.Bytes())

	case "/debug/pprof/threadcreate",
		"/debug/pprof/goroutine",
		"/debug/pprof/heap",
		"/debug/pprof/block":

		req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(fmt.Sprintf("%s", &ctx.Request))))
		if err != nil {
			return
		}

		rw := httptest.NewRecorder()
		pprof.Index(rw, req)

		ctx.Response.Header.Set("Content-Type", "text/plain; charset=utf-8")
		ctx.SetStatusCode(http.StatusOK)
		ctx.Write(rw.Body.Bytes())

		return

	case "/debug/pprof/profile":

		req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(fmt.Sprintf("%s", &ctx.Request))))
		if err != nil {
			return
		}

		rw := httptest.NewRecorder()
		pprof.Profile(rw, req)

		ctx.Response.Header.Set("Content-Type", "text/plain; charset=utf-8")
		ctx.SetStatusCode(http.StatusOK)
		ctx.Write(rw.Body.Bytes())
	}
}
