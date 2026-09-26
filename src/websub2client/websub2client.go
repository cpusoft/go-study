package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	HubURL       = "http://localhost:8080"
	TopicURL     = "https://myblog.com/feed.json"
	CallbackHost = "http://localhost:9090"
	CallbackPath = "/callback"
	Secret       = "my-super-secret-123456789"
	LeaseSeconds = 86400
)

func sendSubscriptionRequest(mode string) error {
	// 手动构造 form data
	var sb strings.Builder
	first := true
	data := map[string]string{
		"hub.mode":          mode,
		"hub.topic":         TopicURL,
		"hub.callback":      CallbackHost + CallbackPath,
		"hub.lease_seconds": fmt.Sprintf("%d", LeaseSeconds),
	}
	if Secret != "" {
		data["hub.secret"] = Secret
	}

	for k, v := range data {
		if !first {
			sb.WriteString("&")
		}
		sb.WriteString(k + "=" + v)
		first = false
	}

	// 用构造好的 body 直接发请求
	resp, err := http.Post(HubURL, "application/x-www-form-urlencoded", strings.NewReader(sb.String()))
	if err != nil {
		return fmt.Errorf("请求 Hub 失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Hub 拒绝: %s, 响应: %s", resp.Status, string(body))
	}
	log.Printf("已向 Hub 发送 %s 请求，等待验证...", mode)
	return nil
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	// 在最开头记录所有到达的请求，不管方法对不对
	log.Printf("[🌐 请求到达] %s %s from %s Headers=%v",
		r.Method, r.URL.String(), r.RemoteAddr, r.Header)

	switch r.Method {
	case http.MethodGet:
		handleVerification(w, r)
	case http.MethodPost:
		handleContentDistribution(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func handleVerification(w http.ResponseWriter, r *http.Request) {
	hubMode := r.URL.Query().Get("hub.mode")
	hubTopic := r.URL.Query().Get("hub.topic")
	challenge := r.URL.Query().Get("hub.challenge")
	leaseSeconds := r.URL.Query().Get("hub.lease_seconds")

	log.Printf("[验证] mode=%s topic=%s challenge=%s lease=%s",
		hubMode, hubTopic, challenge, leaseSeconds)

	if hubTopic != TopicURL {
		log.Printf("[验证] 拒绝未知 topic: %s", hubTopic)
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(challenge))
	log.Printf("[验证] 已确认 %s", hubMode)
}

func handleContentDistribution(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("[推送] 读取 body 失败:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	contentType := r.Header.Get("Content-Type")
	signature := r.Header.Get("X-Hub-Signature")
	linkHeader := r.Header.Get("Link")

	log.Printf("[推送] Content-Type=%s Signature=%s Link=%s BodyLen=%d",
		contentType, signature, linkHeader, len(body))

	if len(body) > 0 {
		log.Printf("[推送] Body %s", string(body))
	}

	if Secret != "" {
		if signature == "" {
			log.Println("[推送] 错误: 期望签名但未收到")
			http.Error(w, "Signature Required", http.StatusForbidden)
			return
		}
		if !verifySignature(body, signature, Secret) {
			log.Println("[推送] 错误: 签名验证失败")
			http.Error(w, "Invalid Signature", http.StatusForbidden)
			return
		}
		log.Println("[推送] 签名验证通过 ✅")
	}

	log.Printf("[推送] 收到完整内容，处理成功")
	w.WriteHeader(http.StatusNoContent)
}

func verifySignature(body []byte, signatureHeader, secret string) bool {
	parts := strings.SplitN(signatureHeader, "=", 2)
	if len(parts) != 2 {
		return false
	}
	algoName := parts[0]
	expectedHex := parts[1]

	var hasher func() hash.Hash
	switch strings.ToLower(algoName) {
	case "sha1":
		hasher = sha1.New
	case "sha256":
		hasher = sha256.New
	case "sha384":
		hasher = sha512.New384
	case "sha512":
		hasher = sha512.New
	default:
		log.Printf("不支持的签名算法: %s", algoName)
		return false
	}

	mac := hmac.New(hasher, []byte(secret))
	mac.Write(body)
	computed := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(computed), []byte(expectedHex))
}

func main() {
	// 先注册一个全局 catch-all，记录所有到达的请求
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[🌐 全局未匹配] %s %s from %s", r.Method, r.URL.String(), r.RemoteAddr)
		http.NotFound(w, r)
	})
	http.HandleFunc(CallbackPath, callbackHandler)

	go func() {
		log.Println("Subscriber callback 启动在 " + CallbackHost + CallbackPath)
		log.Fatal(http.ListenAndServe(":9090", nil))
	}()

	time.Sleep(200 * time.Millisecond)

	if err := sendSubscriptionRequest("subscribe"); err != nil {
		log.Fatal(err)
	}

	log.Println("Subscriber 运行中...")
	select {}
}
