package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/cpusoft/goutil/jsonutil"
	"github.com/gin-gonic/gin"
)

type DistributedDeltaModel struct {
	Id        int    `json:"id"`
	NotifyUrl string `json:"notifyUrl"`
	DeltaUrl  string `json:"deltaUrl"`
	Serial    uint64 `json:"serial"`

	CenterNodeUrl string `json:"centerNodeUrl"`
	Index         uint64 `json:"index"`
	SyncLogId     uint64 `json:"syncLogId"`
}

func main() {
	s := `[
		{
			"id": 17385,
			"notifyUrl": "https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/708aafaf-00b4-485b-854c-0b32ca30f57b/notification.xml",
			"deltaUrl": "https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/708aafaf-00b4-485b-854c-0b32ca30f57b/7c2d14ed-ded9-4777-9885-dc9d35403305/391/delta.xml",
			"serial": 391,
			"centerNodeUrl": "https://10.1.135.104:8071",
			"index": 25,
			"syncLogId": 3
		}
	]
	`
	ds := make([]DistributedDeltaModel, 0)
	err := jsonutil.UnmarshalJson(s, &ds)
	fmt.Println(ds, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(s))
	c.Request.Header.Add("Content-Type", "application/json") // set fake content-type

	ds2 := make([]DistributedDeltaModel, 0)
	err = c.ShouldBindJSON(&ds2)
	fmt.Println(ds2, err)
}
