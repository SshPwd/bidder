package model

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
)

var _ = fmt.Print

type (
	Bid struct {
		Price      float64  `json:"price"`
		H          int      `json:"h,omitempty"`
		W          int      `json:"w,omitempty"`
		Weight     int      `json:"-"`
		Id         string   `json:"id"`
		ImpId      string   `json:"impid"`
		AdId       string   `json:"adid,omitempty"`
		NURL       string   `json:"nurl,omitempty"`
		AdMarkup   string   `json:"adm,omitempty"`
		Bundle     string   `json:"bundle,omitempty"`
		IURL       string   `json:"iurl,omitempty"`
		CampaignId string   `json:"cid,omitempty"`
		CreativeId string   `json:"crid,omitempty"`
		DealId     string   `json:"deal_id,omitempty"`
		Cat        []string `json:"cat,omitempty"`
		Attr       []int    `json:"attr,omitempty"`
		AdvDomain  []string `json:"adomain,omitempty"`
	}

	SeatBid struct {
		Bid   []Bid  `json:"bid"`
		Seat  string `json:"seat,omitempty"`
		Group int    `json:"group"`
	}

	BidResponse struct {
		Id         string    `json:"id"`
		BidId      string    `json:"bidid,omitempty"`
		Currency   string    `json:"cur,omitempty"`
		CustomData string    `json:"customdata,omitempty"`
		Nbr        int       `json:"nbr,omitempty"`
		SeatBid    []SeatBid `json:"seatbid"`
	}

	//

	NativeTitle struct {
		Len  int    `json:"len,omitempty"`
		Text string `json:"text,omitempty"`
	}

	NativeImg struct {
		WMin   int      `json:"wmin,omitempty"`
		HMin   int      `json:"hmin,omitempty"`
		Width  int      `json:"w,omitempty"`
		Height int      `json:"h,omitempty"`
		Type   int      `json:"type,omitempty"`
		Mimes  []string `json:"mimes,omitempty"`
		Url    string   `json:"url,omitempty"`
	}

	NativeVideo struct {
	}

	NativeData struct {
		Type  int         `json:"type,omitempty"`
		Len   int         `json:"len,omitempty"`
		Value interface{} `json:"value,omitempty"`
	}

	NativeLink struct {
		Url string `json:"url,omitempty"`
	}

	NativeAssets struct {
		Id       int          `json:"id"`
		Required int          `json:"required,omitempty"`
		Type     int          `json:"type,omitempty"`
		Title    *NativeTitle `json:"title,omitempty"`
		Img      *NativeImg   `json:"img,omitempty"`
		Data     *NativeData  `json:"data,omitempty"`
		Video    *NativeVideo `json:"video,omitempty"`
	}

	NativeRequest struct {
		Version     string         `json:"ver"`
		Layout      int            `json:"layout,omitempty"`
		Adunit      int            `json:"adunit,omitempty"`
		Plcmtcnt    int            `json:"plcmtcnt,omitempty"`
		Link        NativeLink     `json:"link,omitempty"`
		Assets      []NativeAssets `json:"assets"`
		Imptrackers []string       `json:"imptrackers,omitempty"`
	}

	ResponseNativeRequest struct {
		NativeRequest NativeRequest `json:"native"`
	}

	// ==

	Device struct {
		Ip        string `json:"ip"`
		UserAgent string `json:"ua"`
	}

	User struct {
		Id         string `json:"id"`
		BuyerId    string `json:"buyerid,omitempty"`
		CustomData string `json:"customdata,omitempty"`
	}

	DspExt struct {
		AdTypes []string `json:"ad_types"`
	}

	Site struct {
		Id     string `json:"id"`
		Domain string `json:"domain,omitempty"`
		Page   string `json:"page,omitempty"`
	}

	App struct {
		Name   string `json:"name,omitempty"`
		Bundle string `json:"bundle,omitempty"`
	}

	Native struct {
		Request string `json:"request"`
		Ver     string `json:"ver,omitempty"`
		BAttr   []int  `json:"battr,omitempty"`
	}

	Impression struct {
		Id       string  `json:"id"`
		BidFloor float64 `json:"bidfloor"`
		Secure   int     `json:"secure"`
		Native   Native  `json:"native,omitempty"`
		BAttr    []int   `json:"battr,omitempty"`
	}

	BidRequest struct {
		Id       string       `json:"id"`
		Device   Device       `json:"device,omitempty"`
		User     User         `json:"user,omitempty"`
		Ext      DspExt       `json:"ext,omitempty"`
		Site     Site         `json:"site,omitempty"`
		App      App          `json:"app,omitempty"`
		Currency []string     `json:"cur,omitempty"`
		BCat     []string     `json:"bcat,omitempty"`
		Imp      []Impression `json:"imp,omitempty"`
		Test     int          `json:"test,omitempty"`
	}
)

