package db

import (
	"fmt"
	aerospike "github.com/aerospike/aerospike-client-go"
	m "golang_bidder/model"
	"time"
)

// ==

const (
	inventoryExpire = 10 * 24 * 60 * 60
)

// aimp - bidRequest.bidRequest.placementCount
// fl - bidRequest.bidRequest.placementFloors
// req  = 1

// imp
// clk
// bp
// wp
// rev
// spd

type (
	InventoryData struct {
		imp int
		clk int
		bp  int
		wp  int
		rev int
		spd int
	}

	InventoryData2 struct {
		aimp int
		req  int
		fl   float64
	}
)

func SaveInventoryData(request *m.Request) bool {

	key, err := aerospike.NewKey("bidder", "inventory", request.Targeting.Inventory)
	if err != nil {
		fmt.Println("SaveInventoryData:", err)
		return false
	}

	policy := aerospike.NewWritePolicy(1, inventoryExpire)

	name := ""
	if request.Channel == m.ChannelSite {
		name = request.Site
	} else {
		name = request.App
	}

	bins := []*aerospike.Bin{
		aerospike.NewBin("id", request.Targeting.Inventory),
		aerospike.NewBin("channel", request.Channel.String()),
		aerospike.NewBin("name", name),
		aerospike.NewBin("status", 1),
		aerospike.NewBin("saved", 0),
	}

	err = client.PutBins(policy, key, bins...)

	if err != nil {
		fmt.Println("SaveInventoryData:", err)
		return false
	}

	return true
}

func UpdateInventoryData2(inventory string, country, device int, data InventoryData2) bool {

	ctm := time.Now()

	inventoryKey := fmt.Sprintf("%04d-%02d-%02d_%s_%d_%d", ctm.Year(), ctm.Month(), ctm.Day(), inventory, country, device)

	key, err := aerospike.NewKey("bidder", "inventory_stats", inventoryKey)
	if err != nil {
		fmt.Println("UpdateInventoryData:", err)
		return false
	}

	ops := []*aerospike.Operation{
		aerospike.AddOp(aerospike.NewBin("aimp", data.aimp)),
		aerospike.AddOp(aerospike.NewBin("req", data.req)),
		aerospike.AddOp(aerospike.NewBin("fl", int(data.fl))),
	}

	policy := aerospike.NewWritePolicy(1, inventoryExpire)

	_, err = client.Operate(policy, key, ops...)

	if err != nil {
		fmt.Println("UpdateInventoryData:", err)
		return false
	}

	return true
}

func UpdateInventoryData(inventory string, country, device int, data InventoryData) bool {

	ctm := time.Now()

	inventoryKey := fmt.Sprintf("%04d-%02d-%02d_%s_%d_%d", ctm.Year(), ctm.Month(), ctm.Day(), inventory, country, device)

	key, err := aerospike.NewKey("bidder", "inventory_stats", inventoryKey)
	if err != nil {
		fmt.Println("UpdateInventoryData:", err)
		return false
	}

	ops := []*aerospike.Operation{
		aerospike.AddOp(aerospike.NewBin("imp", data.imp)),
		aerospike.AddOp(aerospike.NewBin("clk", data.clk)),
		aerospike.AddOp(aerospike.NewBin("bp", data.bp)),
		aerospike.AddOp(aerospike.NewBin("wp", data.wp)),
		aerospike.AddOp(aerospike.NewBin("rev", data.rev)),
		aerospike.AddOp(aerospike.NewBin("spd", data.spd)),
	}

	policy := aerospike.NewWritePolicy(1, inventoryExpire)

	_, err = client.Operate(policy, key, ops...)

	if err != nil {
		fmt.Println("UpdateInventoryData:", err)
		return false
	}

	return true
}

func BidSave(cid int, inv string) {

	op := aerospike.AddOp(aerospike.NewBin("bid", 1))
	ctm := time.Now()
	policy := aerospike.NewWritePolicy(1, 60*60*24*8)

	keyStr := fmt.Sprintf("%04d-%02d-%02d_%s_%d", ctm.Year(), ctm.Month(), ctm.Day(), inv, cid)

	key, err := aerospike.NewKey("bidder", "inventory_camp", keyStr)
	if err != nil {
		fmt.Println("FeaturesSave:", err)
		return
	}

	record, err := client.Operate(policy, key, op)
	if err != nil {
		fmt.Println(keyStr, record, err)
	}
}

