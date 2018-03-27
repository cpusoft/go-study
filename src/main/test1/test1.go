package main

import (
	"fmt"
	//	"runtime"
	"errors"

	"regexp"
	"strconv"
	"strings"
	"unsafe"
)

const (
	Pi = 3.14
)

type TypeName struct {
	name string
	age  int
}

type Person TypeName

func (p Person) ModifyAge(age int) {
	p.age = age
}
func (p *Person) ModifyAge2(age int) {
	p.age = age
}
func test(s string, x int) (r string) {
	r = fmt.Sprintf("test: %s %d", s, x)
	//	runtime.Breakpoint()
	return r
}
func Division(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("除数不能为0")
	}

	return a / b, nil
}

type ErrCorrupted struct {
	Fd  string
	Err error
}

type Mem struct {
	db  map[string][]byte
	yes bool
}

func main() {

	var sli []int
	fmt.Println("%v", sli)
	sli2 := make([]int, 2)
	sli = append(sli, 1)
	fmt.Println("%v", sli)
	fmt.Println("%v", sli2)

	sli = make([]int, 10)
	for i := 0; i < 10; i++ {
		sli[i] = i
	}
	fmt.Println("111  %v", sli)
	fmt.Println(":0  %v", sli[:0])
	sli = append(sli[:0], sli[1:]...)
	fmt.Println("222 %v", sli)

	person := &Person{"aaa", 12}
	fmt.Println("%+v", person)
	person.ModifyAge(99)
	fmt.Println("%+v", person)
	person.ModifyAge2(99)
	fmt.Println("%+v", person)
	ip := "1.1.2.0"
	formatIp := ""
	ipsV4 := strings.Split(ip, ".")
	if len(ipsV4) > 1 {
		for _, ipV4 := range ipsV4 {
			ip, _ := strconv.Atoi(ipV4)
			formatIp += fmt.Sprintf("%02x", ip)
		}
		fmt.Println(formatIp)
	}

	astruct := make(map[string]struct{})
	fmt.Println(astruct)
	fmt.Println(unsafe.Sizeof(astruct))
	var aaaa struct{}
	astruct["aaa"] = aaaa
	fmt.Println(astruct)
	astruct["struct{}"] = struct{}{}
	fmt.Println(astruct)

	mem := &Mem{
		db:  make(map[string][]byte),
		yes: true}
	fmt.Println(mem)

	mem2 := Mem{
		db:  make(map[string][]byte),
		yes: true}
	fmt.Println(mem2)

	var mem3 Mem
	mem3.db = make(map[string][]byte)
	mem3.yes = false
	fmt.Println(mem3)

	mem4 := new(Mem)
	mem4.db = make(map[string][]byte)
	mem4.yes = true
	fmt.Println(mem4)

	var mem5 *Mem = new(Mem)
	fmt.Println(mem5)

	fmt.Printf("%T\n%T\n%T\n%T\n%T\n%T\n", mem, mem2, mem3, mem4, mem5)

	var st string
	fmt.Println(st, len(st))

	pat6 := `^\s*((([0-9A-Fa-f]{1,4}:){7}(([0-9A-Fa-f]{1,4})|:))|(([0-9A-Fa-f]{1,4}:){6}(:|((25[0-5]|2[0-4]\d|[01]?\d{1,2})(\.(25[0-5]|2[0-4]\d|[01]?\d{1,2})){3})|(:[0-9A-Fa-f]{1,4})))|(([0-9A-Fa-f]{1,4}:){5}((:((25[0-5]|2[0-4]\d|[01]?\d{1,2})(\.(25[0-5]|2[0-4]\d|[01]?\d{1,2})){3})?)|((:[0-9A-Fa-f]{1,4}){1,2})))|(([0-9A-Fa-f]{1,4}:){4}(:[0-9A-Fa-f]{1,4}){0,1}((:((25[0-5]|2[0-4]\d|[01]?\d{1,2})(\.(25[0-5]|2[0-4]\d|[01]?\d{1,2})){3})?)|((:[0-9A-Fa-f]{1,4}){1,2})))|(([0-9A-Fa-f]{1,4}:){3}(:[0-9A-Fa-f]{1,4}){0,2}((:((25[0-5]|2[0-4]\d|[01]?\d{1,2})(\.(25[0-5]|2[0-4]\d|[01]?\d{1,2})){3})?)|((:[0-9A-Fa-f]{1,4}){1,2})))|(([0-9A-Fa-f]{1,4}:){2}(:[0-9A-Fa-f]{1,4}){0,3}((:((25[0-5]|2[0-4]\d|[01]?\d{1,2})(\.(25[0-5]|2[0-4]\d|[01]?\d{1,2})){3})?)|((:[0-9A-Fa-f]{1,4}){1,2})))|(([0-9A-Fa-f]{1,4}:)(:[0-9A-Fa-f]{1,4}){0,4}((:((25[0-5]|2[0-4]\d|[01]?\d{1,2})(\.(25[0-5]|2[0-4]\d|[01]?\d{1,2})){3})?)|((:[0-9A-Fa-f]{1,4}){1,2})))|(:(:[0-9A-Fa-f]{1,4}){0,5}((:((25[0-5]|2[0-4]\d|[01]?\d{1,2})(\.(25[0-5]|2[0-4]\d|[01]?\d{1,2})){3})?)|((:[0-9A-Fa-f]{1,4}){1,2})))|(((25[0-5]|2[0-4]\d|[01]?\d{1,2})(\.(25[0-5]|2[0-4]\d|[01]?\d{1,2})){3})))(%.+)?\s*$`
	url6 := "2001:DB8::" //
	matched6, err6 := regexp.MatchString(pat6, url6)
	fmt.Println(matched6, err6)

	url := "192.168.2.2" //"2001:DB8::" //
	pattern := `^$(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
	pattern = `^(\d+)\.(\d+)\.(\d+)\.(\d+)$`
	matched, err := regexp.MatchString(pattern, url)
	fmt.Println(matched, err)

	ttt := "2001:DB8::" //"FF01::1101" //"ABCD:EF01:2345:6789:ABCD:EF01:2345:6789" //"2001:DB8::"
	count := strings.Count(ttt, ":")
	if count == 7 {

	} else {
		needCount := 7 - count + 2 //2 is current "::"
		fmt.Println(count, needCount)
		colon := strings.Repeat(":", needCount)
		ttt = strings.Replace(ttt, "::", colon, -1)
		fmt.Println(colon, ttt)
	}

	ipsV6 := strings.Split(ttt, ":")
	fmt.Println(len(ipsV6), 9-len(ipsV6))
	formatIps := ""
	for _, ipV6 := range ipsV6 {
		fmt.Println("ipV6 is:", ipV6, len(ipV6))
		if len(ipV6) > 0 {
			formatIps += fmt.Sprintf("%04s", ipV6)
		} else {
			// ipv6 has total 8,  so should use 9
			for i := 0; i < (9 - len(ipsV6)); i++ {
				formatIps += fmt.Sprintf("%04s", ipV6)
			}
		}
	}
	fmt.Println(formatIps)

	lll := strings.Split(ttt, ".")
	fmt.Println(len(lll))

	fmt.Printf("%04s\n", "2000")
	fmt.Printf("%04s\n", "f")

	fmt.Printf("%x\n", 456)
	formatIp = ""
	formatIp = fmt.Sprintf("%02x", 5)
	fmt.Println(formatIp)

	const d = 3e33 / 20
	fmt.Println(d)

	fmt.Printf("hello, world\n")

	s := "haha"
	i := 1234
	fmt.Printf(test(s, i))

	fmt.Println(Pi)

	var tn TypeName
	tn.name = "zzzaaa"
	tn.age = 33
	fmt.Println("aaa" + tn.name)

	tn2 := TypeName{"James", 23}
	tn3 := TypeName{
		name: "sss",
		age:  33,
	}
	fmt.Println(tn2)
	fmt.Println(tn3)
}