var (
	ErrCrawler                 = errors.New("ErrCrawler")
	ErrDataCenter              = errors.New("ErrDataCenter")
	ErrDataCenter2             = errors.New("ErrDataCenter2")
	ErrIpNotFound              = errors.New("ErrIpNotFound")
	ErrCannotBeResolved        = errors.New("ErrCannotBeResolved")
	ErrUserAgentNotFound       = errors.New("ErrUserAgentNotFound")
	ErrNoDeviceFieldsAvailable = errors.New("ErrNoDeviceFieldsAvailable")

	ErrBidId                 = errors.New("ErrBidId")
	ErrBidDevice             = errors.New("ErrBidDevice")
	ErrAppObject             = errors.New("ErrAppObject")
	ErrSiteObject            = errors.New("ErrSiteObject")
	ErrEmptyAssets           = errors.New("ErrEmptyAssets")
	ErrBidCurrency           = errors.New("ErrBidCurrency")
	ErrBidImpression         = errors.New("ErrBidImpression")
	ErrNativeUnmarshal       = errors.New("ErrNativeUnmarshal")
	ErrMissingSiteOrApp      = errors.New("ErrMissingSiteOrApp")
	ErrInvalidAssetsType     = errors.New("ErrInvalidAssetsType")
	ErrVideoAssetsUnbiddable = errors.New("ErrVideoAssetsUnbiddable")
)

// ==

var (
	bufferPool = sync.Pool{
		New: func() interface{} {
			buf := new(bytes.Buffer)
			buf.Grow(64)
			return buf
		},
	}
)

func getBuffer() (buf *bytes.Buffer) {
	return bufferPool.Get().(*bytes.Buffer)
}

func putBuffer(buf *bytes.Buffer) {
	buf.Truncate(0)
	bufferPool.Put(buf)
}

// ==

var testJson = []byte(`{"id":"1JEJJAE3Ow47SurJOzKV8G","device":{"ip":"90.255.169.124","ua":"Mozilla/5.0 (iPhone; CPU iPhone OS 9_2_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13D15 Safari/601.1"},"user":{"id":"","buyerid":""},"ext":{"ad_types":null},"site":{"id":"27856","domain":"testmenow.co.uk","page":"http://putlocker.is/watch-grown-ups-2-online-free-putlocker.html"},"cur":["USD","RUB","UAH"],"imp":[{"id":"1","bidfloor":0.23653958191554902,"secure":0,"native":{"request":"{\"ver\":\"1\",\"layout\":1,\"adunit\":2,\"plcmtcnt\":5,\"assets\":[{\"id\":1,\"required\":1,\"title\":{\"len\":65}},{\"id\":2,\"required\":1,\"img\":{\"wmin\":200,\"hmin\":200,\"type\":3,\"mimes\":[\"image/jpg\",\"image/png\"]}},{\"id\":3,\"required\":0,\"data\":{\"type\":2,\"len\":75}},{\"id\":6,\"required\":0,\"data\":{\"type\":11}}]}\n","ver":"1","battr":[1,2,3,4,5,6,8,9,10,14]}}]}`)

/*`{
  "id": "e87c12f3-b46a-48f4-bafd-1c93ee6a0d88",
  "imp": [
	{
	  "id": "1",
	  "bidfloor": 1.6,
	  "native": {
		"request": "{\"plcmtcnt\":5,\"assets\":[{\"id\":1,\"required\":1,\"title\":{\"len\":30}},{\"id\":2,\"required\":1,\"data\":{\"type\":2}},{\"id\":3,\"required\":1,\"img\":{\"type\":2,\"h\":50,\"w\":50}},{\"id\":5,\"required\":1,\"data\":{\"type\":1}}]}"
	  }
	},
	{
	  "id": "2",
	  "bidfloor": 1.7,
	  "native": {
		"request": "{\"plcmtcnt\":5,\"assets\":[{\"id\":1,\"required\":1,\"title\":{\"len\":30}},{\"id\":2,\"required\":1,\"data\":{\"type\":2}},{\"id\":3,\"required\":1,\"img\":{\"type\":2,\"h\":50}},{\"id\":5,\"required\":1,\"data\":{\"type\":1}}]}"
	  }
	}
  ],
  "test": 0,
  "device": {
	"ip": "71.167.36.38",
	"ua": "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.97 Safari/537.36"
  },
  "site": {
	"page": "http://gabworthy.com/",
	"domain": "gabworthy.com"
  }
}`*/

