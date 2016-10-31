package ip2location

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"golang_bidder/logs"
	"golang_bidder/utils"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type (
	IP2Location struct {
		Records        []Record
		Types          map[string]uint16
		Countries      map[string]uint16
		Cities         map[string]int
		Regions        map[string]int
		HostingIndex   []bool
		TypesIndex     []string
		CountriesIndex []string
		CitiesIndex    []string
		RegionsIndex   []string
		TestMode       bool
	}

	Record struct {
		LeftIp    uint32
		RightIp   uint32
		RegionId  int
		CityId    int
		CountryId uint16
		TypeId    uint16
		IsHosting bool
	}

	Result struct {
		Country   string
		City      string
		Region    string
		Type      string
		IsHosting bool
		Ok        bool
	}
)

var (
	DefaultIP2Location = IP2Location{TestMode: true}
	ErrNotFound        = errors.New("NotFound")
	hostingWord        = []string{"cdn", "dch", "ses", "rsv"}
)

func Load(filename string) (ip IP2Location) {

	defer utils.WTimeStart().Stop()

	if filename == "" {
		logs.Critical("ip2location db not set")
		os.Exit(1)
	}

	file, err := os.Open(filename)
	if err != nil {
		logs.Critical(fmt.Sprintf("Cannot open ip2location db, file \"%s\" not exist", filename))
		os.Exit(1)
	}
	defer file.Close()

	fmt.Println(filename)

	var reader io.ReadCloser

	if filepath.Ext(filename) == ".gz" {
		reader, err = gzip.NewReader(file)
		defer reader.Close()
	} else {
		reader = file
	}

	scaner := bufio.NewScanner(reader)
	scaner.Split(bufio.ScanLines)

	ip.Types = map[string]uint16{"": 0}
	ip.Cities = map[string]int{"": 0}
	ip.Regions = map[string]int{"": 0}
	ip.Countries = map[string]uint16{"": 0}

	// 13125470
	ip.Records = make([]Record, 0, 15000000)

	for scaner.Scan() {

		line := scaner.Text()
		list := strings.Split(line, "\",\"")

		leftIp := unquote(list[0])
		rightIp := unquote(list[1])
		country := unquote(list[2])
		region := unquote(list[4])
		city := unquote(list[5])
		typ := unquote(list[13])

		countryId := ip.GetCountryId(country)
		regionId := ip.GetRegionId(region)
		cityId := ip.GetCityId(city)
		typeId := ip.GetTypeId(typ)

		record := Record{
			LeftIp:    utils.ToUint32(leftIp),
			RightIp:   utils.ToUint32(rightIp),
			RegionId:  regionId,
			CityId:    cityId,
			CountryId: countryId,
			TypeId:    typeId,
		}

		ip.Records = append(ip.Records, record)
	}

	ip.updateIndex()

	// fmt.Printf("%#v\n", ip.Find("8.8.8.8"))
	return
}

func (this *IP2Location) Find(ipStr string) (result Result) {

	if this.TestMode {
		result.Type = "com"
		result.City = "kiev"
		result.Region = "kyyiv"
		result.Country = "ua"
		result.IsHosting = false
		result.Ok = true
		return
	}

	ip := utils.IpToUint32(ipStr)
	length := len(this.Records)

	index := sort.Search(length, func(i int) bool {

		return this.Records[i].RightIp >= ip
	})

	if index < length && this.Records[index].LeftIp <= ip && ip <= this.Records[index].RightIp {

		record := this.Records[index]

		result.Type = this.TypesIndex[record.TypeId]
		result.City = this.CitiesIndex[record.CityId]
		result.Region = this.RegionsIndex[record.RegionId]
		result.Country = this.CountriesIndex[record.CountryId]
		result.IsHosting = this.HostingIndex[record.TypeId]
		result.Ok = true
	}
	return
}

func isHosting(typ string) bool {
	for n := range hostingWord {
		if pos := strings.Index(typ, hostingWord[n]); pos != -1 {
			return true
		}
	}
	return false
}

func (this *IP2Location) updateIndex() {

	this.TypesIndex = make([]string, len(this.Types))
	this.CitiesIndex = make([]string, len(this.Cities))
	this.RegionsIndex = make([]string, len(this.Regions))
	this.HostingIndex = make([]bool, len(this.Types))
	this.CountriesIndex = make([]string, len(this.Countries))

	for str, index := range this.Cities {
		this.CitiesIndex[index] = strings.ToLower(str)
		delete(this.Cities, str)
	}

	for str, index := range this.Types {
		str = strings.ToLower(str)
		this.TypesIndex[index] = str
		this.HostingIndex[index] = isHosting(str)
		delete(this.Types, str)
	}

	for str, index := range this.Regions {
		this.RegionsIndex[index] = strings.ToLower(str)
		delete(this.Regions, str)
	}

	for str, index := range this.Countries {
		this.CountriesIndex[index] = strings.ToLower(str)
		delete(this.Countries, str)
	}
}

func (this *IP2Location) GetCountryId(str string) uint16 {
	if n, ok := this.Countries[str]; ok {
		return n
	}
	n := uint16(len(this.Countries))
	this.Countries[str] = n
	return n
}

func (this *IP2Location) GetTypeId(str string) uint16 {
	if n, ok := this.Types[str]; ok {
		return n
	}
	n := uint16(len(this.Types))
	this.Types[str] = n
	return n
}

func (this *IP2Location) GetRegionId(str string) int {
	if n, ok := this.Regions[str]; ok {
		return n
	}
	n := len(this.Regions)
	this.Regions[str] = n
	return n
}

func (this *IP2Location) GetCityId(str string) int {
	if n, ok := this.Cities[str]; ok {
		return n
	}
	n := len(this.Cities)
	this.Cities[str] = n
	return n
}

func unquote(str string) string {
	str = strings.TrimSuffix(str, `"`)
	str = strings.TrimPrefix(str, `"`)
	return str
}
