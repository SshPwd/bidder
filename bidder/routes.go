package bidder

import (
	"bytes"
	_ "fmt"
	"github.com/valyala/fasthttp"
	"golang_bidder/browscap"
	"golang_bidder/counters"
	"golang_bidder/ip2location"
	_ "golang_bidder/utils"
	"net/http"
)

var (
	urlBid       = []byte("/bid")
	urlSync      = []byte("/sync")
	urlStat      = []byte("/stat")
	urlEvent     = []byte("/event")
	urlHealth    = []byte("/health")
	urlPayload   = []byte("/payload")
	textNotFound = []byte("404 not found")
	textSlash    = []byte("/")

	debugTestIp    = []byte("/test/ip2location")
	debugTestBc    = []byte("/test/browscap")
	debugTestPprof = []byte("/debug")

	Debug int = 1

	Datacenters = ip2location.DefaultDatacenters
	Browsers    = browscap.DefaultBrowscap
	Geoip       = ip2location.DefaultIP2Location
)

func HttpHandle(ctx *fasthttp.RequestCtx) {

	// defer utils.WTimeStart().Stop()

	path := bytes.TrimSuffix(ctx.Path(), textSlash)

	if Debug != 0 {

		switch {
		case bytes.Equal(path, debugTestIp):
			DebugIp2locationHandle(ctx)
			return

		case bytes.Equal(path, debugTestBc):
			DebugBrowscapHandle(ctx)
			return

		case bytes.HasPrefix(path, debugTestPprof):
			DebugPprof(ctx)
			return
		}
	}

	switch {
	case bytes.Equal(path, urlBid):
		BidHandle(ctx)

	case bytes.HasPrefix(path, urlEvent):
		EventHandle(bytes.TrimPrefix(path, urlEvent), ctx)

	case bytes.Equal(path, urlHealth):
		HealthHandle(ctx)

	case bytes.Equal(path, urlSync):
		SyncHandle(ctx)

	case bytes.Equal(path, urlStat):
		StatHandle(ctx)

	default:
		StaticHandle(ctx)
	}

	counters.Request()
}

func StatusNoContent(ctx *fasthttp.RequestCtx, reason string) { // 204

	if reason != "" {
		ctx.Response.Header.Set("X-Reason", reason)
	}

	// fmt.Print(CB, http.StatusNoContent, " ", reason, CN, "\n")

	ctx.SetStatusCode(http.StatusNoContent)
}

func StatusNotFound(ctx *fasthttp.RequestCtx) { // 404

	// fmt.Print(CB, http.StatusNotFound, CN, "\n")

	ctx.SetStatusCode(http.StatusNotFound)
	ctx.Write(textNotFound)
}

func StatusMethodNotAllowed(ctx *fasthttp.RequestCtx) {

	// fmt.Print(CB, http.StatusMethodNotAllowed, CN, "\n")

	ctx.SetStatusCode(http.StatusMethodNotAllowed)
}
