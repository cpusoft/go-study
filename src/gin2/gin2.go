package main

import (
	"fmt"
	"net/http"

	"github.com/cpusoft/goutil/osutil"
	"github.com/gin-gonic/gin"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)
	//engine := gin.New()
	//engine.Use(gin.Recovery())

	engine := gin.Default()

	router := engine
	router.GET("/someGet", getting)
	router.GET("/user/:name", getname)
	router.GET("/lastname", getlastname)
	router.POST("/upload", upload)

	authorized := router.Group("/", AuthRequired())
	authorized.POST("/login", login)

	router.Any("/startpage", startPage)

	router.GET("/someXML", someXML)
	router.GET("/someJson", someJson)

	//router.Static("/dnsviz_files", "../../templates/dnsviz_files")
	//router.Static("./dnsviz_files", "templates/dnsviz_files")
	router.LoadHTMLGlob("template/*.html")
	//ok
	//router.Static("/dnsviz_files", "E:/Go/go-study/template/static/")
	//router.Static("/static", "E:/Go/go-study/template/static/")
	router.Static("/static", osutil.GetParentPath()+"/template/static/")
	//router.Static("/dnsviz_files", osutil.GetParentPath()+"/template/static/")

	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/login2", login2)
	router.GET("/dnsviz", dnsviz)
	router.GET("/vlabs", vlabs)

	router.Run()
}

func getting(c *gin.Context) {
	fmt.Println("ok")
	c.JSON(http.StatusOK, gin.H{
		"status":  "posted",
		"message": "hello world",
		"nick":    "yes",
	})
}
func getname(c *gin.Context) {
	name := c.Param("name")
	c.String(http.StatusOK, "Hello %s", name)
}
func getlastname(c *gin.Context) {
	name := c.Query("name")
	c.String(http.StatusOK, "lastname %s", name)
}
func upload(c *gin.Context) {
	file, _ := c.FormFile("file")
	fmt.Println(file.Filename)
	c.SaveUploadedFile(file, `G:\Download\z`)
	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}
func login(c *gin.Context) {
	c.String(http.StatusOK, "login")
}
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("before middleware")
		c.Set("request", "clinet_request")
		c.Next()
		fmt.Println("after middleware")
	}
}

type Person struct {
	Name    string `form:"name" json:"name"`
	Address string `form:"address" json:"address"`
}

func startPage(c *gin.Context) {
	var person Person
	// 如果是 `GET`, 只使用 `Form` 绑定引擎 (`query`) 。
	// 如果 `POST`, 首先检查 `content-type` 为 `JSON` 或 `XML`, 然后使用 `Form` (`form-data`) 。
	// 在这里查看更多信息 https://github.com/gin-gonic/gin/blob/master/binding/binding.go#L48
	if c.ShouldBind(&person) == nil {
		fmt.Println(person.Name)
		fmt.Println(person.Address)
	}

	c.String(200, "Success")
}

func someXML(c *gin.Context) {
	c.XML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
}
func someJson(c *gin.Context) {
	p := Person{Name: "name111", Address: "addressss"}
	c.JSON(http.StatusOK, p)
}
func login2(c *gin.Context) {
	c.HTML(http.StatusOK, "login2.html", gin.H{
		"title": "Main website",
	})
}
func dnsviz(c *gin.Context) {
	c.HTML(http.StatusOK, "dnsviz.html", gin.H{
		"title": "Main website",
	})
}
func vlabs(c *gin.Context) {
	c.HTML(http.StatusOK, "vlabs.html", gin.H{})
}
