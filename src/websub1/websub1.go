package websub1


package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"meow.tf/websub"
	"meow.tf/websub/store"
	"meow.tf/websub/worker"
)

// JSONFeed v1.1 标准结构
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
	// 1. 初始化存储
	// 持久化磁盘存储（生产推荐）
	st, err := store.NewBoltStore("subscriptions.db")
	if err != nil {
		log.Fatalf("初始化存储失败: %v", err)
	}
	// 测试临时内存存储（进程重启丢失订阅）
	// st := store.NewMemoryStore()

	// 2. 创建任务worker（异步推送回调）
	w := worker.NewLocalWorker()

	// 3. 构建Hub核心实例
	hub := websub.NewHub(st, w)

	// 4. 启动标准WebSub接口（只用于订阅/取消订阅，没有发布接口）
	http.Handle("/hub", hub)

	// 后台goroutine：模拟业务内部定时产生消息，直接推送
	go func() {
		topicURL := "https://my-feed.example/feed.json" // WebSub规范topic为URL形式
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		seq := 0
		for range ticker.C {
			seq++
			// 构造JSONFeed数据
			feed := JSONFeed{
				Version: "https://jsonfeed.org/version/1.1",
				Title:   "内部消息推送Feed",
				Items: []Item{
					{
						ID:            time.Now().Format("msg-20060102-150405"),
						Title:         "实时消息 #" + string(rune(seq)),
						ContentText:   "由程序内部逻辑直接推送，不经过HTTP publish接口",
						DatePublished: time.Now().Format(time.RFC3339),
					},
				},
			}

			feedBytes, err := json.MarshalIndent(feed, "", "  ")
			if err != nil {
				log.Printf("序列化feed失败: %v", err)
				continue
			}

			// ========== 核心：内部直接推送 ==========
			err = hub.Publish(topicURL, "application/json", feedBytes)
			if err != nil {
				log.Printf("publish调用失败: %v", err)
			} else {
				log.Printf("✅ 成功推送feed，seq=%d", seq)
			}
		}
	}()

	log.Println("WebSub Hub 启动监听 :8080")
	log.Println("订阅地址：POST http://127.0.0.1:8080/hub")
	log.Fatal(http.ListenAndServe(":8080", nil))
}