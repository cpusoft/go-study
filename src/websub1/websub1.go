package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// JSONFeed v1.1
type JSONFeed struct {
	Version string `json:"version"`
	Title   string `json:"title"`
	Items   []Item `json:"items"`
}

type Item struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	ContentText   string `json:"content_text"`
	DatePublished string `json:"date_published"`
}

func main() {
	// 持久化bolt存储
	st, err := websub.NewBoltStore("subscriptions.db")
	if err != nil {
		log.Fatalf("bolt store init err: %v", err)
	}
	// 测试用内存存储
	// st := websub.NewMemoryStore()

	w := websub.NewLocalWorker()
	hub := websub.NewHub(st, w)

	http.Handle("/hub", hub)

	// 内部业务goroutine，直接推送JSONFeed
	go func() {
		topic := "https://my-feed.example/feed.json"
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()

		i := 0
		for range ticker.C {
			i++
			feed := JSONFeed{
				Version: "https://jsonfeed.org/version/1.1",
				Title:   "内部 JSONFeed 推送测试",
				Items: []Item{
					{
						ID:            time.Now().Format("msg-20060102T150405"),
						Title:         "消息 #" + string(rune(i)),
						ContentText:   "内部直接调用 hub.Publish 发送JSONFeed",
						DatePublished: time.Now().Format(time.RFC3339),
					},
				},
			}

			data, err := json.Marshal(feed)
			if err != nil {
				log.Println("marshal feed err:", err)
				continue
			}

			err = hub.Publish(topic, "application/json", data)
			if err != nil {
				log.Println("publish err:", err)
			} else {
				log.Printf("publish success, seq=%d", i)
			}
		}
	}()

	log.Println("listen :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
