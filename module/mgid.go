package module

import (
	"encoding/json"
	m "golang_bidder/model"
	"math"
)

type (
	Mgid struct{}
)

func (_ Mgid) Name() string {

	return "mgid"
}

func (_ Mgid) Read(request *m.Request) error {

	request.Ctx.Response.Header.Set("x-openrtb-version", "2.3")

	if err := m.ParseBidRequest(request); err != nil {
		return err
	}

	request.Seat.MinRating = 3

	if len(request.BidRequest.Ext.AdTypes) > 0 {

		request.Seat.MinRating = 0

		for n := range request.BidRequest.Ext.AdTypes {

			switch request.BidRequest.Ext.AdTypes[n] {
			case "dark", "shock":

				request.Seat.MinRating = 1
				break
			}
		}
	}

	return nil
}

func (_ Mgid) Write(request *m.Request) error {

	request.Ctx.Response.Header.Set("Content-Type", "application/json")

	bidResponse := m.BidResponse{
		Id:       request.BidRequest.Id,
		Currency: "USD",
	}

	curIndex, curPlacement := -1, ""

	for _, ad := range request.Ads1 {

		if curPlacement != ad.ImpId {

			curIndex++
			curPlacement = ad.ImpId
			bidResponse.SeatBid = append(bidResponse.SeatBid, m.SeatBid{Bid: []m.Bid{ad}})

		} else {

			bidResponse.SeatBid[curIndex].Bid = append(bidResponse.SeatBid[curIndex].Bid, ad)
		}
	}

	data, err := json.Marshal(&bidResponse)

	if err != nil {
		// TODO write error
		return ErrJsonMarshal
	}

	request.Ctx.Write(data)
	return nil
}

func (_ Mgid) Win(auctionPrice, revenuePrice float64) float64 {

	if !math.IsNaN(auctionPrice) && !math.IsInf(revenuePrice, 0) {

		if auctionPrice > 0 {

			auctionPrice /= 1000.0

			if auctionPrice < revenuePrice {
				return auctionPrice
			}
		}
	}

	return revenuePrice
}
