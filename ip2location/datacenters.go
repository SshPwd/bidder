package ip2location

import (
	"bufio"
	"compress/gzip"
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
	Datacenters struct {
		Records []DcRecord
	}

	DcRecord struct {
		LeftIp  uint32
		RightIp uint32
	}
)

var (
	DefaultDatacenters = Datacenters{}
)

// https://raw.githubusercontent.com/client9/ipcat/master/datacenters.csv

func LoadDatacenters(filename string) (datacenters Datacenters) {

	if filename == "" {
		logs.Critical("datacenters db not set")
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		logs.Critical(fmt.Sprintf("Cannot open datacenters db, file \"%s\" not exist", filename))
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

	for scaner.Scan() {

		line := scaner.Text()
		list := strings.Split(line, ",")

		if len(list) == 4 {

			record := DcRecord{
				LeftIp:  utils.IpToUint32(list[0]),
				RightIp: utils.IpToUint32(list[1]),
			}

			datacenters.Records = append(datacenters.Records, record)
		}
	}

	return
}

func (this *Datacenters) Find(ipStr string) bool {

	ip := utils.IpToUint32(ipStr)
	length := len(this.Records)

	index := sort.Search(length, func(i int) bool {

		return this.Records[i].RightIp >= ip
	})

	if index < length && this.Records[index].LeftIp <= ip && ip <= this.Records[index].RightIp {

		return true
	}

	return false
}
