package bidder

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"golang_bidder/browscap"
	"golang_bidder/config"
	"golang_bidder/counters"
	"golang_bidder/db"
	"golang_bidder/logs"
	m "golang_bidder/model"
	"golang_bidder/utils"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	paramTest      = []byte("test")
	paramSeatId    = []byte("seat_id")
	paramSecretKey = []byte("secret_key")
)

//

var logger *log.Logger = log.New(os.Stderr, ``, log.Ldate|log.Ltime)

func init() {

	file, err := os.OpenFile("json.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Critical: log file \"json.log\" cannot be written\n")
		os.Exit(1)
	}

	logger = log.New(file, ``, log.Ldate|log.Ltime)
}

// ctx.Response.Header.Set("Content-Type", "application/x-javascript; charset=utf-8")
// ctx.Response.Header.Set("Content-Type", "text/plain; charset=utf-8")

func BidHandle(ctx *fasthttp.RequestCtx) {

	// if !ctx.IsPost() {
	// 	StatusMethodNotAllowed(ctx)
	// 	return
	// }

	// defer func() {
	// 	logger.Printf("<< %s; %s", ctx.URI().FullURI(), ctx.PostBody())
	// 	logger.Printf(">> [%d] %s", ctx.Response.StatusCode(), ctx.Response.Body())
	// }()

	nsec := time.Now().UnixNano()
	defer counters.ProcessingTime(time.Now().UnixNano() - nsec)

	args := ctx.QueryArgs()

	querySeatId := args.PeekBytes(paramSeatId)
	querySecretKey := args.PeekBytes(paramSecretKey)

	if querySeatId == nil || querySecretKey == nil {
		StatusNoContent(ctx, "params not set")
		return
	}

	seatId := utils.ParseInt(querySeatId)

	seat, err := db.GetSeatById2(seatId, querySecretKey)
	if err != nil {
		BidRequestParseError(seatId, utils.ParseInt(args.PeekBytes(paramTest)))
		StatusNoContent(ctx, err.Error())
		return
	}

	request := m.Request{
		Ctx:         ctx,
		Seat:        seat,
		Geoip:       &Geoip,
		Browsers:    &Browsers,
		Datacenters: &Datacenters,
		RateType:    "CPM",
	}

	// fmt.Printf(CG+"Seat: %d, Module: %s\n"+CN, request.Seat.SeatId, request.Seat.Module.Name())

	err = request.Seat.Module.Read(&request)
	if err != nil {
		StatusNoContent(ctx, err.Error())
		return
	}

	request.InventoryData = db.GetInventoryData(&request)
	// fmt.Printf("InventoryData: %#v\n", request.InventoryData)

	request.Campaigns = GetCampaigns(&request)
	// fmt.Printf("Campaigns: %#v\n", len(request.Campaigns))

	if request.Campaigns != nil {

		ctx.SetStatusCode(http.StatusOK)

		BuildResponse(&request)

		request.Seat.Module.Write(&request)

		BidSent(&request)

	} else {

		CampaignsGetError(&request)
		StatusNoContent(ctx, "noСampaigns")
	}
}

// ==

func CampaignsGetError(request *m.Request) {

	// app.globals.fileLog.log('req',bidRequest.bidRequest.test+"\t"+bidRequest.seat.sid+"\t"+bidRequest.bidRequest.id);

	logs.Report("req", request.BidRequest.Test, request.Seat.SeatId, string(request.BidRequest.Id), "")

	tm := time.Now()
	dateTime := fmt.Sprintf("%04d-%02d-%02d %02d:00:00", tm.Year(), tm.Month(), tm.Day(), tm.Hour())

	db.Aggregate("publisher",
		db.AggregateKey{
			Datetime:    dateTime,
			PublisherId: request.Seat.UserId,
			SeatId:      request.Seat.SeatId,
			Test:        request.BidRequest.Test,
		}, db.AggregateMetric{
			Requests:      1,
			RequestErrors: 0,
			Bids:          0,
			Impressions:   0,
			Clicks:        0,
			Revenue:       0,
		})

	if request.BidRequest.Test == 0 && request.Targeting.Inventory != "" && request.Targeting.PlatformId != 0 {

		var phone, tablet, desktop int

		switch request.Targeting.PlatformId {
		case browscap.DESKTOP:
			desktop = request.PlacementCount
		case browscap.TABLET:
			tablet = request.PlacementCount
		case browscap.PHONE:
			phone = request.PlacementCount
		}

		date := fmt.Sprintf("%04d-%02d-%02d", tm.Year(), tm.Month(), tm.Day())

		db.Aggregate("inventory",
			db.AggregateKey{
				Date:        date,
				InventoryId: request.Targeting.Inventory,
			}, db.AggregateMetric{
				PhoneImpressions:   phone,
				TabletImpressions:  tablet,
				DesktopImpressions: desktop,
				Floor:              request.PlacementFloors / 1000.0,
			})
	}
}