func ParseBidRequest(request *Request) error {

	/*request.Ctx.PostBody()*/

	body := request.Ctx.PostBody()

	if len(body) == 0 {
		body = testJson
	}

	if err := json.Unmarshal(body, &request.BidRequest); err != nil {
		return err
	}

	bid := &request.BidRequest

	if bid.Id == "" {
		return ErrBidId
	}

	if len(bid.Imp) == 0 {
		return ErrBidImpression
	}

	var err error

	if err = bid.impressionsParse(request); err != nil {
		return err
	}

	if err = bid.channelParse(request); err != nil {
		return err
	}

	if err = bid.userParse(request); err != nil {
		return err
	}

	if err = bid.deviceParse(request); err != nil {
		return err
	}

	if err = bid.geoParse(request); err != nil {
		return err
	}

	// ==

	buf := getBuffer()
	buf.Grow(len(bid.Device.Ip) + len(bid.Device.UserAgent))
	buf.WriteString(bid.Device.Ip)
	buf.WriteString(bid.Device.UserAgent)
	sum := md5.Sum(buf.Bytes())
	putBuffer(buf)

	request.GlobalUser = hex.EncodeToString(sum[:])

	// ==

	// return app.globals.aerospike.inventory.get(bidRequest,'request',function(inventory_data){
	// 	bidRequest.bidRequest.inventory_data = inventory_data;
	// 	return cb();
	// });

	// ==

	return nil
}

func (bid *BidRequest) deviceParse(request *Request) error {

	if bid.Device.UserAgent == "" {
		return ErrNoDeviceFieldsAvailable
	}

	if deviceInfo := request.Browsers.FindStr(bid.Device.UserAgent); deviceInfo.Ok {

		if deviceInfo.Crawler {
			return ErrCrawler
		}

		request.Targeting.OsId = deviceInfo.OsId
		request.Targeting.PlatformId = deviceInfo.PlatformId
		return nil
	}

	return ErrUserAgentNotFound
}

func (bid *BidRequest) geoParse(request *Request) error {

	if bid.Device.Ip == "" {
		return ErrCannotBeResolved
	}

	if request.Datacenters.Find(bid.Device.Ip) {
		return ErrDataCenter2
	}

	if geo := request.Geoip.Find(bid.Device.Ip); geo.Ok {

		if geo.IsHosting {
			return ErrDataCenter
		}

		request.Targeting.City = geo.City
		request.Targeting.Region = geo.Region
		request.Targeting.Country = geo.Country

		// fmt.Println("geoParse:", geo.City, geo.Region, geo.Country)
		return nil
	}

	return ErrIpNotFound
}