func FeaturesSave(uid, cat []byte, cid int, inv string, userId []byte, country, device int, bp, wp, rev, spd float64, action string) {

	imp, clk := 0, 0

	switch action {
	case "impression":
		imp = 1
	case "click":
		clk = 1
	}

	UpdateInventoryData(inv, country, device, InventoryData{
		imp: imp,
		clk: clk,
		bp:  int(bp * 1000),
		wp:  int(wp * 1000),
		rev: int(rev * 1000),
		spd: int(spd * 1000),
	})

	/*map[string]interface{}{
		"imp": imp,
		"clk": clk,
		"bp":  int(bp * 1000),
		"wp":  int(wp * 1000),
		"rev": int(rev * 1000),
		"spd": int(spd * 1000),
	}*/

	ctm := time.Now()

	keyStr := fmt.Sprintf("%04d-%02d-%02d_%s_%d", ctm.Year(), ctm.Month(), ctm.Day(), uid, cid)

	key, err := aerospike.NewKey("bidder", "frequency", keyStr)
	if err != nil {
		fmt.Println("FeaturesSave:", err)
		return
	}

	ops := []*aerospike.Operation{
		aerospike.AddOp(aerospike.NewBin("imp", imp)),
		aerospike.AddOp(aerospike.NewBin("clk", clk)),
	}

	policy := aerospike.NewWritePolicy(1, 60*60*24)

	record, err := client.Operate(policy, key, ops...)
	if err != nil {
		fmt.Println(keyStr, record, err)
	}

	// ==

	keyStr = fmt.Sprintf("%04d-%02d-%02d_%s_%s", ctm.Year(), ctm.Month(), ctm.Day(), inv, cat)

	key, err = aerospike.NewKey("bidder", "inventory_cat", keyStr)
	if err != nil {
		fmt.Println("FeaturesSave:", err)
		return
	}

	policy = aerospike.NewWritePolicy(1, 60*60*24*8)

	record, err = client.Operate(policy, key, ops...)
	if err != nil {
		fmt.Println(keyStr, record, err)
	}

	// ==

	keyStr = fmt.Sprintf("%04d-%02d-%02d_%s_%d", ctm.Year(), ctm.Month(), ctm.Day(), inv, cid)

	key, err = aerospike.NewKey("bidder", "inventory_camp", keyStr)
	if err != nil {
		fmt.Println("FeaturesSave:", err)
		return
	}

	record, err = client.Operate(policy, key, ops...)
	if err != nil {
		fmt.Println(keyStr, record, err)
	}

	// ==

	daytype := 'w'

	switch ctm.Weekday() {
	case time.Saturday, time.Sunday:
		daytype = 'h'
	}

	hourslice := ctm.Hour() / 4

	keys := make([]string, 0, 10)

	// - inventoryid_campaignid_countryid
	keyPrefix := fmt.Sprintf("%04d-%02d-%02d_%s_%d_%d", ctm.Year(), ctm.Month(), ctm.Day(), inv, cid, country)

	keys = append(keys, keyPrefix)

	// - inventoryid_campaignid_countryid_daytype_hourslice_userip+ua
	keys = append(keys, fmt.Sprintf("%s_%c_%d_%s", keyPrefix, daytype, hourslice, userId))

	// - inventoryid_campaignid_countryid_daytype_userip+ua
	keys = append(keys, fmt.Sprintf("%s_%c_%s", keyPrefix, daytype, userId))

	// - inventoryid_campaignid_countryid_daytype_hourslice
	keys = append(keys, fmt.Sprintf("%s_%c_%d", keyPrefix, daytype, hourslice))

	// - inventoryid_campaignid_countryid_daytype
	keys = append(keys, fmt.Sprintf("%s_%d", keyPrefix, daytype))

	// - inventoryid_campaignid_countryid_userip+ua
	keys = append(keys, fmt.Sprintf("%s_%s", keyPrefix, userId))

	// ==

	// - inventoryid_campaignid
	keyPrefix = fmt.Sprintf("%04d-%02d-%02d_%s_%d", ctm.Year(), ctm.Month(), ctm.Day(), inv, cid)
	keys = append(keys, keyPrefix)

	// - inventoryid_campaignid_userip+ua
	keys = append(keys, fmt.Sprintf("%s_%s", keyPrefix, userId))

	// ==

	for n := range keys {

		if key, err = aerospike.NewKey("bidder", "inventory_campaign", keys[n]); err == nil {

			if record, err = client.Operate(policy, key, ops...); err != nil {

				fmt.Println(keys[n], record, err)
			}
		}
	}
}
