package db

import (
	"fmt"
	"golang_bidder/config"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"sync"
)

const (
	maxHits = 1000
)

type (
	AggregateKey struct {
		Date         string `bson:"date"`
		Datetime     string `bson:"datetime"`
		InventoryId  string `bson:"inventory_id"`
		AdvertiserId int    `bson:"advertiser_id"`
		CampaignId   int    `bson:"campaign_id"`
		CountryId    int    `bson:"country_id"`
		CreativeId   int    `bson:"creative_id"`
		DeviceId     int    `bson:"device_id"`
		Frequency    int    `bson:"frequency"`
		Position     int    `bson:"position"`
		PublisherId  int    `bson:"publisher_id"`
		SeatId       int    `bson:"seat_id"`
		Test         int    `bson:"test"`
	}

	AggregateMetric struct {
		Actions            int     `bson:"actions"`
		Bids               int     `bson:"bids"`
		Clicks             int     `bson:"clicks"`
		Impressions        int     `bson:"impressions"`
		RequestErrors      int     `bson:"request_errors"`
		Requests           int     `bson:"requests"`
		Score              int     `bson:"score"`
		PhoneImpressions   int     `bson:"phone_impressions,omitempty"`
		TabletImpressions  int     `bson:"tablet_impressions,omitempty"`
		DesktopImpressions int     `bson:"desktop_impressions,omitempty"`
		Revenue            float64 `bson:"revenue"`
		RevenueAuction     float64 `bson:"revenue_auction"`
		RevenueWon         float64 `bson:"revenue_won"`
		Spend              float64 `bson:"spend"`
		Floor              float64 `bson:"floor,omitempty"`
	}

	Aggregated struct {
		Hits    int
		Metrics map[AggregateKey]*AggregateMetric
	}

	AggregateData struct {
		Table  string
		Key    AggregateKey
		Setric AggregateMetric
	}
)

var (
	aggregated = map[string]*Aggregated{}
	aggrSync   sync.Mutex
)

func AggregateV(aggregateData []AggregateData) {

	// TODO optimize
}

func Aggregate(table string, keys AggregateKey, metric AggregateMetric) {

	// fmt.Println("Aggregate:", table)

	aggrSync.Lock()

	aggrPtr, ok := aggregated[table]
	if !ok {
		aggrPtr = &Aggregated{Hits: 1, Metrics: map[AggregateKey]*AggregateMetric{}}
	} else {
		aggrPtr.Hits++
	}

	metricPtr, ok := aggrPtr.Metrics[keys]
	if !ok {
		metricPtr = &AggregateMetric{}
	}

	metricPtr.Bids += metric.Bids
	metricPtr.Spend += metric.Spend
	metricPtr.Score += metric.Score
	metricPtr.Clicks += metric.Clicks
	metricPtr.Actions += metric.Actions
	metricPtr.Revenue += metric.Revenue
	metricPtr.Requests += metric.Requests
	metricPtr.RevenueWon += metric.RevenueWon
	metricPtr.Impressions += metric.Impressions
	metricPtr.RequestErrors += metric.RequestErrors
	metricPtr.RevenueAuction += metric.RevenueAuction

	aggrPtr.Metrics[keys] = metricPtr

	if aggrPtr.Hits > maxHits {

		go pushToMongo(table, aggrPtr)
		delete(aggregated, table)

	} else {

		aggregated[table] = aggrPtr
	}

	aggrSync.Unlock()
}

// ==

var indexNames = []string{
	"date",
	"datetime",
	"inventory_id",
	"advertiser_id",
	"campaign_id",
	"country_id",
	"creative_id",
	"device_id",
	"frequency",
	"position",
	"publisher_id",
	"seat_id",
	"test",
}

func pushToMongo(table string, aggrPtr *Aggregated) {

	sess := globalMgoSession.Clone()
	defer sess.Close()

	sess.Refresh()

	conf := config.Get()

	collection := sess.DB(conf.Mongodb.DbName).C(table)

	bulk := collection.Bulk()

	for key, metric := range aggrPtr.Metrics {
		bulk.Upsert(key, bson.M{"$inc": metric})
	}

	if result, err := bulk.Run(); err != nil {

		fmt.Println("try.1", table, result, err, len(aggrPtr.Metrics))

		sess.Refresh()

		collection = sess.DB(conf.Mongodb.DbName).C(table)

		collection.EnsureIndex(mgo.Index{
			Key:        indexNames,
			Unique:     false,
			Background: true,
		})

		bulk = collection.Bulk()

		for key, metric := range aggrPtr.Metrics {
			bulk.Upsert(key, bson.M{"$inc": metric})
		}

		if result, err := bulk.Run(); err != nil {
			fmt.Println("try.2", table, result, err, len(aggrPtr.Metrics))
		}
	}
}
