package main

import (
	"fmt"
	"strings"
)

func main() {
	var builder strings.Builder
	// 写入换行符
	builder.WriteString("第一行")
	builder.WriteString("\n") // 实际写入换行符
	builder.WriteString("第二行")
	builder.WriteString("\n") // 实际写入换行符
	builder.WriteString("第二行")
	result := builder.String()

	fmt.Println("直接打印：")
	fmt.Print(result)

}
