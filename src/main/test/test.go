package main

import (
	"bytes"
	"errors"
	"fmt"
)

func GetIndexLast0000(oldb []byte) int {
	endbytes := []byte{0x00, 0x00}
	pos := bytes.Index(oldb, endbytes)
	fmt.Println("pos:", pos)
	for len(oldb) > 0 &&
		pos > 0 &&
		len(oldb) > pos+2*len(endbytes) &&
		len(oldb) > pos+4 &&
		bytes.Equal(oldb[pos+2:pos+4], endbytes) {
		pos += 2
		fmt.Println("GetTopHierarchyFor00(): pos:", pos)
	}
	return pos

}

func GetTopHierarchyFor00(oldb []byte) int {
	top := 0
	endbytes := []byte{0x00, 0x00}
	pos := bytes.LastIndex(oldb, endbytes)
	fmt.Println("pos:", pos, "   len(oldb):", len(oldb), " pos+len(endbytes) ", pos+len(endbytes))
	for pos > 0 && len(oldb) == pos+len(endbytes) {
		oldb = oldb[:pos]
		top += 1
		pos = bytes.LastIndex(oldb, endbytes)
	}
	fmt.Println("GetTopHierarchyFor00(): top:", top)
	return top
}

// UTC 类型是 短的年
func DecodeUTCTime(data []byte) (ret string, err error) {
	if len(data) < 13 {
		return "", errors.New("DecodeUTCTime fail")
	}
	year := "20" + string(data[0:2])
	month := string(data[2:4])
	day := string(data[4:6])
	hour := string(data[6:8])
	minute := string(data[8:10])
	second := string(data[10:12])
	z := string(data[12])
	return year + "-" + month + "-" + day + " " + hour + ":" + minute + ":" + second + z, nil
}

// Generalized 是长的年
func DecodeGeneralizedTime(data []byte) (ret string, err error) {
	if len(data) < 15 {
		return "", errors.New("DecodeGeneralizedTime fail")
	}
	year := string(data[0:4])
	month := string(data[4:6])
	day := string(data[6:8])
	hour := string(data[8:10])
	minute := string(data[10:12])
	second := string(data[12:14])
	z := string(data[14])
	return year + "-" + month + "-" + day + " " + hour + ":" + minute + ":" + second + z, nil
}

func PrintBytes(name string, buf *[]byte, indent string) {
	if len(name) > 0 {
		fmt.Print(name)
	}
	data_lines := make([]string, (len(*buf)/30)+1)

	for i, b := range *buf {
		data_lines[i/30] += fmt.Sprintf("%02x ", b)
	}

	for i := 0; i < len(data_lines); i++ {
		//belogs.Debug(indent + data_lines[i])
		fmt.Print(indent + data_lines[i] + "\n\n")
	}
}
func InversionCount(a *[]int) int {
	if len(*a) <= 1 {
		return 0
	}
	mid := len(*a) / 2
	left := (*a)[:mid]
	right := (*a)[mid:]
	leftCount := InversionCount(&left)   //not being sorted properly due to reference issues
	rightCount := InversionCount(&right) //not being sorted properly due to reference issues

	res := make([]int, 0, len(right)+len(left)) //temp slice to hold the sorted left side and right side

	iCount := 1 //mergeCount(left, right, &res)

	*a = res
	fmt.Println(a) //a in the end is not sorted properly for most cases
	return iCount + leftCount + rightCount
}
func FmtBytes(buf *[]byte) {
	fmt.Printf("%p\r\n", buf)
}
func main() {
	data1 := []byte("sss")
	data := []byte{0x31, 0x38, 0x00, 0x31, 0x38, 0x00}
	fmt.Printf("%p\r\n", &data1)
	FmtBytes(&data1)
	left := data1[:]
	FmtBytes(&left)

	fmt.Printf("%p\r\n", &data)
	FmtBytes(&data)
	left = data[:]
	FmtBytes(&left)
	/*
		pos := GetIndexLast0000(data)
		fmt.Println(pos)

			data := []byte{0x31, 0x38, 0x31, 0x30, 0x30, 0x33, 0x30, 0x37, 0x34, 0x34, 0x35, 0x36, 0x5A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
			pos := GetTopHierarchyFor00(data)
			fmt.Println(pos)
			/*
				data = []byte{0x31, 0x38, 0x31, 0x30, 0x30, 0x33, 0x30, 0x37, 0x34, 0x34, 0x35, 0x36, 0x5A}
				ret, _ := DecodeUTCTime(data)
				fmt.Println(ret)

				data = []byte{0x32, 0x30, 0x31, 0x38, 0x30, 0x36, 0x32, 0x38, 0x32, 0x30, 0x33, 0x31, 0x31, 0x36, 0x5A}
				ret, _ = DecodeGeneralizedTime(data)
				fmt.Println(ret)
	*/

}