func (bid *BidRequest) impressionsParse(request *Request) error {

	for n := range bid.Imp {

		if bid.Imp[n].Id != "" && bid.Imp[n].Native.Request != "" {

			nativeRequest := NativeRequest{}

			if err := json.Unmarshal([]byte(bid.Imp[n].Native.Request), &nativeRequest); err != nil {
				return ErrNativeUnmarshal
			}

			if len(nativeRequest.Assets) == 0 {
				return ErrEmptyAssets
			}

			assets := Assets{}
			hasDescription := false

			for i := range nativeRequest.Assets {

				if nativeRequest.Assets[i].Video != nil {
					return ErrVideoAssetsUnbiddable
				}

				if nativeRequest.Assets[i].Id != 0 {

					if nativeRequest.Assets[i].Title != nil {

						assets.Title.Id = nativeRequest.Assets[i].Id
						assets.Title.Len = nativeRequest.Assets[i].Title.Len
					}

					if nativeRequest.Assets[i].Img != nil {

						assets.Image.Id = nativeRequest.Assets[i].Id

						img := nativeRequest.Assets[i].Img

						if img.Type != 0 {
							assets.Image.Type = img.Type
						} else {
							assets.Image.Type = 3
						}

						if img.Width != 0 {
							assets.Image.Width = img.Width
						} else {
							assets.Image.Width = img.WMin
						}

						if img.Height != 0 {
							assets.Image.Heigth = img.Height
						} else {
							assets.Image.Heigth = img.HMin
						}
					}

					if nativeRequest.Assets[i].Data != nil {

						switch nativeRequest.Assets[i].Data.Type {
						case 1:
							assets.Sponsored.Id = nativeRequest.Assets[i].Id
							assets.Sponsored.Len = nativeRequest.Assets[i].Data.Len
						case 2:
							if nativeRequest.Assets[i].Required == 1 {
								hasDescription = true
								assets.Description.Id = nativeRequest.Assets[i].Id
								assets.Description.Len = nativeRequest.Assets[i].Data.Len
							}
						case 3:
							if nativeRequest.Assets[i].Required == 1 {
								assets.Rating.Id = nativeRequest.Assets[i].Id
								assets.Rating.Len = nativeRequest.Assets[i].Data.Len
							}
						case 4:
							if nativeRequest.Assets[i].Required == 1 {
								assets.Likes.Id = nativeRequest.Assets[i].Id
								assets.Likes.Len = nativeRequest.Assets[i].Data.Len
							}
						case 5:
							if nativeRequest.Assets[i].Required == 1 {
								assets.Downloads.Id = nativeRequest.Assets[i].Id
								assets.Downloads.Len = nativeRequest.Assets[i].Data.Len
							}
						case 10:
							if nativeRequest.Assets[i].Required == 1 {
								assets.Description2.Id = nativeRequest.Assets[i].Id
								assets.Description2.Len = nativeRequest.Assets[i].Data.Len
							}
						case 11:
							assets.Hostname.Id = nativeRequest.Assets[i].Id
							assets.Hostname.Len = nativeRequest.Assets[i].Data.Len
						case 12:
							if nativeRequest.Assets[i].Required == 1 {
								assets.CallToAction.Id = nativeRequest.Assets[i].Id
								assets.CallToAction.Len = nativeRequest.Assets[i].Data.Len
							}
						default:

							// fmt.Println("skip AssetsId:", nativeRequest.Assets[i].Data.Type)

							if nativeRequest.Assets[i].Required == 1 {
								return ErrInvalidAssetsType
							}
						}
					}
				}
			}

			placement := Placement{
				Id:             bid.Imp[n].Id,
				Count:          1,
				Floor:          bid.Imp[n].BidFloor,
				Assets:         assets,
				HasDescription: hasDescription,
			}

			// fmt.Println("bid.Imp[n].BidFloor", bid.Imp[n].BidFloor)

			if nativeRequest.Plcmtcnt > 1 {
				placement.Count = nativeRequest.Plcmtcnt
			}

			request.PlacementCount += placement.Count
			request.PlacementFloors += bid.Imp[n].BidFloor

			request.Placement = append(request.Placement, placement)
		}
	}
	return nil
}

func (bid *BidRequest) channelParse(request *Request) error {

	keyBuf := getBuffer()

	if bid.Site.Domain != "" || bid.Site.Page != "" || bid.Site.Id != "" {

		if bid.Site.Domain != "" {

			request.Site = bid.Site.Domain

		} else if bid.Site.Page != "" {

			if parsedUrl, err := url.Parse(trim(bid.Site.Page)); err == nil {
				request.Site = parsedUrl.Host
			}

		} else if bid.Site.Id != "" {

			request.Site = bid.Site.Id
		}

		request.Site = trim(request.Site)

		if request.Site == "" {

			return ErrSiteObject
		}

		request.Site = strings.ToLower(request.Site)
		request.Channel = ChannelSite

		keyBuf.WriteString(request.Channel.String())
		keyBuf.WriteString("_")
		keyBuf.WriteString(request.Site)

	} else if bid.App.Name != "" || bid.App.Bundle != "" {

		if bid.App.Name != "" {
			request.App = bid.App.Name
		} else if bid.App.Bundle != "" {
			request.App = bid.App.Bundle
		} else {
			return ErrAppObject
		}

		request.Channel = ChannelApp

		keyBuf.WriteString(request.Channel.String())
		keyBuf.WriteString("_")
		keyBuf.WriteString(request.App)

	} else {

		return ErrMissingSiteOrApp
	}

	sum := md5.Sum(keyBuf.Bytes())

	keyBuf.Truncate(0)
	hashBuf := keyBuf

	encoder := base64.NewEncoder(base64.RawURLEncoding, hashBuf)
	encoder.Write(sum[:])
	encoder.Close()

	// fmt.Println("Targeting.Inventory:", hashBuf.String())

	request.Targeting.Inventory = hashBuf.String()
	return nil
}

func (bid *BidRequest) userParse(request *Request) error {
	if bid.User.Id != "" {
		request.User = bid.User.Id
	}
	return nil
}

func trim(str string) string {
	return strings.TrimSuffix(strings.TrimPrefix(str, " "), " ")
}