func BidRequestParseError(seatId, test int) {

	// if(bidRequest.bidRequest && bidRequest.bidRequest.id)
	//     app.globals.fileLog.log('req',test+"\t"+bidRequest.seat.sid+"\t"+bidRequest.bidRequest.id);
	// else
	//     app.globals.fileLog.log('req',test+"\t"+bidRequest.seat.sid+"\t");

	logs.Report("req", test, seatId, "", "")

	tm := time.Now()
	dateTime := fmt.Sprintf("%04d-%02d-%02d %02d:00:00", tm.Year(), tm.Month(), tm.Day(), tm.Hour())

	db.Aggregate("publisher",
		db.AggregateKey{
			Datetime:    dateTime,
			PublisherId: 0,
			SeatId:      seatId,
			Test:        test,
		}, db.AggregateMetric{
			Requests:      1,
			RequestErrors: 1,
			Bids:          0,
			Impressions:   0,
			Clicks:        0,
			Revenue:       0,
		})
}

func BidSent(request *m.Request) {

	buf := utils.GetBuffer()
	defer utils.PutBuffer(buf)

	for i := range request.Ads0 {

		buf.WriteString(request.Ads1[i].ImpId)
		buf.WriteString(",")
		buf.WriteString(strconv.FormatInt(int64(request.Ads0[i].CampaignId), 10))
		buf.WriteString(",")
		buf.WriteString(strconv.FormatInt(int64(request.Ads0[i].CreativeId), 10))
		buf.WriteString(",")
		buf.WriteString(strconv.FormatFloat(request.Ads0[i].AdvertiserSpendImpression, 'f', 6, 64))
		buf.WriteString(",")
		buf.WriteString(strconv.FormatFloat(request.Ads0[i].PublisherRevenueImpression, 'f', 6, 64))
		buf.WriteString("|")
	}

	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}

	logs.Report("req", request.BidRequest.Test, request.Seat.SeatId, string(request.BidRequest.Id), buf.String())

	//

	tm := time.Now()

	date := fmt.Sprintf("%04d-%02d-%02d", tm.Year(), tm.Month(), tm.Day())
	dateTime := fmt.Sprintf("%04d-%02d-%02d %02d:00:00", tm.Year(), tm.Month(), tm.Day(), tm.Hour())
	dateTimeMin := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:00", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute())

	adsLength := len(request.Ads1)
	impressions := 0

	if request.Seat.NoWinNotice == 1 {
		impressions = 1
	}

	db.Aggregate("publisher",
		db.AggregateKey{
			Datetime:    dateTime,
			PublisherId: request.Seat.UserId,
			SeatId:      request.Seat.SeatId,
			Test:        request.BidRequest.Test,
		}, db.AggregateMetric{
			Requests:      1,
			RequestErrors: 0,
			Bids:          adsLength,
			Impressions:   impressions,
			Clicks:        0,
			Revenue:       0,
		})

	for i := range request.Ads0 {

		db.BidSave(request.Ads0[i].CampaignId, request.Targeting.Inventory)

		db.Aggregate("advertiser",
			db.AggregateKey{
				Date:         date,
				CountryId:    request.Ads0[i].CountryId,
				DeviceId:     request.Ads0[i].DeviceId,
				AdvertiserId: request.Ads0[i].AdvertiserId,
				CampaignId:   request.Ads0[i].CampaignId,
				CreativeId:   request.Ads0[i].CreativeId,
				PublisherId:  request.Seat.UserId,
				SeatId:       request.Seat.SeatId,
				InventoryId:  request.Targeting.Inventory,
				Test:         request.BidRequest.Test,
			}, db.AggregateMetric{
				Bids:        1,
				Impressions: impressions,
				Clicks:      0,
				Actions:     0,
				Revenue:     0,
				Spend:       0,
			})

		if request.BidRequest.Test == 0 {

			db.Aggregate("advertiser_minute",
				db.AggregateKey{
					Datetime:     dateTimeMin,
					AdvertiserId: request.Ads0[i].AdvertiserId,
					CampaignId:   request.Ads0[i].CampaignId,
				}, db.AggregateMetric{
					Bids:        1,
					Impressions: impressions,
					Clicks:      0,
					Revenue:     0,
					Spend:       0,
				})
		}
	}

	if request.BidRequest.Test == 0 && request.Targeting.Inventory != "" && request.Targeting.PlatformId != 0 {

		var phone, tablet, desktop int

		switch request.Targeting.PlatformId {
		case browscap.DESKTOP:
			desktop = request.PlacementCount
		case browscap.TABLET:
			tablet = request.PlacementCount
		case browscap.PHONE:
			phone = request.PlacementCount
		}

		db.Aggregate("inventory",
			db.AggregateKey{
				Date:        date,
				InventoryId: request.Targeting.Inventory,
			}, db.AggregateMetric{
				PhoneImpressions:   phone,
				TabletImpressions:  tablet,
				DesktopImpressions: desktop,
				Floor:              request.PlacementFloors / 1000.0,
			})
	}
}

