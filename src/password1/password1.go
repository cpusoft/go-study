package main

import (
	"fmt"

	"github.com/cpusoft/goutil/hashutil"
)

func main() {
	password := `admin123`
	salt := `d16c1ee7-6811-4e29-af53-29d9ad7bb9ac`
	haspassword := hashutil.Sha256([]byte(password + salt))
	fmt.Println(haspassword)
}
