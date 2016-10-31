package model

type Ad0 struct {
	AdvertiserId               int     `json:"advertiser_id"`
	CampaignId                 int     `json:"campaign_id"`
	CreativeId                 int     `json:"creative_id"`
	PublisherId                int     `json:"publisher_id"`
	SeatId                     int     `json:"seat_id"`
	CountryId                  int     `json:"country_id"`
	DeviceId                   int     `json:"device_id"`
	InventoryId                string  `json:"inventory_id"`
	PublisherBid               float64 `json:"publisher_bid"`
	PublisherBidType           string  `json:"publisher_bid_type"`
	AdvertiserSpendImpression  float64 `json:"advertiser_spend_impression"`
	PublisherRevenueImpression float64 `json:"publisher_revenue_impression"`
	AdvertiserSpendClick       float64 `json:"advertiser_spend_click"`
	PublisherRevenueClick      float64 `json:"publisher_revenue_click"`
}

// type Ad1 struct {
// 	Id    string  `json:"id"`
// 	ImpId string  `json:"impid"`
// 	Price float64 `json:"price"`
// 	Cid   string  `json:"cid"`
// 	Crid  string  `json:"crid"`
// 	Nurl  string  `json:"nurl"`
// 	Adm   string  `json:"adm"`
// }
