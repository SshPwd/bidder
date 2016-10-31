package bidder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"golang_bidder/db"
	"golang_bidder/logs"
	m "golang_bidder/model"
	"golang_bidder/utils"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	urlEventImpression = []byte("/impression")
	urlEventValidate   = []byte("/validate")
	urlEventClick      = []byte("/click")
)

var (
	eventLogger *log.Logger = log.New(os.Stderr, ``, log.Ldate|log.Ltime)
)

func init() {

	file, err := os.OpenFile("event.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Critical: log file \"json.log\" cannot be written\n")
		os.Exit(1)
	}

	eventLogger = log.New(file, ``, log.Ldate|log.Ltime)
}

func EventHandle(path []byte, ctx *fasthttp.RequestCtx) {

	defer func() {

		eventLogger.Println(string(path), ctx.URI().String())
	}()

	// fmt.Println("=URI:", ctx.URI().String())

	switch {
	case bytes.HasPrefix(path, urlEventImpression):
		EventImpression(ctx)

	case bytes.HasPrefix(path, urlEventClick):
		EventClick(ctx)

	case bytes.HasPrefix(path, urlEventValidate):
		EventValidate(ctx, nil, nil)
	}
}

// ==

type ClickData struct {
	Date              int     // da
	CountryId         int     // cid
	DeviceId          int     // did
	AdvertiserId      int     // ad
	CampaignId        int     // ca
	CreativeId        int     // cr
	PublisherId       int     // pu
	SeatId            int     // se
	PlacementCount    int     // +pc
	Test              int     // t
	Position          int     // p
	Frequency         int     // f
	SpendImpression   float64 // spi
	RevenueImpression float64 // rvi
	SpendClick        float64 // spc
	RevenueClick      float64 // rvc
	WinPrice          float64 // wp
	RequestId         []byte  // rid
	CategoryTargeting []byte  // ct
	CreativeUrl       []byte  // r
	UserId            []byte  // uid
	Category          []byte  // cat
	InventoryId       string  // in
	Hash              string  //
}

var (
	requireEventParams = []string{
		"rid", "da", "cid", "did", "ad", "ca", "cr", "pu", "se", "in", "t", "p", "f", "uid", "rs", "pc", "cat",
	}

	hashEventParams = []string{
		"rid", "uid", "da", "cid", "did", "ad", "ca", "cr", "pu", "se", "pc", "t", "p", "f",
	}

	bytesContent = []byte("Content")
)

func eventData(event string, ctx *fasthttp.RequestCtx) (data ClickData, ok bool) {

	args := ctx.QueryArgs()

	value := []byte{}

	for _, name := range requireEventParams {
		if value = args.Peek(name); len(value) == 0 {
			return
		}
	}

	data = ClickData{
		RequestId:         args.Peek("rid"),
		CategoryTargeting: args.Peek("ct"),
		CreativeUrl:       args.Peek("r"),
		UserId:            args.Peek("uid"),
		Category:          args.Peek("cat"),
		Date:              utils.ParseInt(args.Peek("da")),
		CountryId:         utils.ParseInt(args.Peek("cid")),
		DeviceId:          utils.ParseInt(args.Peek("did")),
		AdvertiserId:      utils.ParseInt(args.Peek("ad")),
		CampaignId:        utils.ParseInt(args.Peek("ca")),
		CreativeId:        utils.ParseInt(args.Peek("cr")),
		PublisherId:       utils.ParseInt(args.Peek("pu")),
		SeatId:            utils.ParseInt(args.Peek("se")),
		PlacementCount:    utils.ParseInt(args.Peek("pc")),
		InventoryId:       string(args.Peek("in")),
		Test:              utils.ParseInt(args.Peek("t")),
		Position:          utils.ParseInt(args.Peek("p")),
		Frequency:         utils.ParseInt(args.Peek("f")),
	}

	if event == "impression" {
		if value = args.Peek("wp"); len(value) > 0 {
			data.WinPrice = utils.ToFloat64(string(value))
		}
	}

	decBuf := utils.Decode(args.Peek("rs"))
	defer utils.PutBuffer(decBuf)

	if value = utils.Decrpyt(decBuf.Bytes()); len(value) > 0 {

		if zeroIndex := bytes.IndexByte(value, 0); zeroIndex != -1 {
			value = value[:zeroIndex]
		}

		// eventLogger.Printf("rs: %s", value)

		if list := strings.Split(string(value), "_"); len(list) == 4 {

			data.SpendImpression = utils.ToFloat64(list[0])
			data.RevenueImpression = utils.ToFloat64(list[1])
			data.SpendClick = utils.ToFloat64(list[2])
			data.RevenueClick = utils.ToFloat64(list[3])

			// eventLogger.Println("rs: ", list,
			// 	data.SpendImpression,
			// 	data.RevenueImpression,
			// 	data.SpendClick,
			// 	data.RevenueClick)
		}

	}

	buf := utils.GetBuffer()
	defer utils.PutBuffer(buf)

	buf.WriteString("hash_")
	buf.WriteString(event)

	for _, name := range hashEventParams {
		buf.WriteString("_")
		buf.Write(args.Peek(name))
	}

	data.Hash = buf.String()

	ok = true
	return
}

