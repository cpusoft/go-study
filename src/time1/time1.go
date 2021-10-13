package main

import (
	"fmt"
	"time"
)

func main() {
	date := "2021-09-24 18:57:22"
	//	t, err := time.Parse.EST(time.RFC3339, date)
	t, err := time.ParseInLocation("2006-01-02 15:04:05", date, time.Local)

	fmt.Println(t, err)

	curNow := time.Now().Local()
	fmt.Println(curNow)
}
