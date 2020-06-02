package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {

	router := gin.Default()

	router.GET("/login", func(c *gin.Context) {

		//c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
		ResponseOk(c, "aaa")
	})
	// Listen and serve on 0.0.0.0:8080
	router.Run(":8080")
}

type ResponseModel struct {
	Result string      `json:"result"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data,omitempty"`
}

func ResponseOk(c *gin.Context, v interface{}) {
	ret := ResponseModel{Result: "ok", Msg: "", Data: v}
	responseJSON(c, http.StatusOK, &ret)
}

func ResponseFail(c *gin.Context, err error, v interface{}) {
	ret := ResponseModel{Result: "fail", Msg: err.Error(), Data: v}
	responseJSON(c, http.StatusOK, &ret)
}

func responseJSON(c *gin.Context, status int, v interface{}) {
	c.JSON(status, v)
}

func loginEndpoint(c *gin.Context) {

}

func submitEndpoint(c *gin.Context) {

}
func readEndpoint(c *gin.Context) {

}
func analyticsEndpoint(c *gin.Context) {

}