// ==

func GetCampaigns(request *m.Request) []m.Campaign {

	bidderData := db.GetBidderData()

	campaigns := FilterCampaigns(request, bidderData)

	// fmt.Printf("Placement1: %#v\n", request.Placement)
	// fmt.Println("Campaigns1:", len(campaigns))

	// for n := range campaigns {
	// 	fmt.Println("FilterCampaigns1.1:  ", campaigns[n].CampaignId, campaigns[n].PubRate)
	// }

	campaigns = OptimizeCampaigns(request, bidderData, campaigns)

	// fmt.Println("Campaigns2:", len(campaigns))

	// for n := range campaigns {
	// 	fmt.Println("OptimizeCampaigns2.1:", campaigns[n].CampaignId, campaigns[n].PubRate)
	// }

	// ==

	type PlacementFloors struct {
		adv float64
		pub float64
	}

	placementFloors := map[int]PlacementFloors{}
	biddablePlacements := 0

	for i := range campaigns {

		found := false

		for x := range request.Placement {

			if request.Placement[x].Count > len(request.Placement[x].Campaigns) {

				found = true

				pf, ok := placementFloors[x]
				if !ok {
					pf.pub = request.Placement[x].Floor
				}

				pf.adv += campaigns[i].PubRate

				placementFloors[x] = pf

				request.Placement[x].Campaigns = append(request.Placement[x].Campaigns, campaigns[i])
				biddablePlacements++
			}
		}

		for index, pf := range placementFloors {

			if pf.adv < pf.pub {

				if request.Seat.IgnoreFloor == 0 {

					biddablePlacements -= len(request.Placement[index].Campaigns)
					request.Placement[index].Campaigns = nil
				}
			}
		}

		if biddablePlacements > 5 || !found {
			break
		}
	}

	// for n := range request.Placement {
	// 	fmt.Printf("Placement2: %d\n", len(request.Placement[n].Campaigns))
	// }

	// ==

	if biddablePlacements > 0 {

		// success
		return campaigns

	}

	// err
	return nil
}

// ==

func roundMoney(usd float64) float64 {
	return math.Floor(usd*1000000) / 1000000
}