//

type ActionGUID struct {
	PublisherId  int    `json:"publisher_id"`
	SeatId       int    `json:"seat_id"`
	AdvertiserId int    `json:"advertiser_id"`
	CampaignId   int    `json:"campaign_id"`
	CreativeId   int    `json:"creative_id"`
	CountryId    int    `json:"country_id"`
	DeviceId     int    `json:"device_id"`
	InventoryId  string `json:"inventory_id"`
}

func buildActionGUID(data *ClickData) (bufEnc *bytes.Buffer) {

	bufSrc := utils.GetBuffer()
	enc := json.NewEncoder(bufSrc)

	actionGUID := ActionGUID{
		PublisherId:  data.PublisherId,
		SeatId:       data.SeatId,
		InventoryId:  data.InventoryId,
		AdvertiserId: data.AdvertiserId,
		CampaignId:   data.CampaignId,
		CreativeId:   data.CreativeId,
		CountryId:    data.CountryId,
		DeviceId:     data.DeviceId,
	}

	enc.Encode(&actionGUID)

	bufEnc = utils.Encode(utils.Encrpyt(bufSrc.Bytes()))
	return
}

func parseActionGUID(data []byte) (actionGUID ActionGUID, ok bool) {

	decBuf := utils.Decode(data)
	defer utils.PutBuffer(decBuf)

	if value := utils.Decrpyt(decBuf.Bytes()); len(value) > 0 {

		if zeroIndex := bytes.IndexByte(value, 0); zeroIndex != -1 {
			value = value[:zeroIndex]
		}

		if err := json.Unmarshal(value, &actionGUID); err != nil {
			fmt.Println(err)
			return
		}

		ok = true
	}
	return
}

//

func EventClick(ctx *fasthttp.RequestCtx) {

	if data, ok := eventData("click", ctx); ok {

		if err := db.ExistSeatId(data.SeatId); err == nil {

			actionGUID := buildActionGUID(&data)
			defer utils.PutBuffer(actionGUID)

			cookie := fasthttp.Cookie{}
			cookie.SetKey(fmt.Sprint("__guid_", data.CampaignId))
			cookie.SetValueBytes(actionGUID.Bytes())
			cookie.SetPathBytes(cookieDomain)
			cookie.SetExpire(time.Now().Add(60 * 60 * 24 * 30 * time.Second))
			ctx.Response.Header.SetCookie(&cookie)

			fmt.Printf("data.CategoryTargeting: %s\n", data.CategoryTargeting)
			fmt.Printf("bytesContent: %s\n", bytesContent)

			if len(data.CategoryTargeting) > 0 && !bytes.Equal(data.CategoryTargeting, bytesContent) {

				ctx.PostArgs().Set("ias_score", "1000")
				EventValidate(ctx, &data, actionGUID)
				return
			}

			var url = strings.Replace(ctx.URI().String(), "/event/click", "/event/validate", 1)
			html := strings.Replace(redirectTemplate, "{{{url}}}", strconv.Quote(url), 1)

			ctx.SetStatusCode(302)
			ctx.Response.Header.Set("Content-Type", "text/html; charset=utf-8")
			ctx.WriteString(html)
			return
		}
	}

	StatusNotFound(ctx)
	ctx.WriteString("Error00")
}

