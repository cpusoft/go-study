package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	websub "github.com/cpusoft/goutil/websubutil"
	"github.com/cpusoft/goutil/websubutil/model"
	"github.com/cpusoft/goutil/websubutil/store/memory"
)

type JSONFeed struct {
	Version     string     `json:"version"`
	Title       string     `json:"title"`
	HomePageURL string     `json:"home_page_url"`
	FeedURL     string     `json:"feed_url"`
	Items       []FeedItem `json:"items"`
}

type FeedItem struct {
	ID            string `json:"id"`
	URL           string `json:"url"`
	Title         string `json:"title"`
	ContentText   string `json:"content_text"`
	DatePublished string `json:"date_published"`
}

func main() {
	subStore := memory.New()
	hub := websub.New(subStore,
		websub.WithURL("http://localhost:8080"),
	)

	// ========== 关键诊断：监听验证事件，打印详细错误 ==========
	hub.AddHandler(func(e *websub.VerificationFailed) {
		log.Printf("❌ [事件] 验证失败: Topic=%s Callback=%s Error=%v",
			e.Subscription.Topic, e.Subscription.Callback, e.Error)
	})

	hub.AddHandler(func(e *websub.Verified) {
		log.Printf("✅ [事件] 验证成功: Topic=%s Callback=%s Expires=%s",
			e.Subscription.Topic, e.Subscription.Callback, e.Subscription.Expires)
	})

	hub.AddHandler(func(e *websub.Publish) {
		log.Printf("📢 [事件] 发布请求: Topic=%s ContentType=%s DataLen=%d",
			e.Topic, e.ContentType, len(e.Data))
	})

	http.Handle("/", hub)
	go func() {
		log.Println("Hub 启动在 http://localhost:8080")
		log.Println(">>> 请在 60 秒内启动 Subscriber 并完成订阅 <<<")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// 异步等待 60 秒后发布
	go func() {
		topic := "https://myblog.com/feed.json"

		log.Println("将在 60 秒后自动发布...")
		time.Sleep(60 * time.Second)

		// 检查 Store
		subs, err := subStore.All(topic)
		if err != nil {
			log.Println("❌ 获取订阅列表失败:", err)
		} else {
			log.Printf("📊 Topic %s 的订阅数量: %d", topic, len(subs))
		}

		if len(subs) == 0 {
			log.Println("⚠️ Store 中没有订阅，尝试手动注入测试订阅...")
			// ========== Workaround：手动添加订阅 ==========
			testSub := model.Subscription{
				Topic:    topic,
				Callback: "http://localhost:9090/callback",
				Secret:   "my-super-secret-123456789",
				Expires:  time.Now().Add(24 * time.Hour),
			}
			if err := subStore.Add(testSub); err != nil {
				log.Println("❌ 手动注入也失败:", err)
				return
			}
			log.Println("✅ 手动注入订阅成功")
		}

		feed := JSONFeed{
			Version:     "https://jsonfeed.org/version/1.1",
			Title:       "我的博客",
			HomePageURL: "https://myblog.com",
			FeedURL:     topic,
			Items: []FeedItem{
				{
					ID:            fmt.Sprintf("post-%d", time.Now().Unix()),
					URL:           "https://myblog.com/posts/new-post",
					Title:         "新文章标题",
					ContentText:   "这是文章内容...",
					DatePublished: time.Now().UTC().Format(time.RFC3339),
				},
			},
		}

		feedJSON, err := json.Marshal(feed)
		if err != nil {
			log.Println("序列化失败:", err)
			return
		}

		err = hub.Publish(topic, "application/feed+json", feedJSON)
		if err != nil {
			log.Println("发布失败:", err)
			return
		}

		log.Println("✅ JSON Feed 已推送到所有订阅者")
	}()

	select {}
}