func buildEventLinks(request *m.Request, ad *m.Bid, _ad *m.Ad0, campaign *m.Campaign, creative *m.Creative) (impression, click string) {

	conf := config.Get()

	params := url.Values{}

	params.Add("rid", request.BidRequest.Id)
	params.Add("da", strconv.FormatInt(time.Now().Unix(), 10))
	params.Add("cid", strconv.FormatInt(int64(_ad.CountryId), 10))
	params.Add("did", strconv.FormatInt(int64(_ad.DeviceId), 10))

	params.Add("ad", strconv.FormatInt(int64(_ad.AdvertiserId), 10))
	params.Add("ca", strconv.FormatInt(int64(_ad.CampaignId), 10))
	params.Add("cr", strconv.FormatInt(int64(_ad.CreativeId), 10))
	params.Add("pu", strconv.FormatInt(int64(_ad.PublisherId), 10))

	params.Add("se", strconv.FormatInt(int64(_ad.SeatId), 10))
	params.Add("in", _ad.InventoryId)
	params.Add("t", strconv.FormatInt(int64(request.BidRequest.Test), 10))
	params.Add("p", ad.Id)

	params.Add("f", strconv.FormatInt(int64(campaign.ImpFreq), 10))
	params.Add("uid", request.GlobalUser)
	params.Add("pc", strconv.FormatInt(int64(request.PlacementCount), 10))
	params.Add("cat", campaign.Category)

	params.Add("ct", campaign.CategoryTargeting)

	bufSrc := utils.GetBuffer()
	bufSrc.WriteString(strconv.FormatFloat(_ad.AdvertiserSpendImpression, 'f', 6, 64))
	bufSrc.WriteString("_")
	bufSrc.WriteString(strconv.FormatFloat(_ad.PublisherRevenueImpression, 'f', 6, 64))
	bufSrc.WriteString("_")
	bufSrc.WriteString(strconv.FormatFloat(_ad.AdvertiserSpendClick, 'f', 6, 64))
	bufSrc.WriteString("_")
	bufSrc.WriteString(strconv.FormatFloat(_ad.PublisherRevenueClick, 'f', 6, 64))

	bufEnc := utils.Encode(utils.Encrpyt(bufSrc.Bytes()))

	params.Add("rs", bufEnc.String())

	utils.PutBuffer(bufEnc)
	utils.PutBuffer(bufSrc)

	// ==

	impression = conf.Event.BaseUrl + "/impression?" + params.Encode() + "&wp=${AUCTION_PRICE}"

	// ==

	params.Add("r", buildCreativeUrl(campaign, creative))
	click = conf.Event.BaseUrl + "/click?" + params.Encode()

	return
}

// ==

// var (
// 	urlReplacer = strings.NewReplacer("%7B", "{", "%7D", "}")
// )

func buildCreativeUrl(campaign *m.Campaign, creative *m.Creative) string {

	if creative.FinalLink != "" {
		return creative.FinalLink
	}

	link := strings.TrimSpace(creative.Link)

	parsedLink, err := url.Parse(link)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	appendQuery := strings.TrimSpace(campaign.UrlParameters)
	if appendQuery == "" {
		creative.FinalLink = link
		return link
	}

	parsedAppendQuery, err := url.ParseQuery(appendQuery)
	if err != nil {
		fmt.Println(err.Error())
		creative.FinalLink = link
		return link
	}

	resultLink := url.URL{
		Host:     parsedLink.Host,
		Path:     parsedLink.Path,
		Scheme:   parsedLink.Scheme,
		Fragment: parsedLink.Fragment,
	}

	resultQuery := parsedLink.Query()

	for paramName, paramValues := range parsedAppendQuery {
		for _, value := range paramValues {
			resultQuery.Add(paramName, value)
		}
	}

	creative.FinalLink = resultLink.String() + "?" + resultQuery.Encode()

	return creative.FinalLink
}

// ==

func shortenString(str string, limit int, sufix string) string {

	if limit > 0 {

		if text := []rune(str); len(text) > limit {

			return string(text[:limit-len(sufix)]) + sufix // FIXME
		}
	}
	return str
}

