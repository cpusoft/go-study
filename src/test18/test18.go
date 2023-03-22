package main

import (
	"fmt"

	"github.com/cpusoft/goutil/jsonutil"
)

type Param struct {
	Rp []string `json:"rpstir2-rp"`
	Vc []string `json:"vc"`
}

func main() {
	//"rpstir2-rp:tal&entiresync&directsync&parsevalidate&chainvalidate&clear&sys&statistic&roacompete&slurm&rss&synclocal"
	rp := make([]string, 0)
	rp = append(rp, "tal")
	rp = append(rp, "entiresync")
	rp = append(rp, "directsync")
	rp = append(rp, "parsevalidate")
	rp = append(rp, "chainvalidate")
	rp = append(rp, "clear")
	rp = append(rp, "sys")
	rp = append(rp, "statistic")
	rp = append(rp, "roacompete")
	rp = append(rp, "slurm")
	rp = append(rp, "rss")
	rp = append(rp, "synclocal")

	vc := make([]string, 0)
	vc = append(vc, "rtrp")
	vc = append(vc, "rtrtcp")

	param := Param{}
	param.Rp = rp
	param.Vc = vc
	fmt.Println(jsonutil.MarshalJson(param))
}
