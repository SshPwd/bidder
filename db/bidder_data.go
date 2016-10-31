package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang_bidder/config"
	"golang_bidder/logs"
	m "golang_bidder/model"
	"golang_bidder/module"
	"golang_bidder/utils"
	"sync/atomic"
	"time"
)

const (
	bidderDataTTL = 1
)

var (
	globalBidderData atomic.Value
	globalSeat       atomic.Value

	modules = map[string]m.Module{
		"JS":       module.Js{},
		"mgid":     module.Mgid{},
		"Internal": module.Internal{},
	}
)

var (
	ErrSeatNotFound   = errors.New("ErrSeatNotFound")
	ErrSeatInvalidKey = errors.New("ErrSeatInvalidKey")
)

// ==

func init() {

	globalBidderData.Store(&m.BidderData{})
	globalSeat.Store(map[int]m.Seat{})

	go updateBidderData()
}

// ==

func updateBidderData() {

	// TODO recover

	for {
		time.Sleep(bidderDataTTL * time.Minute)
		LoadBidderData()
	}
}

// ==

func LoadBidderData() {

	bidderData := m.BidderData{}

	conf := config.Get()

	buf, err := utils.HttpGet(conf.BidderDataUrl, 30*time.Second)
	if err != nil {
		logs.Critical(err.Error())
		return
	}
	defer utils.PutBuffer(buf)

	if err := json.Unmarshal(buf.Bytes(), &bidderData); err != nil {
		logs.Critical(err.Error())
		return
	}

	// ==

	if frequency, ok := bidderData.Features.FrequencyData.([]float64); ok {

		bidderData.Features.Frequency = frequency

	} else if frequency, ok := bidderData.Features.FrequencyData.(map[string]float64); ok {

		maxIndex := 0

		for key, _ := range frequency {
			if n := utils.ToInt(key); n > maxIndex {
				maxIndex = n
			}
		}

		if maxIndex > 0 {

			array := make([]float64, maxIndex+1)

			for key, val := range frequency {

				n := utils.ToInt(key)
				array[n] = val
			}

			bidderData.Features.Frequency = array
		}
	}

	// bidderData.Features.Frequency = []float64{1.0, 1.0, 2.4269}

	for n := range bidderData.Campaigns {

		var value interface{}

		err := json.Unmarshal(bidderData.Campaigns[n].InventoryTargetingRaw, &value)

		bidderData.Campaigns[n].InventoryTargetingRaw = nil

		if err != nil {
			logs.Critical(err.Error())
			continue
		}

		if inventoryTargeting, ok := value.(map[string]interface{}); ok {

			if exclude, ok := inventoryTargeting["exclude"].(map[string]interface{}); ok {
				bidderData.Campaigns[n].InventoryTargeting.Exclude = map[string]int{}
				for key, _ := range exclude {
					bidderData.Campaigns[n].InventoryTargeting.Exclude[key] = 1
				}
			}

			if include, ok := inventoryTargeting["include"].(map[string]interface{}); ok {
				bidderData.Campaigns[n].InventoryTargeting.Include = map[string]int{}
				for key, _ := range include {
					bidderData.Campaigns[n].InventoryTargeting.Include[key] = 1
				}
			}
		}
	}

	seats := map[int]m.Seat{}

	for _, seat := range bidderData.Seats {

		if seatModule, ok := modules[seat.ModuleName]; ok {

			fmt.Println(seat.SeatId, seat.SecretKey)

			seat.Module = seatModule
			seats[seat.SeatId] = seat

		} else {

			logs.Critical(fmt.Sprintf("SeatId=%d module %q not found", seat.SeatId, seat.ModuleName))
		}
	}

	for n := range bidderData.Campaigns {
		bidderData.Campaigns[n].OptimizeCreatives()
	}

	globalSeat.Store(seats)
	globalBidderData.Store(&bidderData)
}

// ==

func GetBidderData() *m.BidderData {

	return globalBidderData.Load().(*m.BidderData)
}

func GetCampaigns() (list []m.Campaign) {

	bidderData := globalBidderData.Load().(*m.BidderData)

	return bidderData.Campaigns
}

// ==

func ExistSeatId(seatId int) error {

	seatById := globalSeat.Load().(map[int]m.Seat)

	if _, ok := seatById[seatId]; ok {

		return nil
	}

	return ErrSeatNotFound
}

func GetSeatById(seatId int) (m.Seat, error) {

	seatById := globalSeat.Load().(map[int]m.Seat)

	if seat, ok := seatById[seatId]; ok {

		return seat, nil
	}

	return m.Seat{}, ErrSeatNotFound
}

func GetSeatById2(seatId int, secretKey []byte) (m.Seat, error) {

	seatById := globalSeat.Load().(map[int]m.Seat)

	if seat, ok := seatById[seatId]; ok {

		if seat.SecretKey == string(secretKey) {

			return seat, nil

		} else {

			return m.Seat{}, ErrSeatInvalidKey
		}
	}

	return m.Seat{}, ErrSeatNotFound
}
