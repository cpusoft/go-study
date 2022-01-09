package main

import (
	"bytes"
	"encoding/asn1"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/cpusoft/asn1/der"
	"github.com/cpusoft/goutil/fileutil"
)

func main() {

	fn := derFile
	//fn := derHex
	//fn := testUint64
	//fn := testIntDER
	//fn := testIntJSON
	//fn := testPersone
	//fn := testFloat

	if err := fn(); err != nil {
		fmt.Println(err)
	}
}
func derFile() error {
	file := `apnic-rpki-root-iana-origin.cer`
	data1, err := fileutil.ReadFileToBytes(file)
	fmt.Println(len(data1), err)
	n := new(der.Node)

	_, err = der.DecodeNode(data1, n)
	if err != nil {
		return err
	}

	s, err := der.ConvertToString(n)
	if err != nil {
		return err
	}

	fmt.Println(s)
	return nil
}
func derHex() error {

	const hexDump = `30-2E-A0-03-02-01-01-A1 03-02-01-01-A2-03-02-01
01-A3-08-0C-06-31-32-33 34-35-36-A4-13-17-11-31
35-31-32-31-37-31-37-34 38-34-34-2B-30-33-30-30`

	s := onlyHex(hexDump)

	data1, err := hex.DecodeString(s)
	if err != nil {
		return err
	}

	n := new(der.Node)

	_, err = der.DecodeNode(data1, n)
	if err != nil {
		return err
	}

	s, err = der.ConvertToString(n)
	if err != nil {
		return err
	}

	fmt.Println(s)

	data2, err := der.EncodeNode(nil, n)
	if err != nil {
		return err
	}

	fmt.Printf("equal: %t\n", bytes.Equal(data1, data2))

	return nil
}

func byteIsHex(b byte) bool {

	if (b >= '0') && (b <= '9') {
		return true
	}

	if (b >= 'a') && (b <= 'f') {
		return true
	}

	if (b >= 'A') && (b <= 'F') {
		return true
	}

	return false
}

func onlyHex(s string) string {

	data := []byte(s)

	var res []byte
	for _, b := range data {
		if byteIsHex(b) {
			res = append(res, b)
		}
	}

	return string(res)
}

type uint64Sample struct {
	val  uint64
	data []byte
}

func newUint64Sample(v uint64, s string) *uint64Sample {
	s = onlyHex(s)
	data, err := hex.DecodeString(s)
	if err != nil {
		panic(err.Error())
	}
	return &uint64Sample{v, data}
}

func testUint64() error {

	var as = []uint64{0, 1, 2}

	var x uint64 = 4
	for i := 0; i < 64; i++ {

		as = append(as, x-1)
		as = append(as, x)
		as = append(as, x+1)

		x *= 2
	}

	for _, a := range as {

		data, err := der.Marshal(a)
		if err != nil {
			return err
		}

		fmt.Printf("newUint64Sample(%d, \"%X\"),\n", a, data)
	}

	return nil
}

type int64Sample struct {
	val  int64
	data []byte
}

func newInt64Sample(v int64, s string) *int64Sample {
	s = onlyHex(s)
	data, err := hex.DecodeString(s)
	if err != nil {
		panic(err.Error())
	}
	return &int64Sample{v, data}
}

func testInt64() error {

	var as = []int64{0, 1, -1, 2, -2}

	var x int64 = 4
	for i := 0; i < 64; i++ {

		as = append(as, x-1)
		as = append(as, -(x - 1))
		as = append(as, x)
		as = append(as, -x)
		as = append(as, x+1)
		as = append(as, -(x + 1))

		x *= 2
	}

	for _, a := range as {

		data, err := der.Marshal(a)
		if err != nil {
			return err
		}

		fmt.Printf("newInt64Sample(%d, \"%X\"),\n", a, data)
	}

	return nil
}

func testIntDER() error {

	var a int64 = -100000

	data, err := der.Marshal(a)
	if err != nil {
		return err
	}

	fmt.Printf("%X\n", data)

	var b int32

	err = der.Unmarshal(data, &b)
	if err != nil {
		return err
	}

	fmt.Println(b)

	return nil
}

func testIntJSON() error {

	var a int = -108987

	data, err := json.Marshal(a)
	if err != nil {
		return err
	}

	var b int

	err = json.Unmarshal(data, &b)
	if err != nil {
		return err
	}

	fmt.Println(b)

	return nil
}

type Persone struct {
	Name string `asn1:"tag:0" der:"tag:0"`
	Age  int    `asn1:"tag:1" der:"tag:1"`
	Desc string `asn1:"tag:2" der:"tag:2,optional"`
}

func newString(s string) *string {
	return &s
}

func testPersone() error {

	a := Persone{
		Name: "John",
		Age:  -97,
		Desc: "Ароза упала на лапу Азора",
	}

	data, err := der.Marshal(&a)
	if err != nil {
		return err
	}

	fmt.Printf("%X\n", data)

	var b Persone

	err = der.Unmarshal(data, &b)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", b)
	fmt.Println("desc:", b.Desc)

	data, err = asn1.Marshal(a)
	if err != nil {
		return err
	}
	fmt.Printf("%X\n", data)

	return nil
}

func testFloat() error {

	var a float64 = 3.14

	data, err := der.Marshal(a)
	//data, err := asn1.Marshal(a)
	if err != nil {
		return err
	}
	fmt.Printf("%X\n", data)

	return nil
}
