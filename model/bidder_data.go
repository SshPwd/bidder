package model

import (
	"encoding/json"
	"math"
)

type (
	Creative struct {
		CreativeId       int     `json:"crid"`
		ShortTitle       string  `json:"st"`
		Title            string  `json:"t"`
		Image            string  `json:"i"`
		Link             string  `json:"l"` // landing page
		FinalLink        string  `json:"-"`
		Ctr              float64 `json:"c"`
		CreativeApproval int     `json:"r"` // 0 = G , 1 = G+PG, 2 = G+PG+PG13 , 3 = G+PG+PG13+R
	}

	InventoryTargeting struct {
		Include map[string]int
		Exclude map[string]int
	}

	Campaign struct {
		CampaignId            int                `json:"cid"`
		UserId                int                `json:"uid"`
		BrandingText          string             `json:"br"`
		BrandLogo             string             `json:"bl"`
		UrlParameters         string             `json:"up"`
		Category              string             `json:"cat"`
		Flowrate              int                `json:"f"`
		BidType               string             `json:"bt"`
		Bid                   float64            `json:"bid"`
		CreativeOptimizer     int                `json:"co"`
		InventoryOtimizer     int                `json:"io"`
		Rating                int                `json:"r"` //0 = G , 1 = G+PG, 2 = G+PG+PG13 , 3 = G+PG+PG13+R
		CategoryTargeting     string             `json:"ct"`
		ImpressionPixel       string             `json:"ip"`
		LocationTargeting     map[string]int     `json:"lt"` // location
		DeviceTargeting       map[string]int     `json:"dt"` // device
		Creatives             []Creative         `json:"cr"`
		CreativesTop          []int              `json:"-"`
		InventoryTargetingRaw json.RawMessage    `json:"it"`
		InventoryTargeting    InventoryTargeting `json:"-"`
		ImpFreq               int                `json:"-"`
		Optimized             int                `json:"-"`
		PubRate               float64            `json:"-"`
	}

	FeatureCtr      map[int]float32
	FeatureRevshare map[string]float32

	Features struct {
		Frequency     []float64          `json:"-"`
		Ctr           map[int]FeatureCtr `json:"ctr"`
		FrequencyData interface{}        `json:"frequency"`
		Revshares     interface{}        `json:"revshares"` // FIXME
	}

	Seat struct {
		SeatId          int     `json:"sid"`
		UserId          int     `json:"uid"`
		ModuleName      string  `json:"m"`
		SecretKey       string  `json:"s"`
		A               string  `json:"a"`
		B               string  `json:"b"`
		C               string  `json:"c"`
		D               string  `json:"d"`
		E               string  `json:"e"`
		F               string  `json:"f"`
		NoWinNotice     int     `json:"nwn"`
		Revshare        int     `json:"r"`
		IgnoreFloor     int     `json:"ig"`
		MinRating       int     `json:"mr"` // 0 = G , 1 = G+PG, 2 = G+PG+PG13 , 3 = G+PG+PG13+R
		LearnCtrContent float32 `json:"lcc"`
		LearnCtrOffer   float32 `json:"lco"`
		LearnCtrApp     float32 `json:"lca"`
		Module          Module  `json:"-"`
		Test            string  `json:"test"`
	}

	Inventory struct {
		Iid int `json:"iid"`
		S   int `json:"s"`
	}

	BidderData struct {
		Date      string          `json:"date"`
		Features  Features        `json:"features"`
		Seats     map[string]Seat `json:"seats"`
		Campaigns []Campaign      `json:"campaigns"`
		Inventory []Inventory     `json:"inventory"`
	}

	ByPubRate []Campaign
)

// ==

func (a ByPubRate) Len() int           { return len(a) }
func (a ByPubRate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPubRate) Less(i, j int) bool { return a[i].PubRate > a[j].PubRate }

// ==

func (campaign *Campaign) OptimizeCreatives() {

	var ctr float64
	var maxCtr float64 = -1.0
	var topCreatives = make([]int, 0, 8)

	for i := range campaign.Creatives {

		ctr = math.Floor(campaign.Creatives[i].Ctr*100.0) / 100.0

		if ctr > maxCtr {

			maxCtr = ctr
			topCreatives = topCreatives[:0]
			topCreatives = append(topCreatives, i)

		} else if ctr == maxCtr {

			topCreatives = append(topCreatives, i)
		}
	}

	campaign.CreativesTop = topCreatives
}
