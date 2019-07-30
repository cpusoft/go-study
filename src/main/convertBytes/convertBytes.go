package main

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

func main() {
	str := `00:9c:71:13:39:cc:2a:32:0c:5f:2d:1e:5d:b8:df:fc:b3:6e:27:39:68:c0:42:a1:b9:13:3f:e9:73:86`
	split := strings.Split(str, ":")
	b := make([]int, 0)
	for _, one := range split {
		bb, _ := strconv.ParseInt(one, 16, 0)
		b = append(b, int(bb))
	}
	fmt.Println(b)

	str1 := strings.Replace(str, ":", "", -1)
	n := new(big.Int)
	n, ok := n.SetString(str1, 16)
	if !ok {
		fmt.Println("SetString: error")
		return
	}
	fmt.Println(n)
}
