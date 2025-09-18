package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func main() {
	r := gin.Default()

	// 顺序关键：Tracing中间件需在最前
	r.Use(TracingMiddleware())
	r.Use(AuthMiddleware())   // JWT鉴权
	r.Use(LoggerMiddleware()) // 访问日志

	// 受保护的路由
	r.GET("/protected", func(c *gin.Context) {
		userID := c.GetString("userID")
		c.JSON(200, gin.H{"user": userID})
	})

	r.Run(":8080")
}
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Missing token"})
			return
		}

		// 解析JWT并验证
		token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte("your-secret-key"), nil
			})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}

		// 将解析后的用户信息存入上下文
		if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
			c.Set("userID", claims.Subject)
		}

		c.Next() // 继续执行后续中间件
	}
}
func TracingMiddleware() gin.HandlerFunc {
	start := time.Now()
	return func(c *gin.Context) {
		// 生成唯一traceID（优先从请求头获取，适用于跨服务调用）
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String() // 生成新traceID
		}

		// 存入Gin上下文和标准Context
		ctx := context.WithValue(c.Request.Context(), "traceID", traceID)
		c.Set("traceID", traceID)
		c.Request = c.Request.WithContext(ctx)

		// 继续处理请求
		c.Next()

		go test1(ctx)
		// 请求结束后记录日志（含状态码、延迟等）
		log.Printf("[%s] %s %d - %d | TraceID=%s",
			c.Request.Method, c.Request.URL, c.Writer.Status(), time.Since(start), traceID)
	}
}
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		path := c.Request.URL.Path
		method := c.Request.Method

		logMessage := fmt.Sprintf("| %3d | %13v | %15s | %s  %s\n",
			status,
			latency,
			c.ClientIP(),
			method,
			path,
		)

		// 将日志信息输出到控制台
		fmt.Print(logMessage)
	}
}

func test1(ctx context.Context) {
	start := time.Now()
	traceID, _ := ctx.Value("traceID").(string)
	log.Printf(" %d | TraceID=%s",
		time.Since(start), traceID)

}