func buildImageObject(image, brandLogo string, assetsImage m.AssetsImage) m.NativeAssets {

	nativeAssets := m.NativeAssets{
		Id:  assetsImage.Id,
		Img: new(m.NativeImg),
	}

	conf := config.Get()

	if assetsImage.Type == 2 {
		nativeAssets.Img.Url = conf.Cdn.ImgBaseUrl + brandLogo
	} else {
		nativeAssets.Img.Url = conf.Cdn.ImgBaseUrl + image
	}

	query := make([]string, 0, 6)
	query = append(query, "auto=enhance")
	query = append(query, "vib=25")

	crop := false

	if assetsImage.Width > 0 {
		crop = true
		query = append(query, "w="+strconv.FormatInt(int64(assetsImage.Width), 10))
		nativeAssets.Img.Width = assetsImage.Width

	}

	if assetsImage.Heigth > 0 {
		crop = true
		query = append(query, "h="+strconv.FormatInt(int64(assetsImage.Heigth), 10))
		nativeAssets.Img.Height = assetsImage.Heigth
	}

	if crop {
		query = append(query, "crop=faces")
		query = append(query, "fit=crop")
	}

	nativeAssets.Img.Url += "?" + strings.Join(query, "&")

	return nativeAssets
}

func buildAdm(request *m.Request, placement *m.Placement, campaign *m.Campaign, creative *m.Creative, _ad *m.Ad0, click string) string {

	native := m.NativeRequest{}

	native.Version = "1"
	native.Link.Url = click

	// == title
	title := creative.Title

	if placement.HasDescription || placement.Assets.Title.Len <= 30 {
		title = creative.ShortTitle
	}

	native.Assets = append(native.Assets, m.NativeAssets{
		Id:    placement.Assets.Title.Id,
		Title: &m.NativeTitle{Text: shortenString(title, placement.Assets.Title.Len, "…")},
	})

	// == image
	native.Assets = append(native.Assets, buildImageObject(creative.Image, campaign.BrandLogo, placement.Assets.Image))

	// == sponsored
	if placement.Assets.Sponsored.Id != 0 {
		native.Assets = append(native.Assets, m.NativeAssets{
			Id:   placement.Assets.Sponsored.Id,
			Data: &m.NativeData{Value: shortenString(campaign.BrandingText, placement.Assets.Sponsored.Len, "")},
		})
	}

	// == hostname
	if placement.Assets.Hostname.Id != 0 {
		native.Assets = append(native.Assets, m.NativeAssets{
			Id:   placement.Assets.Hostname.Id,
			Data: &m.NativeData{Value: shortenString(campaign.BrandingText, placement.Assets.Hostname.Len, "")},
		})
	}

	// == description
	if placement.Assets.Description.Id != 0 {
		native.Assets = append(native.Assets, m.NativeAssets{
			Id:   placement.Assets.Description.Id,
			Data: &m.NativeData{Value: shortenString(creative.Title, placement.Assets.Description.Len, "…")}, // FIXME
		})
	}

	// == description2
	if placement.Assets.Description2.Id != 0 {
		native.Assets = append(native.Assets, m.NativeAssets{
			Id:   placement.Assets.Description2.Id,
			Data: &m.NativeData{Value: shortenString(creative.Title, placement.Assets.Description2.Len, "…")}, // FIXME
		})
	}

	// == rating
	if placement.Assets.Rating.Id != 0 {
		native.Assets = append(native.Assets, m.NativeAssets{
			Id:   placement.Assets.Rating.Id,
			Data: &m.NativeData{Value: "5"},
		})
	}

	// == likes
	if placement.Assets.Likes.Id != 0 {
		native.Assets = append(native.Assets, m.NativeAssets{
			Id:   placement.Assets.Likes.Id,
			Data: &m.NativeData{Value: strconv.FormatInt(int64(5000+rand.Intn(60000)), 10)},
		})
	}

	// == downloads
	if placement.Assets.Downloads.Id != 0 {
		native.Assets = append(native.Assets, m.NativeAssets{
			Id:   placement.Assets.Downloads.Id,
			Data: &m.NativeData{Value: strconv.FormatInt(int64(150000+rand.Intn(1350000)), 10)},
		})
	}

	// == call_to_action
	if placement.Assets.CallToAction.Id != 0 {
		native.Assets = append(native.Assets, m.NativeAssets{
			Id:   placement.Assets.CallToAction.Id,
			Data: &m.NativeData{Value: "Learn More"},
		})
	}

	// ==

	data, _ := json.Marshal(&m.ResponseNativeRequest{NativeRequest: native})

	return string(data)
}

