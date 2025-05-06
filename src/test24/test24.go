package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	syncTransferTime, err := time.ParseInLocation("2006-01-02 15:04:05", "2025-02-27 10:53:01", time.Local)
	if err != nil {
		fmt.Println(err)
		return

	}
	notBeforTime := now.Local().Add(-8 * time.Hour)
	notAfterTime := now.Local().Add(8 * time.Hour)
	if syncTransferTime.Before(notBeforTime) || syncTransferTime.After(notAfterTime) {
		fmt.Println(err)
		return
	}
	fmt.Println(notBeforTime, syncTransferTime, notAfterTime)
}