// rid: 'request_id',
// da:  'date',
// cid: 'country_id',
// did: 'device_id',
// ad:  'advertiser_id',
// ca:  'campaign_id',
// cr:  'creative_id',
// pu:  'publisher_id',
// se:  'seat_id',
// in:  'inventory_id',
// t:   'test',
// spi: 'spend_impression',
// rvi: 'revenue_impression',
// spc: 'spend_click',
// rvc: 'revenue_click',
// p:   'position',
// f:   'frequency'

func EventValidate(ctx *fasthttp.RequestCtx, data *ClickData, actionGUID *bytes.Buffer) {

	if data == nil {

		evData, ok := eventData("click", ctx)

		if !ok {

			StatusNotFound(ctx)
			ctx.WriteString("Error01")
			return
		}

		data = &evData
	}

	seat, err := db.GetSeatById(data.SeatId)
	if err != nil {

		StatusNotFound(ctx)
		ctx.WriteString("Error02")
		return
	}

	creativeUrl := string(data.CreativeUrl)

	// fmt.Printf("data.CreativeUrl: %s\n", creativeUrl)

	parsedLink, err := url.Parse(creativeUrl)
	if err != nil {
		fmt.Println(err)
		return
	}

	resultLink := url.URL{
		Host:     parsedLink.Host,
		Path:     parsedLink.Path,
		Scheme:   parsedLink.Scheme,
		Fragment: parsedLink.Fragment,
	}

	resultQuery := parsedLink.Query()
	paramValue := ""

	for paramName, paramValues := range resultQuery {

		paramValue = paramValues[0]

		switch {
		case strings.Index(paramValue, "{device_id}") != -1:
			paramValue = strings.Replace(paramValue, "{device_id}", strconv.FormatInt(int64(data.DeviceId), 10), 1)
			resultQuery.Del(paramName)
			resultQuery.Add(paramName, paramValue)

		case strings.Index(paramValue, "{campaign_id}") != -1:
			paramValue = strings.Replace(paramValue, "{campaign_id}", strconv.FormatInt(int64(data.CampaignId), 10), 1)
			resultQuery.Del(paramName)
			resultQuery.Add(paramName, paramValue)

		case strings.Index(paramValue, "{creative_id}") != -1:
			paramValue = strings.Replace(paramValue, "{creative_id}", strconv.FormatInt(int64(data.CreativeId), 10), 1)
			resultQuery.Del(paramName)
			resultQuery.Add(paramName, paramValue)

		case strings.Index(paramValue, "{inventory_id}") != -1:
			paramValue = strings.Replace(paramValue, "{inventory_id}", data.InventoryId, 1)
			resultQuery.Del(paramName)
			resultQuery.Add(paramName, paramValue)

		case strings.Index(paramValue, "{guid}") != -1:

			if actionGUID == nil {
				actionGUID = buildActionGUID(data)
				defer utils.PutBuffer(actionGUID)
			}

			paramValue = strings.Replace(paramValue, "{guid}", actionGUID.String(), 1)
			resultQuery.Del(paramName)
			resultQuery.Add(paramName, paramValue)
		}
	}

	creativeUrl = resultLink.String()

	if len(resultQuery) > 0 {
		creativeUrl += "?"
		creativeUrl += resultQuery.Encode()
	}

	fmt.Println(creativeUrl)

	ctx.SetStatusCode(302)
	ctx.Response.Header.Set("Location", creativeUrl)

	postArgs := ctx.PostArgs()
	iasScore := utils.ParseInt(postArgs.Peek("ias_score"))

	if true || len(data.CategoryTargeting) > 0 && !bytes.Equal(data.CategoryTargeting, bytesContent) {

		clickEvent(data, &seat, iasScore)

	} else {

		db.NotDuplicate(data.Hash, 60*60*6, func() {
			clickEvent(data, &seat, iasScore)
		})
	}

}

