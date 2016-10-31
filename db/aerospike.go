package db

import (
	"fmt"
	"golang_bidder/config"
	"golang_bidder/logs"
	m "golang_bidder/model"
	"golang_bidder/utils"
	"math"
	"os"
	"sync"
	"time"

	aerospike "github.com/aerospike/aerospike-client-go"
)

const (
	C0 = "\x1b[30;1m"
	CR = "\x1b[31;1m"
	CG = "\x1b[32;1m"
	CY = "\x1b[33;1m"
	CB = "\x1b[34;1m"
	CM = "\x1b[35;1m"
	CC = "\x1b[36;1m"
	CW = "\x1b[37;1m"
	CN = "\x1b[0m"

	namespace = "test"
	set       = ""
)

var (
	client *aerospike.Client
)

func InitAerospike() {

	conf := config.Get()

	var err error

	clientPolicy := aerospike.NewClientPolicy()
	clientPolicy.Timeout = 30 * time.Second

	client, err = aerospike.NewClientWithPolicy(clientPolicy, conf.Aerospike.Host, conf.Aerospike.Port)
	if err != nil {
		logs.Critical(fmt.Sprintf("Cannot connect to Aerospike server: %s:%d (%s)",
			conf.Aerospike.Host, conf.Aerospike.Port, err.Error()))
		os.Exit(1)
		return
	}

	fmt.Println(client)
}

// ==

func GetInventoryData(request *m.Request) m.InventoryData {

	lastSeven := utils.GetLastWeek()

	keys := make([]*aerospike.Key, 0, 8)

	key, err := aerospike.NewKey("bidder", "inventory", request.Targeting.Inventory)
	if err == nil {
		// fmt.Println(key)
		keys = append(keys, key)
	}

	countryId := MapCountry.Get(request.Targeting.Country)

	keyFmt := fmt.Sprintf("%%s_%s_%d_%d", request.Targeting.Inventory, countryId, request.Targeting.PlatformId)

	for _, date := range lastSeven {

		keyStr := fmt.Sprintf(keyFmt, date)
		// fmt.Println("inventory_stats:", keyStr)

		key, err := aerospike.NewKey("bidder", "inventory_stats", keyStr)
		if err == nil {
			keys = append(keys, key)
		}
	}

	if results, err := client.BatchGet(nil, keys); err == nil {

		// if request.Targeting.Inventory == "oVcCP0ZxN4SLEI4w5dhDaQ" {
		// 	for n := range results {
		// 		if results[n] != nil {
		// 			fmt.Println(CM+"<<<", keys[n].Value(), CN, CG+">>>", results[n].Bins, CN)
		// 		}
		// 	}
		// }

		inventoryData := m.InventoryData{}

		if results[0] == nil {

			SaveInventoryData(request)

			UpdateInventoryData2(request.Targeting.Inventory, countryId, int(request.Targeting.PlatformId), InventoryData2{
				aimp: request.PlacementCount,
				fl:   request.PlacementFloors,
				req:  1,
			})

		} else {

			if readFieldInt(results[0].Bins, "status") == 0 {
				return inventoryData
			}

			UpdateInventoryData2(request.Targeting.Inventory, countryId, int(request.Targeting.PlatformId), InventoryData2{
				aimp: request.PlacementCount,
				fl:   request.PlacementFloors,
				req:  1,
			})

			for n := 1; n < len(results); n++ {

				if results[n] != nil {

					inventoryData.Aimp += readFieldInt(results[n].Bins, "aimp")
					inventoryData.Fl += readFieldInt(results[n].Bins, "fl")
					inventoryData.Req += readFieldInt(results[n].Bins, "req")
					inventoryData.Imp += readFieldInt(results[n].Bins, "imp")
					inventoryData.Clk += readFieldInt(results[n].Bins, "clk")
					inventoryData.Bp += readFieldInt(results[n].Bins, "bp")
					inventoryData.Wp += readFieldInt(results[n].Bins, "wp")
					inventoryData.Rev += readFieldInt(results[n].Bins, "rev")
					inventoryData.Spd += readFieldInt(results[n].Bins, "spd")
				}
			}
		}

		return inventoryData
	}

	// TODO

	// data := map[string]int{
	// 	"imp": 1,
	// 	"clk": 0,
	// 	"bp":  0,
	// 	"wp":  0,
	// 	"rev": 0,
	// 	"spd": 0,
	// }

	// UpdateInventory(request.Targeting.Inventory, request.Targeting.Country, platformName, data)

	return m.InventoryData{}
}

// ==

func readFieldInt(bins aerospike.BinMap, fieldName string) int {

	if bins != nil {

		if fieldData, ok := bins[fieldName]; ok {

			value, _ := fieldData.(int)
			return value
		}
	}
	return 0
}

// ==

type Counter struct {
	clk, imp, bid int
}

