package main

import (
	"fmt"
	"net/http"

	"github.com/cpusoft/gin-contrib-sessions"
	"github.com/cpusoft/gin-contrib-sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	igorePaths := []string{"/login"}

	r := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store), AuthMiddleWare(igorePaths))

	r.GET("/login", login)

	r.GET("/hello", hello)

	r.Run(":8000")
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

		session := sessions.Default(c)
		if session.Get("test@mail.com") != "test.access.token" {
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

func login(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("test@mail.com", "test.access.token")
	session.Save()
	fmt.Println("login and save session")
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
