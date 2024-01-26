package main

import (
	"fmt"

	"github.com/cpusoft/goutil/conf"
	"github.com/cpusoft/goutil/jsonutil"
)

type LocationPerformanceModel struct {
	IsInMainland bool   `json:"isInMainland"`
	City         string `json:"city"`
	Region       string `json:"region"`
	Rir          string `json:"rir"`
}

func main() {
	locationPerformanceModel := &LocationPerformanceModel{}
	location := conf.String("performance::location")
	fmt.Println(location)
	err := jsonutil.UnmarshalJson(location, locationPerformanceModel)
	fmt.Println(locationPerformanceModel, err)
}
