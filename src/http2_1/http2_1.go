package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// https://blog.csdn.net/testapl/article/details/131089133
func main() {
	//r := gin.Default()
	r := gin.New()
	r.UseH2C = true
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	log.Fatal(r.RunTLS(":443", "cert.pem", "key.pem"))
}
