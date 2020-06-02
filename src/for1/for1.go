package main

import (
	"fmt"
)

func main() {
	ss := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	var text string
	for i := 0; i < len(ss); i = i + 3 {
		text = ss[i]
		if i+1 < len(ss) {
			text = text + "  |   " + ss[i+1]
		}
		if i+2 < len(ss) {
			text = text + "  |   " + ss[i+2]
		}
		fmt.Println(text)
	}
}
