package main

import (
	"go-study/log4/zaplog"
	"net/http"
	"time"

	"github.com/cpusoft/goutil/belogs"
	"github.com/cpusoft/goutil/conf"
	"github.com/cpusoft/goutil/ginsession"
	"github.com/cpusoft/goutil/jwtutil"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func main() {
	serverHttpPort := conf.String("rpkix-admin::serverHttpPort")
	serverHttpsPort := conf.String("rpkix-admin::serverHttpsPort")
	belogs.Info("startServer(): will start server on ", serverHttpPort, serverHttpsPort)

	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)
	engine := gin.New()
	engine.Use(gin.Recovery())

	ginsession.RegisterJwt(engine)

}

// loginHandler 处理用户登录并生成JWT令牌
func loginHandler(c *gin.Context) {

	// 设置令牌过期时间为2小时
	expirationTime := time.Now().Add(2 * time.Hour)

	// 创建JWT声明
	claims := &jwtutil.CustomClaims{
		UserName: "zhangsan"
		UserId: 1,
		TraceId: "abcdefg",
		LogId:2,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "gin-jwt-example",
		},
	}
	tokenString, err := jwtutil.GenToken(claims, "my_secret_key")
	if err != nil {
		zaplog.ErrorJw("loginHandler(): GenToken fail:", "claims:",cliaims,"err:",err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}
	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	// 返回令牌给客户端
	c.JSON(http.StatusOK, gin.H{
		"token":     tokenString,
		"expiresAt": expirationTime.Unix(),
	})
}
