package main

import (
	"context"
	zlog "go-study/log4/zaplog"
	"net/http"

	"go.uber.org/zap"
)

func main() {

	simpleHttpGet("www.sogo.com")
	simpleHttpGet("http://www.sogo.com")
}

func simpleHttpGet(url string) {
	defer zlog.DeferSync()

	zlog.DebugJw(context.TODO(), "Trying to hit GET request for", "url", url)
	resp, err := http.Get(url)
	if err != nil {
		zlog.ErrorJw(context.TODO(), "Error fetching URL:", zap.String("url", url), zap.Errors("err", []error{err}))
	} else {
		zlog.InfoJw(context.TODO(), "Success! statusCode", "status", resp.Status, "url", url)
		resp.Body.Close()
	}
}
