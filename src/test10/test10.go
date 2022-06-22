package main

import (
	"fmt"
	"time"
)

func main() {

	ss := make([]string, 0)
	ss = append(ss, "")
	for _, k := range ss {
		fmt.Println("k", k)
	}

	fmt.Printf("时间戳（秒）：%v;\n", time.Now().Unix())
	fmt.Printf("时间戳（纳秒）：%v;\n", time.Now().UnixNano())
	fmt.Printf("时间戳（微妙）：%v;\n", time.Now().UnixNano()/1e3)
	fmt.Printf("时间戳（毫秒）：%v;\n", time.Now().UnixNano()/1e6)
	fmt.Printf("时间戳（纳秒转换为秒）：%v;\n", time.Now().UnixNano()/1e9)

}
