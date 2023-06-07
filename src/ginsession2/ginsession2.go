package main

import (
	"fmt"
	"net/http"

	"github.com/cpusoft/gin-contrib-sessions"
	"github.com/cpusoft/gin-contrib-sessions/cookie"
	"github.com/cpusoft/gin-contrib-sessions/memstore"
	"github.com/gin-gonic/gin"
)

func main() {
	igorePaths := []string{"/login"}

	r := gin.Default()

	//store := cookie.NewStore([]byte("secret"))
	//r.Use(sessions.Sessions("mysession", store), AuthMiddleWare(igorePaths))
	UseCookie("cookieKey", "memKey", 36*60, false, false, r)
	r.Use(AuthMiddleWare(igorePaths))

	r.GET("/login", login)

	r.GET("/hello", hello)

	r.Run(":8000")
}

// maxAge(seconds): 30 * 60 ;
// secure: only use cookie on https ;
// httpOnly: only in http, cannot in js ;
func UseCookie(cookieKey, memKey string, maxAge int, secure, httpOnly bool, engine *gin.Engine) {

	cookieStore := cookie.NewStore([]byte("ginserver-cookie-secret"))
	fmt.Println("cookieStore:", cookieStore)
	cookieStore.Options(sessions.Options{
		Path:     "/",
		Domain:   "/",
		MaxAge:   maxAge,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
	engine.Use(sessions.SessionsContextKey(cookieKey, cookieKey, cookieStore))

	memStore := memstore.NewStore([]byte("ginserver-mem-secret"))
	fmt.Println("memStore:", memStore)
	engine.Use(sessions.SessionsContextKey(memKey, memKey, memStore))

}

func AuthMiddleWare(ignorePaths []string) gin.HandlerFunc {
	return func(c *gin.Context) {

		path := c.Request.URL.Path
		fmt.Println(ignorePaths)
		fmt.Println(path)
		for _, s := range ignorePaths {
			if s == path {
				fmt.Println("in ignorePaths:", s, path)
				c.Next()
				return
			}
		}

		//cookieSession := sessions.DefaultContextKey("cookieKey", c)
		//cookieGet := cookieSession.Get("test@cookie.com")
		//fmt.Println(cookieGet)

		memSession := sessions.DefaultContextKey("memKey", c)
		memGet := memSession.Get("test@mem.com")
		fmt.Println(memGet)

		if memGet != "test.access.mem" {
			fmt.Println("fail, have not logined")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		} else {
			fmt.Println("success")
			//c.JSON(http.StatusOK, gin.H{"hello": "success"})
			c.Next()
		}

	}
}

/*
func login(c *gin.Context) {
	//cookieSession := sessions.DefaultContextKey("cookieKey", c)
	//cookieSession.Set("test@cookie.com", "test.access.cookie")
	//cookieSession.Save()
	//fmt.Println("login and save cookie")

	memSession := sessions.DefaultContextKey("memKey", c)
	memSession.Set("test@mem.com", "test.access.mem")
	memSession.Save()
	fmt.Println("login and save mem")

	c.JSON(http.StatusOK, nil)
}
*/
func login(c *gin.Context) {

	test2 := `
互联网、电子商务以及新的商业模式出现，让传统的商业模式受到巨大的冲击。商业世界的后浪们也逐步把前浪都拍在了沙滩上。

作为全球知名的消费品（服装、玩具等）设计、开发、采购及物流跨国企业利丰（00494，HK）退市的时间也被确定，将定在5月27日收盘后。而让人唏嘘的是，就在几年前利丰还是恒生指数成分股，港股中响当当的蓝筹股。



图片来源：摄图网（图文无关）

利丰的管理层表示，疫情的影响让公司业绩受阻。但是每日经济新闻（微信号：nbdnews）记者还是注意到，利丰近年来的业绩一直在萎缩，在最近五年时间里其营业额都在下降。而让投资者尴尬的是，利丰的股价在过去十年来更是经历了雪崩式的下滑。

私有化引爆短期股价

从2011年的22.779港元，一路下滑至今年最低的0.435港元（股价已前复权处理），利丰最近十年的股价一路向下，如果你持有了它，恐怕不是一件愉快的事。

而就在今年3月23日，利丰的股价却顿时大涨88%。为何出现这样的情况？原来是在3月20日（周五），利丰发布了公告称，要约人拟将公司私有化并退市。此次私有化的回收价为每股1.25港元，较此前0.5港元的股价溢价了150%。利丰指出，对股东而言，高溢价的回收价是具有吸引力的。

而随着私有化的消息披露后，利丰的股价也在3月23日爆升之后稳步上涨，逼近了1.25港元的回收价，截至5月14日收盘，最新股价为1.24港元`
	cookieSession := sessions.DefaultContextKey("cookieKey", c)
	fmt.Println("cookieSession:", cookieSession)
	cookieSession.Set("test@cookie.com", test2)
	err := cookieSession.Save()
	fmt.Println("login and save cookie,", err)

	memSession := sessions.DefaultContextKey("memKey", c)
	fmt.Println("memSession:", memSession)
	memSession.Set("test@mem.com", test2)
	memSession.Save()
	fmt.Println("login and save mem2")

	c.JSON(http.StatusOK, nil)
}

/* check session in func */
/*
func hello(c *gin.Context) {
	session := sessions.Default(c)

	if session.Get("test@mail.com") != "test.access.token" {
		fmt.Println("hello fail, have not logined")
		c.JSON(http.StatusOK, gin.H{"hello": "fail, have not logined"})
	} else {
		c.JSON(http.StatusOK, gin.H{"hello": "success"})
	}

}
*/

func hello(c *gin.Context) {
	fmt.Println("found session")
	c.JSON(http.StatusOK, gin.H{"hello": "found session, success"})

}
