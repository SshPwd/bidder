package module

import (
	"encoding/json"
	"errors"
	m "golang_bidder/model"
	"math"
)

type (
	Internal struct{}
)

var (
	ErrJsonMarshal = errors.New("ErrJsonMarshal")
)

func (_ Internal) Name() string {

	return "internal"
}

func (_ Internal) Read(request *m.Request) error {

	request.Ctx.Response.Header.Set("x-openrtb-version", "2.3")

	return m.ParseBidRequest(request)
}

func (_ Internal) Write(request *m.Request) error {

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

func (_ Internal) Win(auctionPrice, revenuePrice float64) float64 {

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