type CategoriInfo struct {
	Name    string
	Counter Counter
}

func GetFeatures(request *m.Request, campaigns []m.Campaign) []m.Campaign {

	if len(campaigns) == 0 {
		return campaigns
	}

	keys := make([]*aerospike.Key, 0, len(campaigns))
	camp_inv := map[int]Counter{}
	lastSeven := utils.GetLastWeek()
	categories := []CategoriInfo{}
	campaignCIDs := make([]int, 0, len(campaigns))

	ctm := time.Now()
	keyFmt := fmt.Sprintf("%04d-%02d-%02d_%s_%%d", ctm.Year(), ctm.Month(), ctm.Day(), request.GlobalUser)

	for i := range campaigns {

		campaignCIDs = append(campaignCIDs, campaigns[i].CampaignId)

		keyStr := fmt.Sprintf(keyFmt, campaigns[i].CampaignId)

		// fmt.Println("frequency:", keyStr)

		key, err := aerospike.NewKey("bidder", "frequency", keyStr)
		if err == nil {
			keys = append(keys, key)
		}

		if campaigns[i].Category != "" {
			// categories[campaigns[i].Category] = Counter{}
			categories = append(categories, CategoriInfo{Name: campaigns[i].Category})
		}
	}

	for _, cat := range categories {

		keyFmt = fmt.Sprintf("%%s_%s_%s", request.Targeting.Inventory, cat.Name)

		for _, date := range lastSeven {

			keyStr := fmt.Sprintf(keyFmt, date)

			// fmt.Println("inventory_cat:", keyStr)

			key, err := aerospike.NewKey("bidder", "inventory_cat", keyStr)
			if err == nil {
				keys = append(keys, key)
			}
		}
	}

	camp_inv_start := len(keys)

	for i := range campaigns {

		camp_inv[campaigns[i].CampaignId] = Counter{}
		keyFmt = fmt.Sprintf("%%s_%s_%d", request.Targeting.Inventory, campaigns[i].CampaignId)

		for _, date := range lastSeven {

			keyStr := fmt.Sprintf(keyFmt, date)

			// fmt.Println("inventory_camp:", keyStr)

			key, err := aerospike.NewKey("bidder", "inventory_camp", keyStr)
			if err == nil {
				keys = append(keys, key)
			}
		}
	}

	// fmt.Println("keys:", len(keys))

	if results, err := client.BatchGet(nil, keys); err == nil {

		// ==

		// if request.Targeting.Inventory == "oVcCP0ZxN4SLEI4w5dhDaQ" {

		// 	for n := range results {
		// 		if results[n] != nil {
		// 			fmt.Println(CM+"<<<", keys[n].Value(), CN)
		// 			fmt.Println(CG+">>>", results[n].Bins, CN)
		// 		}
		// 	}
		// }

		// ==

		campaignLength := len(campaigns)

		skipIndex := map[int]struct{}{}

		for i := range campaigns {

			// fmt.Println("GetFeatures:", i, results[i])

			if results[i] != nil {

				if imp := readFieldInt(results[i].Bins, "imp"); imp >= 10 {
					skipIndex[i] = struct{}{}
				} else {
					campaigns[i].ImpFreq = imp
				}
			}
		}

		if len(skipIndex) > 0 {

			tmp := make([]m.Campaign, 0, len(campaigns)-len(skipIndex))

			for i := range campaigns {
				if _, ok := skipIndex[i]; !ok {
					tmp = append(tmp, campaigns[i])
				}
			}

			campaigns = tmp

			for i := range skipIndex {
				delete(skipIndex, i)
			}

			if len(campaigns) == 0 {
				return campaigns
			}
		}

		categoriesTotal := Counter{}

		for n := range categories {
			for i := 0; i < 7; i++ {

				if results[campaignLength] != nil {

					if imp := readFieldInt(results[campaignLength].Bins, "imp"); imp > 0 {
						categoriesTotal.imp += imp
						categories[n].Counter.imp += imp
					}

					if clk := readFieldInt(results[campaignLength].Bins, "clk"); clk > 0 {
						categoriesTotal.clk += clk
						categories[n].Counter.clk += clk
					}
				}
				campaignLength++
			}
		}

		for x, rotation, cid := camp_inv_start, 0, 0; x < len(results); x++ {

			if results[x] != nil {

				campaignId := campaignCIDs[cid]
				counter := camp_inv[campaignId]

				counter.imp += readFieldInt(results[x].Bins, "imp")
				counter.clk += readFieldInt(results[x].Bins, "clk")
				counter.bid += readFieldInt(results[x].Bins, "bid")

				camp_inv[campaignId] = counter
			}

			rotation++
			if rotation == 7 {
				cid++
				rotation = 0
			}
		}

		inventoryCamp := request.Targeting.Inventory

		request.FeaturesCtr = calculateFeatures(inventoryCamp, camp_inv, categories, categoriesTotal, request.InventoryData)

		return campaigns

	} else {

		fmt.Println("BatchGet:", err)
	}

	return []m.Campaign{}
}

