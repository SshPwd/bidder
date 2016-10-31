package model

import (
	"github.com/valyala/fasthttp"
	"golang_bidder/browscap"
	"golang_bidder/ip2location"
)

/*
	var data = {
		id: null,
		ip: null,
		ua: null,
		site: null,
		app: null,
		channel: null,
		inventory_data: null,
		user: null,
		global_user: null,
		global_users: null,
		test: 0,
		placements: [],
		placementCount: 0,
		rate_type: rateType,
		targeting: {
			country: '',
			region: '',
			city: '',
			device: '',
			os: '',
			inventory: null
		},
		map: {
			device: null,
			country: null
		}
	}
*/

type (
	Channel     int
	AssetsData  struct{ Id, Len int }
	AssetsImage struct{ Id, Type, Width, Heigth int }

	Targeting struct {
		Inventory, Country, Region, City string
		OsId, PlatformId                 uint8
	}

	Assets struct {
		Image        AssetsImage // img
		Title        AssetsData  // title
		Sponsored    AssetsData  // 1
		Description  AssetsData  // 2
		Rating       AssetsData  // 3
		Likes        AssetsData  // 4
		Downloads    AssetsData  // 5
		Description2 AssetsData  // 10
		Hostname     AssetsData  // 11
		CallToAction AssetsData  // 12
	}

	Placement struct {
		Id             string
		HasDescription bool
		Count          int
		Floor          float64
		Assets         Assets
		Campaigns      []Campaign
	}

	InventoryData struct {
		Aimp int // aimp - bidRequest.bidRequest.placementCount
		Fl   int // fl   - bidRequest.bidRequest.placementFloors
		Req  int // req  - +1
		Imp  int // imp
		Clk  int // clk
		Bp   int // bp
		Wp   int // wp
		Rev  int // rev
		Spd  int // spd
	}

	FeaturesCtr struct {
		InventoryCtr      float64
		CategoriesTotal   float64
		CategoriesDefault float64
		Categories        map[string]float64
		CampInv           map[int]float64
	}

	Request struct {
		Ctx             *fasthttp.RequestCtx
		Browsers        *browscap.Browscap
		Geoip           *ip2location.IP2Location
		Datacenters     *ip2location.Datacenters
		BidRequest      BidRequest
		Seat            Seat
		Campaigns       []Campaign
		Channel         Channel
		Site            string
		App             string
		User            string
		GlobalUser      string
		RateType        string
		Targeting       Targeting
		InventoryData   InventoryData // FIXME
		FeaturesCtr     FeaturesCtr
		Placement       []Placement
		PlacementCount  int
		PlacementFloors float64
		Ads0            []Ad0
		Ads1            []Bid
	}
)

var (
	ChannelSite = Channel(1)
	ChannelApp  = Channel(2)
)

func (ch Channel) String() (name string) {
	switch ch {
	case ChannelSite:
		name = "site"
	case ChannelApp:
		name = "app"
	}
	return
}
