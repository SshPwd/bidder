package main

import (
	"flag"
	"golang_bidder/bidder"
	"golang_bidder/browscap"
	"golang_bidder/config"
	"golang_bidder/db"
	"golang_bidder/ip2location"
	"golang_bidder/logs"
	// "golang_bidder/utils"
	"runtime"
	"runtime/debug"
	"sync"
)

var (
	build, date string
	production  string = "true"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	debug.SetGCPercent(20)

	configNamePtr := flag.String("c", config.DefaultConfigName, "config file")
	flag.Parse()

	config.Load(*configNamePtr)
	conf := config.Get()

	logs.Init(conf.LogFile)

	logs.Debug("START")
	logs.Debug("BUILD:", build, date)

	db.InitAerospike()
	db.InitMongodb()

	if production == "true" {

		var wg sync.WaitGroup
		wg.Add(4)

		go func() {
			bidder.Geoip = ip2location.Load(conf.GeoipDb)
			wg.Done()
		}()

		go func() {
			bidder.Browsers = browscap.Load(conf.UserAgentsDb)
			wg.Done()
		}()

		go func() {
			bidder.Datacenters = ip2location.LoadDatacenters(conf.DatacentersDb)
			wg.Done()
		}()

		go func() {
			db.LoadBidderData()
			wg.Done()
		}()

		wg.Wait()
	}

	bidder.BuildDate = date
	bidder.BuildName = build

	bidder.Serve()
}
