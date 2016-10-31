package utils

import (
	"fmt"
	"time"
)

type WTime struct {
	Start time.Time
}

func WTimeStart() WTime {
	return WTime{Start: time.Now()}
}

func (wtm WTime) Stop() {
	fmt.Println(time.Now().Sub(wtm.Start))
}
