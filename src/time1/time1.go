package main

import (
	"fmt"
	"time"

	"github.com/cpusoft/goutil/jsonutil"
)

type Test struct {
	Key time.Time `json:"key"`
}

func main() {
	/*
		date := "2021-09-24 18:57:22"
		//	t, err := time.Parse.EST(time.RFC3339, date)
		t, err := time.ParseInLocation("2006-01-02 15:04:05", date, time.Local)

		fmt.Println(t, err)

		curNow := time.Now().Local()
		fmt.Println(curNow)

		now := time.Now()
		s := now.Local().Format("2006-01-02T15:04:05-0700")
		fmt.Println(s)
	*/
	expireTimeStr := `{"key":"2023-06-30T19:36:56.571798033+08:00"}`
	nowStr := `{"key":"2023-06-30T18:47:08.9244639+08:00"}`

	var expireT Test
	var nowT Test
	err := jsonutil.UnmarshalJson(expireTimeStr, &expireT)
	fmt.Println("expireT", expireT, err)

	err = jsonutil.UnmarshalJson(nowStr, &nowT)
	fmt.Println("nowT", nowT, err)
	b := nowT.Key.After(expireT.Key)
	fmt.Println(b)
}