// ==

func calculateFeatures(inventoryCamp string, camp_inv map[int]Counter, categories []CategoriInfo, categoriesTotal Counter,
	inventory m.InventoryData) m.FeaturesCtr {

	// fmt.Println("calculateFeatures(")
	// fmt.Printf("  camp_inv: %#v\n", camp_inv)
	// fmt.Printf("  categories: %#v\n", categories)
	// fmt.Printf("  categoriesTotal: %#v\n", categoriesTotal)
	// fmt.Printf("  inventory: %#v\n", inventory)
	// fmt.Println(")")

	features := m.FeaturesCtr{
		CampInv:           map[int]float64{},
		Categories:        map[string]float64{},
		InventoryCtr:      calculateSmartCTR("", inventory.Imp, inventory.Clk, 0, 1000, 1),
		CategoriesTotal:   calculateSmartCTR("", categoriesTotal.imp, categoriesTotal.clk, 0, 1000, 1),
		CategoriesDefault: defCTR,
	}

	var delta float64

	if inventory.Imp > 1000 && inventory.Rev > 0 {
		delta = float64(inventory.Spd) / float64(inventory.Rev)
	} else {
		delta = 1.0
	}

	adjust := features.InventoryCtr / features.CategoriesTotal * math.Min(1.5, delta)

	for n := range categories {

		features.Categories[categories[n].Name] = calculateSmartCTR("", categories[n].Counter.imp,
			categories[n].Counter.clk, 0, 1000, adjust)
	}

	for cid, counter := range camp_inv {
		if counter.imp > 100 || counter.clk > 0 {
			features.CampInv[cid] = calculateSmartCTR(inventoryCamp, counter.imp, counter.clk, counter.bid, 0, 1)
		}
	}

	// fmt.Printf(CW+">>> features: %#v\n\n"+CN, features)

	return features
}

// ==

const (
	defCTR      float64 = 1
	leastCTR    float64 = 0.01
	greatestCTR float64 = 5
)

func calculateSmartCTR(key string, impressions, clicks, bid, peak int, offset float64) float64 {

	switch {
	case impressions < peak:

		if key != "" {
			return Corrector(key, defCTR, impressions, clicks, bid)
		}

		return defCTR

	case clicks == 0:
		return leastCTR

	case clicks > impressions:
		return greatestCTR
	}

	if offset <= 0 {
		offset = 1
	}

	clicks0, impressions0, offset0 := float64(clicks), float64(impressions), float64(offset)

	smartCTR := clicks0 / (impressions0 - (impressions0 / (clicks0 + 1.0))) * 100.0 * offset0

	return math.Min(math.Max(smartCTR, leastCTR), greatestCTR)
}

// ==

var (
	allCtx  = map[string]CtrCtx{}
	lockCtx sync.Mutex
)

func Corrector(key string, ctr float64, imp, clk, bid int) float64 {

	if imp > 1000 || clk > 50 {
		return ctr
	}

	lockCtx.Lock()

	ctx, ok := allCtx[key]
	if !ok {
		ctx.Time = time.Now().Unix()
	}

	ctr = ctx.Process(ctr, imp, bid)

	allCtx[key] = ctx

	lockCtx.Unlock()
	return ctr
}

type CtrCtx struct {
	Time     int64
	Delta    float64
	Progress int
	Step     int
	Deep     int
}

func (ctx *CtrCtx) CalcDelta(ctr float64) (delta float64) {

	const minstep = 0.01
	k := 2.0

	for i := 0; i < ctx.Deep; i++ {

		delta += ctr / k
		k *= 2
	}

	if delta > 0.95 {
		delta = 0.95
	}

	delta = -delta
	return
}

func (ctx *CtrCtx) Process(ctr float64, imp, bid int) float64 {

	if imp < 1000 {

		step := bid / 100

		if ctx.Step != step {

			ctx.Step = step

			if ctx.Progress == 0 {
				ctx.Progress = imp
			}

			if ctx.Progress < imp {
				if ctx.Deep < 12 {
					ctx.Deep++
				}
			} else if ctx.Progress >= imp {
				if ctx.Deep > 1 {
					ctx.Deep--
				}
			}

			ctx.Delta = ctx.CalcDelta(ctr)

			logs.Debug(fmt.Sprintf("Process(ctr:%.3f, deep:%d, bid:%d): step:%d, imp:%d, delta:%.3f, ctr:%.3f",
				ctr, ctx.Deep, bid, step, ctx.Progress-imp, ctx.Delta, ctr+ctr*ctx.Delta))

			ctx.Progress = imp
		}

		ctr += ctr * ctx.Delta

	} else {

		ctx.Delta = 0
	}
	return ctr
}
