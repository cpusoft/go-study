package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type Field struct {
	Field interface{} `json:"field"`
	Value interface{} `json:"value"`
}

func main() {
	fields := make([]Field, 0)
	f1 := Field{Field: "pageIndex", Value: 1}
	f2 := Field{Field: "group", Value: 0}
	f3 := Field{Field: "Searchtype", Value: 1}
	f4 := Field{Field: "keyword", Value: "孙大文"}
	f5 := Field{Field: "recommend", Value: 1}
	f6 := Field{Field: 4, Value: ""}
	f7 := Field{Field: 5, Value: ""}
	f8 := Field{Field: 6, Value: ""}
	f9 := Field{Field: 7, Value: ""}
	f10 := Field{Field: 8, Value: ""}
	f11 := Field{Field: 9, Value: ""}
	f12 := Field{Field: 10, Value: ""}

	fields = append(fields, f1)
	fields = append(fields, f2)
	fields = append(fields, f3)
	fields = append(fields, f4)
	fields = append(fields, f5)
	fields = append(fields, f6)
	fields = append(fields, f7)
	fields = append(fields, f8)
	fields = append(fields, f9)
	fields = append(fields, f10)
	fields = append(fields, f11)
	fields = append(fields, f12)
	fmt.Println(fields)
	b, _ := json.Marshal(fields)
	searchInfos := base64.URLEncoding.EncodeToString(b)
	fmt.Println(searchInfos)
}
