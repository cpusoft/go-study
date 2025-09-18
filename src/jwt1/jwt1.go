package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	mapClaims := jwt.MapClaims{
		"iss": "程序员陈明勇",
		"sub": "chenmingyong.cn",
		"aud": "Programmer",
		"exp": time.Now().Add(time.Second * 10).UnixMilli(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	fmt.Println(token) // true
	fmt.Println(token.Claims.GetIssuer())

	jwtKey := "test"
	signedToken, err := token.SignedString([]byte(jwtKey))
	fmt.Println(signedToken, err) // eyJ0eXAiOiJKV1QiLCJhbGciOiJITU9fMjU2In0.eyJpc3MiOiL

	// 解析 jwt
	token1, err1 := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	}, jwt.WithExpirationRequired())
	if err1 != nil {
		fmt.Println("Error parsing token:", err1)
		return
	}
	if !token1.Valid {
		fmt.Println("token1 is not valid:", err1)
		return
	}
	fmt.Println(token1.Claims)

}
