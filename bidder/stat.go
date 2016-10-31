package bidder

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"golang_bidder/counters"
	"golang_bidder/db"
	"strconv"
)

var (
	BuildDate, BuildName string
)

func StatHandle(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Set("Content-Type", "text/plain; charset=utf-8")

	bidderData := db.GetBidderData()

	ctx.WriteString("bidderData: ")
	ctx.WriteString(bidderData.Date)
	ctx.WriteString("\n")

	ctx.WriteString("BuildDate: ")
	ctx.WriteString(BuildDate)
	ctx.WriteString("\n")

	ctx.WriteString("BuildName: ")
	ctx.WriteString(BuildName)
	ctx.WriteString("\n\n")

	ctx.WriteString("RequestsCount: ")
	ctx.WriteString(strconv.FormatUint(counters.RequestsCount, 10))
	ctx.WriteString("\n")

	var sum int64
	for n := range counters.ProcessedRequests {
		sum += counters.ProcessedRequests[n]
	}
	for n := range counters.ProcessedRequests {
		counters.ProcessedRequests[n] *= 100
		counters.ProcessedRequests[n] /= sum
	}

	ctx.WriteString("ProcessedRequests: ")
	ctx.WriteString(fmt.Sprint(counters.ProcessedRequests))
	ctx.WriteString("\n")
}