func BuildResponse(request *m.Request) {

	ads, _ads := []m.Bid{}, []m.Ad0{}

	for x := range request.Placement {

		if len(request.Placement[x].Campaigns) == 0 {
			continue
		}

		campaigns := request.Placement[x].Campaigns
		placementId := request.Placement[x].Id

		for y := range campaigns {

			campaign := &campaigns[y]

			creative, ok := GetCreative(campaign)

			if !ok {
				continue
			}

			_ad := m.Ad0{
				AdvertiserId:     campaign.UserId,
				CampaignId:       campaign.CampaignId,
				CreativeId:       creative.CreativeId,
				PublisherId:      request.Seat.UserId,
				SeatId:           request.Seat.SeatId,
				InventoryId:      request.Targeting.Inventory,
				CountryId:        db.MapCountry.Get(request.Targeting.Country),
				DeviceId:         int(request.Targeting.PlatformId),
				PublisherBid:     campaign.PubRate,
				PublisherBidType: request.RateType,
			}

			switch campaign.BidType {
			case "CPM":
				_ad.AdvertiserSpendImpression = campaign.Bid / 1000.0
			case "CPC":
				_ad.AdvertiserSpendClick = campaign.Bid
			}

			switch _ad.PublisherBidType {
			case "CPM":
				_ad.PublisherRevenueImpression = _ad.PublisherBid / 1000.0
			case "CPC":
				_ad.PublisherRevenueClick = _ad.PublisherBid
			}

			_ad.PublisherBid = roundMoney(_ad.PublisherBid)
			_ad.AdvertiserSpendClick = roundMoney(_ad.AdvertiserSpendClick)
			_ad.PublisherRevenueClick = roundMoney(_ad.PublisherRevenueClick)
			_ad.AdvertiserSpendImpression = roundMoney(_ad.AdvertiserSpendImpression)
			_ad.PublisherRevenueImpression = roundMoney(_ad.PublisherRevenueImpression)

			// ==

			ad := m.Bid{
				Id:         strconv.FormatInt(int64(len(ads)+1), 10),
				ImpId:      placementId,
				Price:      _ad.PublisherBid,
				CampaignId: strconv.FormatInt(int64(campaign.CampaignId), 10),
				CreativeId: strconv.FormatInt(int64(creative.CreativeId), 10),
			}

			impressionLink, clickLink := buildEventLinks(request, &ad, &_ad, campaign, &creative)

			// fmt.Println(impressionLink)
			// fmt.Println(clickLink)

			ad.NURL = impressionLink
			ad.AdMarkup = buildAdm(request, &request.Placement[x], campaign, &creative, &_ad, clickLink)

			// fmt.Println(ad.AdMarkup)

			ads = append(ads, ad)
			_ads = append(_ads, _ad)
		}
	}

	request.Ads0 = _ads
	request.Ads1 = ads
}

// ==

func GetCreative(campaign *m.Campaign) (m.Creative, bool) {

	// fmt.Println("OptimizeCreatives", campaign.CampaignId)

	creativesLength := len(campaign.Creatives)

	switch creativesLength {
	case 0:
		return m.Creative{}, false

	case 1:
		return campaign.Creatives[0], true
	}

	if campaign.CreativeOptimizer == 0 {
		return campaign.Creatives[rand.Intn(creativesLength)], true
	}

	if rand.Float32() <= 0.75 {
		return campaign.Creatives[campaign.CreativesTop[rand.Intn(len(campaign.CreativesTop))]], true
	}

	return campaign.Creatives[rand.Intn(creativesLength)], true
}

//

