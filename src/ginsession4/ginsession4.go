package main

import (
	"time"

	"github.com/cpusoft/goutil/belogs"
	_ "github.com/cpusoft/goutil/conf"
	"github.com/cpusoft/goutil/ginserver"
	"github.com/cpusoft/goutil/ginsession"
	"github.com/cpusoft/goutil/jwtutil"
	_ "github.com/cpusoft/goutil/logs"
	"github.com/cpusoft/goutil/zaplogs"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/sync/errgroup"
)

const secrect = "mysecret"

// start server
func main() {
	start := time.Now()
	var g errgroup.Group

	serverHttpPort := "1024"

	gin.SetMode(gin.DebugMode)
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	public := engine.Group("/public")
	{
		public.POST("/hello", Hello)
		public.POST("/login", Login)
	}
	auth := engine.Group("/auth")
	{
		ginsession.RouterGroupRegisterJwt(auth)
		auth.POST("/work", Work)
	}
	g.Go(func() error {
		belogs.Info("startRpServer(): server run http on :", serverHttpPort)
		err := engine.Run(":" + serverHttpPort)
		return err
	})

	if err := g.Wait(); err != nil {
		belogs.Error("startRpServer(): fail, will exit, err:", err)
	}
	belogs.Info("startRpServer(): server end, time(s):", time.Since(start))

}
func Hello(c *gin.Context) {
	ginserver.String(c, "hello")
}

func Login(c *gin.Context) {
	belogs.Info("ChainValidatePartialStart(): use partial chainvalidate")
	m := make(map[string]string)
	m["ownerId"] = "1001"
	m["ownerName"] = "beijing-user1"
	m["opUserId"] = "2002"
	m["opUserName"] = "beijing-user2"
	m["traceId"] = "550e8400-e29b-41d4-a716-446655440000"
	m["opLogId"] = "3003"

	claims := jwtutil.CustomClaims{
		Infos: m,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)), //过期时间
			NotBefore: jwt.NewNumericDate(time.Now()),                    //生效时间（立即生效）
			IssuedAt:  jwt.NewNumericDate(time.Now()),                    //签发时间
		},
	}
	token, _ := jwtutil.GenToken(&claims, secrect)
	ginserver.String(c, token)
}
func Work(c *gin.Context) {
	cxt, err := ginsession.SetToContextWithValue(c)
	if err != nil {
		zaplogs.ErrorJw(nil, "work(): SetToContextWithValue fail", "err", err.Error())
		ginserver.ResponseFail(c, err, nil)
		c.Abort()
		return
	}
	zaplogs.InfoJw(cxt, "Work(): ", "start")
	ginserver.String(c, "work ok")
}
