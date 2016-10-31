package browscap

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"github.com/digitalcrab/browscap_go"
	"golang_bidder/logs"
	"golang_bidder/utils"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

const (
	DESKTOP = 1 // default
	TABLET  = 2
	PHONE   = 3

	WINDOWS  = 1 // default
	MACOSX   = 2
	IOS      = 3
	ANDROID  = 4
	WINPHONE = 5
)

type (
	Browscap struct {
		expressionTree *browscap_go.ExpressionTree
		deviceInfo     map[string]Record
		TestMode       bool
	}

	Record struct {
		OsId       uint8
		PlatformId uint8 // desktop = 0, tablet = 1, phone = 2
		Crawler    bool
		Ok         bool
	}
)

var (
	DefaultBrowscap = Browscap{
		expressionTree: browscap_go.NewExpressionTree(),
		deviceInfo:     map[string]Record{},
		TestMode:       true,
	}

	PlatformName = map[uint8]string{
		DESKTOP: "desktop",
		TABLET:  "tablet",
		PHONE:   "phone",
	}

	OsName = map[uint8]string{
		WINDOWS:  "windows",
		MACOSX:   "macosx",
		ANDROID:  "android",
		IOS:      "ios",
		WINPHONE: "winphone",
	}
)

// desktop: MacOSX, Windows, Other
// tablet:  IOS, Android, Other
// phone:   IOS, Android, Windows Phone, Other

// "tv device"
// "console"

func Load(filename string) (bc Browscap) {

	defer utils.WTimeStart().Stop()

	bc.expressionTree = browscap_go.NewExpressionTree()
	bc.deviceInfo = map[string]Record{}

	file, err := os.Open(filename)
	if err != nil {
		logs.Critical(fmt.Sprintf("Cannot open browscap db, file \"%s\" not exist", filename))
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

	if scaner.Scan() { // version

		if scaner.Scan() { // date

			if scaner.Scan() {

				columnStr := scaner.Text()
				columnName := map[string]int{}

				for n, str := range strings.Split(columnStr, ",") {
					columnName[unquote(str)] = n
				}

				if scaner.Scan() {

					defaultStr := scaner.Text()
					defaultValue := map[int]string{}

					for n, str := range strings.Split(defaultStr, ",") {
						defaultValue[n] = unquote(str)
					}

					line := make([]string, 0, 64)
					platform, browserType, deviceType := "", "", ""

					indexCrawler := columnName["Crawler"]

					indexPlatform := columnName["Platform"]
					indexDeviceType := columnName["Device_Type"]
					indexBrowserType := columnName["Browser_Type"]

					defaultPlatform := strings.ToLower(defaultValue[indexPlatform])
					defaultDeviceType := strings.ToLower(defaultValue[indexDeviceType])
					defaultBrowserType := strings.ToLower(defaultValue[indexBrowserType])

					indexIsTablet := columnName["isTablet"]
					// indexIsMobileDevice := columnName["isMobileDevice"]

					columnByIndex := func(index int, defaultValue string) string {
						value := line[index]
						if value == "" {
							return defaultValue
						}
						return strings.ToLower(value)
					}

					for scaner.Scan() {

						record := Record{}

						for _, str := range strings.Split(scaner.Text(), "\",\"") {
							line = append(line, unquote(str))
						}

						platform = columnByIndex(indexPlatform, defaultPlatform)
						deviceType = columnByIndex(indexDeviceType, defaultDeviceType)
						browserType = columnByIndex(indexBrowserType, defaultBrowserType)

						switch {
						case strings.Index(platform, "ios") != -1:
							record.OsId = IOS
						case strings.Index(platform, "android") != -1:
							record.OsId = ANDROID
						case strings.Index(platform, "macosx") != -1:
							record.OsId = MACOSX
						case strings.Index(platform, "winphone") != -1:
							record.OsId = WINPHONE
						case strings.Index(platform, "win") != -1:
							record.OsId = WINDOWS
						}

						if deviceType == "desktop" {

							record.PlatformId = DESKTOP

						} else if line[indexIsTablet] == "true" ||
							deviceType == "tablet" {

							record.PlatformId = TABLET

						} else if deviceType == "mobile phone" ||
							deviceType == "mobile device" {

							record.PlatformId = PHONE
						}

						if line[indexCrawler] == "true" ||
							browserType == "bot/crawler" {

							record.Crawler = true
						}

						expression := line[0]

						bc.deviceInfo[expression] = record
						bc.expressionTree.Add(expression)

						line = line[:0]
					}
				}
			}
		}
	}
	return
}

func (bc *Browscap) Find(userAgent []byte) (record Record) {

	if bc.TestMode {
		record.Ok = true
		record.OsId = WINDOWS
		record.PlatformId = DESKTOP
		return
	}

	agent := mapBytes(unicode.ToLower, userAgent)
	defer bytesPool.Put(agent)

	name := bc.expressionTree.Find(agent)
	if name == "" {
		return
	}

	record = bc.deviceInfo[name]
	record.Ok = true
	return
}

func (bc *Browscap) FindStr(userAgent string) (record Record) {

	if bc.TestMode {
		record.Ok = true
		record.OsId = WINDOWS
		record.PlatformId = DESKTOP
		return
	}

	agent := mapToBytes(unicode.ToLower, userAgent)
	defer bytesPool.Put(agent)

	name := bc.expressionTree.Find(agent)
	if name == "" {
		return
	}

	record = bc.deviceInfo[name]
	record.Ok = true
	return
}

func unquote(str string) string {
	str = strings.TrimSuffix(str, `"`)
	str = strings.TrimPrefix(str, `"`)
	return str
}