func clickEvent(data *ClickData, seat *m.Seat, iasScore int) {

	tm := time.Now()

	date := fmt.Sprintf("%04d-%02d-%02d", tm.Year(), tm.Month(), tm.Day())
	dateTime := fmt.Sprintf("%04d-%02d-%02d %02d:00:00", tm.Year(), tm.Month(), tm.Day(), tm.Hour())
	dateTimeMin := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:00", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute())

	db.Aggregate("inventory_ias",
		db.AggregateKey{
			Date:        date,
			InventoryId: data.InventoryId,
		}, db.AggregateMetric{
			Clicks: 1,
			Score:  iasScore,
		})

	db.Aggregate("publisher",
		db.AggregateKey{
			Datetime:    dateTime,
			PublisherId: data.PublisherId,
			SeatId:      data.SeatId,
			Test:        data.Test,
		}, db.AggregateMetric{
			Requests:      0,
			RequestErrors: 0,
			Bids:          0,
			Impressions:   0,
			Clicks:        1,
			Revenue:       data.RevenueClick,
		})

	db.Aggregate("advertiser",
		db.AggregateKey{
			Date:         date,
			CountryId:    data.CountryId,
			DeviceId:     data.DeviceId,
			AdvertiserId: data.AdvertiserId,
			CampaignId:   data.CampaignId,
			CreativeId:   data.CreativeId,
			PublisherId:  data.PublisherId,
			SeatId:       data.SeatId,
			InventoryId:  data.InventoryId,
			Test:         data.Test,
		}, db.AggregateMetric{
			Bids:        0,
			Impressions: 0,
			Clicks:      1,
			Actions:     0,
			Revenue:     data.RevenueClick,
			Spend:       data.SpendClick,
		})

	if data.Test == 0 {

		db.Aggregate("advertiser_minute",
			db.AggregateKey{
				Datetime:     dateTimeMin,
				AdvertiserId: data.AdvertiserId,
				CampaignId:   data.CampaignId,
			}, db.AggregateMetric{
				Bids:        0,
				Impressions: 0,
				Clicks:      1,
				Spend:       data.SpendClick,
				Revenue:     data.RevenueClick,
			})

		db.Aggregate("revshare",
			db.AggregateKey{
				Datetime:   dateTimeMin,
				SeatId:     data.SeatId,
				CampaignId: data.CampaignId,
			}, db.AggregateMetric{
				Spend:   data.SpendClick,
				Revenue: data.RevenueClick,
			})

		if data.Frequency == 0 {
			data.Frequency = 1
		}

		db.Aggregate("features",
			db.AggregateKey{
				Datetime:  dateTime,
				Position:  data.Position,
				Frequency: data.Frequency,
			}, db.AggregateMetric{
				Impressions:    0,
				Clicks:         1,
				Spend:          data.SpendClick,
				RevenueAuction: data.RevenueClick,
				RevenueWon:     data.RevenueClick,
			})
	}

	logs.Report("click", data.Test, data.SeatId, string(data.RequestId),
		fmt.Sprintf("%d,%d,%.06f,%.06f", data.CampaignId, data.CreativeId, data.SpendClick, data.RevenueClick))

	db.FeaturesSave(
		data.UserId,
		data.Category,
		data.CampaignId,
		data.InventoryId,
		data.UserId,
		data.CountryId,
		data.DeviceId,
		0,
		0,
		data.RevenueClick,
		data.SpendClick,
		"click")
}

//

func EventImpression(ctx *fasthttp.RequestCtx) {

	if data, ok := eventData("impression", ctx); ok && int64(data.Date) > time.Now().Unix()-60000 {

		if seat, err := db.GetSeatById(data.SeatId); err == nil {

			db.NotDuplicate(data.Hash, 60*60*6, func() {
				impressionEvent(&data, &seat)
			})
		}
	}
}

