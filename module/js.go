package module

import (
	"encoding/json"
	m "golang_bidder/model"
)

type (
	Js struct{}
)

func (_ Js) Name() string {

	return "js"
}

func (_ Js) Read(request *m.Request) error {

	request.Ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	request.Ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST")

	if err := m.ParseBidRequest(request); err != nil {
		return err
	}

	request.BidRequest.Device.Ip = request.Ctx.RemoteIP().String()
	request.BidRequest.Device.UserAgent = string(request.Ctx.UserAgent())

	return nil
}

func (_ Js) Write(request *m.Request) error {

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

func (_ Js) Win(_, revenuePrice float64) float64 {

	return revenuePrice
}
