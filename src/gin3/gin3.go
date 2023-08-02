package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)
	//engine := gin.New()
	//engine.Use(gin.Recovery())

	engine := gin.Default()

	router := engine
	router.POST("/someGet", getting)

	router.Run()
}

type DistributedDeltaModel struct {
	Id        int    `json:"id"`
	NotifyUrl string `json:"notifyUrl"`
	DeltaUrl  string `json:"deltaUrl"`
	Serial    uint64 `json:"serial"`

	CenterNodeUrl string `json:"centerNodeUrl"`
	Index         uint64 `json:"index"`
	SyncLogId     uint64 `json:"syncLogId"`
}

func getting(c *gin.Context) {
	fmt.Println("getting")
	data, _ := c.GetRawData()
	fmt.Println(data, "\r\n", string(data))
	ds := make([]DistributedDeltaModel, 0)

	_ = json.Unmarshal(data, &ds)
	fmt.Println("body:", ds)

	c.JSON(http.StatusOK, gin.H{
		"status":  "posted",
		"message": "hello world",
		"nick":    "yes",
	})
}
