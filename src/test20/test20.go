package main

import (
	"bytes"
	"fmt"

	"github.com/cpusoft/goutil/convert"
)

func Int64sToInString(s []int64) string {
	if len(s) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	buffer.WriteString("(")
	for i := 0; i < len(s); i++ {
		buffer.WriteString(convert.ToString(s[i]))
		if i < len(s)-1 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString(")")
	return buffer.String()

}

func main() {
	s := make([]int64, 0)
	s = append(s, 1)
	s = append(s, 2)
	str := Int64sToInString(s)
	fmt.Println(str)
}
