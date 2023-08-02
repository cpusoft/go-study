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
			"deltaUrl": "https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/c3cd7c24-12cb-4abc-8fd2-5e2bcbb85ae6/7d81c8e0-e560-4b46-8ebd-89dca0f8ce06/390/delta.xml",
			"errMsg": "",
			"notifyUrl": "https://rpki-rrdp.us-east-2.amazonaws.com/rrdp/c3cd7c24-12cb-4abc-8fd2-5e2bcbb85ae6/notification.xml",
			"publishCount": 2,
			"serial": 390,
			"sessionId": "7d81c8e0-e560-4b46-8ebd-89dca0f8ce06",
			"syncLogId": 6,
			"syncTime": "2023-07-13T16:33:20.497430115+08:00",
			"withdrawCount": 0
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