func impressionEvent(data *ClickData, seat *m.Seat) {

	// TODO

	secondPriceAuction := struct {
		original float64
		won      float64
	}{
		original: data.RevenueImpression,
		won:      0,
	}

	data.RevenueImpression = seat.Module.Win(data.WinPrice, data.RevenueImpression)
	secondPriceAuction.won = data.RevenueImpression

	// // We Allow Second Price Auctions For CPM Campaigns
	if data.SpendImpression > 0 {
		data.SpendImpression = data.SpendImpression * (secondPriceAuction.won / secondPriceAuction.original)
	}

	tm := time.Now()

	date := fmt.Sprintf("%04d-%02d-%02d", tm.Year(), tm.Month(), tm.Day())
	dateTime := fmt.Sprintf("%04d-%02d-%02d %02d:00:00", tm.Year(), tm.Month(), tm.Day(), tm.Hour())
	dateTimeMin := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:00", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute())

	// ==

	impressions := 1

	if seat.NoWinNotice == 1 {
		impressions = 0
	}

	db.Aggregate("publisher",
		db.AggregateKey{
			Datetime:    dateTime,
			PublisherId: data.PublisherId,
			SeatId:      data.SeatId,
			Test:        data.Test,
		}, db.AggregateMetric{
			Requests:      0,
			RequestErrors: 0,
			Bids:          0,
			Impressions:   impressions,
			Clicks:        0,
			Revenue:       data.RevenueImpression,
		})

	db.Aggregate("advertiser",
		db.AggregateKey{
			Date:         date,
			CountryId:    data.CountryId,
			DeviceId:     data.DeviceId,
			AdvertiserId: data.AdvertiserId,
			CampaignId:   data.CampaignId,
			CreativeId:   data.CreativeId,
			PublisherId:  data.PublisherId,
			SeatId:       data.SeatId,
			InventoryId:  data.InventoryId,
			Test:         data.Test,
		}, db.AggregateMetric{
			Bids:        0,
			Impressions: impressions,
			Clicks:      0,
			Actions:     0,
			Revenue:     data.RevenueImpression,
			Spend:       data.SpendImpression,
		})

	if data.Test == 0 {

		db.Aggregate("advertiser_minute",
			db.AggregateKey{
				Datetime:     dateTimeMin,
				AdvertiserId: data.AdvertiserId,
				CampaignId:   data.CampaignId,
			}, db.AggregateMetric{
				Bids:        0,
				Clicks:      0,
				Impressions: impressions,
				Spend:       data.SpendImpression,
				Revenue:     data.RevenueImpression,
			})

		db.Aggregate("revshare",
			db.AggregateKey{
				Datetime:   dateTimeMin,
				SeatId:     data.SeatId,
				CampaignId: data.CampaignId,
			}, db.AggregateMetric{
				Spend:   data.SpendImpression,
				Revenue: data.RevenueImpression,
			})

		if data.Frequency == 0 {
			data.Frequency = 1
		}

		db.Aggregate("features",
			db.AggregateKey{
				Datetime:  dateTime,
				Position:  data.Position,
				Frequency: data.Frequency,
			}, db.AggregateMetric{
				Impressions:    impressions,
				Clicks:         0,
				Spend:          data.SpendImpression,
				RevenueAuction: secondPriceAuction.original,
				RevenueWon:     secondPriceAuction.won,
			})
	}

	logs.Report("win", data.Test, data.SeatId, string(data.RequestId),
		fmt.Sprintf("%d,%d,%.06f,%.06f", data.CampaignId, data.CreativeId, data.SpendClick, data.RevenueClick))

	db.FeaturesSave(
		data.UserId,
		data.Category,
		data.CampaignId,
		data.InventoryId,
		data.UserId,
		data.CountryId,
		data.DeviceId,
		secondPriceAuction.original,
		secondPriceAuction.won,
		data.RevenueClick,
		data.SpendClick,
		"impression")
}

// ==

var (
	// FIXME
	redirectTemplate = `<!DOCTYPE html>
<html>
<head>
</head>
<body>

<script>

	var __fired = false;

	var __IntegralASConfig = {onAPIResult: function(ias){
		if(__fired == false){
			__fired = true;
			__redirect(ias.rsa);
		}
	}};

	setTimeout(function(){
		if(__fired == false) __redirect(0);
	},2500);

	function __redirect(score){
		var redirectForm = document.getElementById('redirect'),
			iasInput = document.getElementById('ias_score');
		iasInput.value = score;
		redirectForm.submit();
	}

</script>

<script type="text/javascript" src="//7077.bapi.adsafeprotected.com/bapi?anId=7077&advId=Macro&campId=MommaDot&pubId=Macro3&placementId=Macro4&chanId=Marco5&chanId=dexiM&adsafe_par&ext_impid=Macro7"></script>

<form id="redirect" method="POST" action={{{url}}} style="display:none;">
	<input id="ias_score" type="text" name="ias_score" value="0"></input>
	<input type="submit"></input>
</form>

</body>
</html>`
)
