package utils

import (
	"fmt"
	"sync/atomic"
	"time"
)

var (
	cacheWeek atomic.Value
)

func init() {
	cacheWeek.Store(makeLastWeek())
	go midnight()
}

func leftBeforeMidnight() time.Duration {
	tm := time.Now()
	return time.Date(tm.Year(), tm.Month(), tm.Day()+1, 0, 0, 0, 0, tm.Location()).Sub(tm) + 5*time.Second
}

func midnight() {
	for {
		time.Sleep(leftBeforeMidnight())
		cacheWeek.Store(makeLastWeek())
	}
}

func makeLastWeek() []string {

	ctm := time.Now()
	weekDates := make([]string, 7)

	for i := range weekDates {
		tm := time.Date(ctm.Year(), ctm.Month(), ctm.Day()-i, 0, 0, 0, 0, ctm.Location())
		weekDates[i] = fmt.Sprintf("%04d-%02d-%02d", tm.Year(), tm.Month(), tm.Day())
	}

	return weekDates
}

func GetLastWeek() []string {
	return cacheWeek.Load().([]string)
}
