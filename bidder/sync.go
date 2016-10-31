package bidder

import (
	"github.com/valyala/fasthttp"
	"golang_bidder/db"
	"golang_bidder/utils"
	"time"
)

// irtb.io/sync?seat_id=9&user_id=123abc

var (
	cookieTTL    = time.Duration(10 * 365 * 24 * time.Hour)
	cookieName   = []byte("__uid")
	cookieDomain = []byte("/")
	paramUserId  = []byte("user_id")
)

func SyncHandle(ctx *fasthttp.RequestCtx) {

	args := ctx.QueryArgs()

	queryUserId := args.PeekBytes(paramUserId)
	querySeatId := args.PeekBytes(paramSeatId)

	if queryUserId == nil || querySeatId == nil {
		ctx.Response.Header.Set("X-Reason", "params not set")
		return
	}

	seatId := utils.ParseInt(querySeatId)

	if err := db.ExistSeatId(seatId); err != nil {
		ctx.Response.Header.Set("X-Reason", err.Error())
		return
	}

	userAgent := ctx.UserAgent()
	ipStr := ctx.RemoteIP().String()

	guid := ctx.Request.Header.CookieBytes(cookieName)

	if guid != nil {

		db.SaveUserCombination(seatId, queryUserId, userAgent, guid, ipStr)

	} else {

		guid := utils.GUID()

		cookie := fasthttp.Cookie{}
		cookie.SetKeyBytes(cookieName)
		cookie.SetValueBytes(guid[:])
		cookie.SetPathBytes(cookieDomain)
		cookie.SetExpire(time.Now().Add(cookieTTL))

		ctx.Response.Header.SetCookie(&cookie)

		db.SaveUserCombination(seatId, queryUserId, userAgent, guid[:], ipStr)
	}
}