func OptimizeCampaigns(request *m.Request, bidderData *m.BidderData, campaigns []m.Campaign) []m.Campaign {

	if request.RateType == "CPM" {

		features := request.FeaturesCtr

		for i := range campaigns {

			if campaigns[i].BidType == "CPC" {

				var ctr float64

				if inventoryCtr, ok := features.CampInv[campaigns[i].CampaignId]; ok {

					ctr = inventoryCtr

				} else if categoryCtr, ok := features.Categories[campaigns[i].Category]; ok {

					ctr = categoryCtr

				} else {

					ctr = features.CategoriesDefault
				}

				if campaigns[i].ImpFreq < len(bidderData.Features.Frequency) {

					ctr *= bidderData.Features.Frequency[campaigns[i].ImpFreq]
				}

				campaigns[i].PubRate = campaigns[i].Bid * ctr * 10

			} else {

				campaigns[i].PubRate = campaigns[i].Bid
			}

			campaigns[i].PubRate *= float64(request.Seat.Revshare) / 100.0
		}

		sort.Sort(m.ByPubRate(campaigns))

		return campaigns

	} else {

		// TODO
	}

	return nil
}

// ==

func FilterCampaigns(request *m.Request, bidderData *m.BidderData) (finalCampaigns []m.Campaign) {

	// features := bidderData.Features
	campaigns := bidderData.Campaigns
	found := false

	test := 0

	if request.Seat.Test == "1" {
		test = 1
	}

	for n := range campaigns {

		if test == 1 && campaigns[n].UserId != 50 {
			continue
		}

		if test == 0 && campaigns[n].UserId == 50 {
			continue
		}

		if campaigns[n].Rating > request.Seat.MinRating {
			// fmt.Println(campaigns[n].CampaignId, "Rating:", campaigns[n].Rating, request.Seat.MinRating)
			continue
		}

		if campaigns[n].Flowrate <= rand.Intn(500000) {
			// fmt.Println(campaigns[n].CampaignId, "Flowrate:", campaigns[n].Flowrate)
			continue
		}

		if len(campaigns[n].InventoryTargeting.Include) > 0 {

			if _, ok := campaigns[n].InventoryTargeting.Include[request.Targeting.Inventory]; !ok {
				// fmt.Println(campaigns[n].CampaignId, "Include", request.Targeting.Inventory, campaigns[n].InventoryTargeting.Include)
				continue
			}

		} else if len(campaigns[n].InventoryTargeting.Exclude) > 0 {

			if _, ok := campaigns[n].InventoryTargeting.Exclude[request.Targeting.Inventory]; ok {
				// fmt.Println(campaigns[n].CampaignId, "Exclude", request.Targeting.Inventory, campaigns[n].InventoryTargeting.Exclude)
				continue
			}
		}

		if len(campaigns[n].LocationTargeting) > 0 { // if(campaign.lt)

			found = false

			targeting := request.Targeting.Country + "|" + request.Targeting.Region + "|" + request.Targeting.City

			locationTargetingCriteria := [3]string{
				targeting[:len(request.Targeting.Country)],
				targeting[:len(request.Targeting.Country)+len(request.Targeting.Region)+1],
				targeting,
			}

			for _, location := range locationTargetingCriteria {
				if _, ok := campaigns[n].LocationTargeting[location]; ok {
					found = true
					// fmt.Println("LocationTargeting:", location)
					break
				}
			}

			if !found {
				// fmt.Println(campaigns[n].CampaignId, "LocationTargeting", targeting, campaigns[n].LocationTargeting)
				continue
			}
		}

		if len(campaigns[n].DeviceTargeting) > 0 { // if(campaign.dt)

			found = false

			platformName := browscap.PlatformName[request.Targeting.PlatformId]
			osName := browscap.OsName[request.Targeting.OsId]

			targeting := platformName + "|" + osName

			deviceTargetingCriteria := [2]string{
				targeting[:len(platformName)],
				targeting,
			}

			for _, deviceTargeting := range deviceTargetingCriteria {
				if _, ok := campaigns[n].DeviceTargeting[deviceTargeting]; ok {
					found = true
					// fmt.Println("DeviceTargeting:", deviceTargeting)
					break
				}
			}

			if !found {
				// fmt.Println(campaigns[n].CampaignId, "DeviceTargeting", targeting, campaigns[n].DeviceTargeting)
				continue
			}
		}

		finalCampaigns = append(finalCampaigns, campaigns[n])
	}

	// fmt.Println("finalCampaigns:", len(finalCampaigns))

	// ==

	finalCampaigns = db.GetFeatures(request, finalCampaigns)
	return
}
