package config

import (
	"encoding/json"
	"fmt"
	"golang_bidder/logs"
	"io/ioutil"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type (
	MysqlConfig struct {
		DbName   string `json:"dbname"`
		Login    string `json:"login"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     string `json:"port"`
	}

	AerospikeConfig struct {
		Host    string `json:"host"`
		Port    int    `json:"port"`
		Timeout int64  `json:"timeout"`
	}

	MongodbConfig struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Host     string `json:"host"`
		DbName   string `json:"db_name"`
	}

	Event struct {
		BaseUrl string `json:"base_url"`
	}

	Cdn struct {
		ImgBaseUrl string `json:"img_base_url"`
	}

	ConfigData struct {
		Bind          []string        `json:"bind"`
		GeoipDb       string          `json:"geoip_db"`
		DatacentersDb string          `json:"datacenters_db"`
		UserAgentsDb  string          `json:"user_agents_db"`
		BidderDataUrl string          `json:"bidder_data_url"`
		LogFile       string          `json:"log_file"`
		PidFile       string          `json:"pid_file"`
		DebugMode     int             `json:"debug_mode"`
		MysqlMaster   MysqlConfig     `json:"mysql"`
		Aerospike     AerospikeConfig `json:"aerospike"`
		Mongodb       MongodbConfig   `json:"mongodb"`
		Event         Event           `json:"event"`
		Cdn           Cdn             `json:"cdn"`
	}
)

const (
	DefaultConfigName = "bidder.conf"
	DefaultPidName    = "bidder.pid"
)

var (
	configName       string = DefaultConfigName
	globalConfigData atomic.Value
	reloadCbMutex    sync.Mutex
	reloadCbList     []func(*ConfigData)
	once             sync.Once
)

func init() {
	globalConfigData.Store(&ConfigData{})
}

func load(first bool) {

	config := new(ConfigData)

	filename := configName

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		logs.Critical(fmt.Sprintf("Config file \"%s\" not found", filename))
		if first {
			// os.Exit(1)
		} else {
			return
		}
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		logs.Critical(fmt.Sprintf("Fail to parse config file \"%s\"", filename))
		if first {
			// os.Exit(1)
		} else {
			return
		}
	}

	if len(config.Bind) == 0 {
		logs.Critical("Config param \"bind\" not set")
		if first {
			// os.Exit(1)
		} else {
			return
		}
	}

	// emptyMysqlConfig := MysqlConfig{}
	emptyAerospikeConfig := AerospikeConfig{}

	// if config.MysqlMaster == emptyMysqlConfig {
	// 	logs.Critical("Config param \"mysql\" not set")
	// 	if first {
	// 		// os.Exit(1)
	// 	} else {
	// 		return
	// 	}
	// }

	if config.Aerospike == emptyAerospikeConfig {
		logs.Critical("Config param \"aerospike\" not set")
		if first {
			// os.Exit(1)
		} else {
			return
		}
	}

	if config.UserAgentsDb == "" {
		logs.Critical("Config param \"user_agents_db\" not set")
	}

	if config.PidFile == "" {
		logs.Critical("Config param \"pid_file\" not set")
		config.PidFile = DefaultPidName
	}

	globalConfigData.Store(config)

	fmt.Println("loaded config:", filename)

	go once.Do(checkUpdate)
}

func checkUpdate() {

	info, err := os.Stat(configName)
	if err != nil {
		logs.Critical(err.Error())
		return
	}

	startTime := info.ModTime()

	for {

		time.Sleep(1 * time.Minute)

		if info, err = os.Stat(configName); err == nil {

			if info.ModTime().Sub(startTime) > 0 {

				reloadConfig()
				startTime = info.ModTime()
			}
		}
	}
}

func reloadConfig() {

	load(false)

	configData := globalConfigData.Load().(*ConfigData)

	reloadCbMutex.Lock()

	for _, cb := range reloadCbList {
		cb(configData)
	}

	reloadCbMutex.Unlock()
}

func OnReloadCb(cb func(*ConfigData)) {

	reloadCbMutex.Lock()
	reloadCbList = append(reloadCbList, cb)
	reloadCbMutex.Unlock()
}

func Load(filename string) {

	if filename != "" {
		configName = filename
	}

	load(true)
}

func Get() *ConfigData {

	return globalConfigData.Load().(*ConfigData)
}
