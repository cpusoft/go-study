package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

type PeopleList struct {
	Name  string `json:"name"`
	Craft string `json:"craft"`
}

type Peoples struct {
	Number  int          `json:"number"`
	Message string       `json:"message"`
	People  []PeopleList `json:"people"`
}
type T struct {
	Return []Desc `json:"return"`
}
type Desc struct {
	Field  []string `json:"field"`
	Start  int64    `json:"start"`
	Token  string   `json:"token"`
	Expire int64    `json:"expire"`
	User   string   `json:"user"`
	Eauth  string   `json:"eauth"`
}
type Person struct {
	Name string
	Age  int
}

func addAge(p Person) {
	_, ok := interface{}(p).(Person)
	if ok {
		p.Age = p.Age + 1
	}
}
func addAgeRef(p *Person) {
	_, ok := interface{}(*p).(Person)
	if ok {
		(*p).Age = (*p).Age + 1
	}
}

type Rectangle struct {
	width, height float64
}
type Circle struct {
	radius float64
}

func add1(a *int) int {
	*a = *a + 1
	return *a
}
func (r Rectangle) area() float64 {
	area := r.width * r.height
	(r).width = 0
	(r).height = 0
	return area
}
func (c Circle) area() float64 {
	area := c.radius * c.radius * math.Pi
	(c).radius = 0
	return area
}
func (r *Rectangle) areaRef() float64 {

	area := (*r).width * (*r).height
	(*r).width = 0
	(*r).height = 0
	return area
}
func (c *Circle) areaRef() float64 {
	area := (*c).radius * (*c).radius * math.Pi
	(*c).radius = 0
	return area
}
func isOdd(i int) bool {
	return i%2 != 0
}
func isEven(i int) bool {
	return i%2 == 0
}

type testInt func(int) bool

func filter(slice []int, f testInt) []int {
	var result []int
	for _, value := range slice {
		if f(value) {
			result = append(result, value)
		}
	}
	return result
}

////////////////////////
type Human struct {
	name  string
	age   int
	phone string
}
type Student struct {
	Human
	school string
	loan   float32
}

type Employee struct {
	Human
	company string
	money   float32
}

func (h *Human) Sing(lyrics string) {
	fmt.Println("human sing.....", lyrics)
}
func (h *Human) SayHi() {
	fmt.Printf("name is %s, phone is %s\n", h.name, h.phone)
}
func (s *Student) SayHi() {
	fmt.Printf("name is %s, phone is %s, study in %s\n", s.name, s.phone, s.school)
}
func (s *Student) Sing(lyrics string) {
	fmt.Println("student sing.....", lyrics)
}
func (s *Student) BorrowMoney(amount float32) {
	s.loan += amount
}
func (e *Employee) SayHi() {
	fmt.Printf("name is %s, phone is %s, work in %s\n", e.name, e.phone, e.company)
}
func (e *Employee) SpendSalary(amount float32) {
	e.money -= amount
}
func (e *Employee) String() string {
	return e.name + " - " + strconv.Itoa(e.age) + " years" + " - " + e.phone
}
func (e Employee) Strings() string {
	return e.name + " - " + strconv.Itoa(e.age) + " years" + " - " + e.phone
}

type HumanInf interface {
	SayHi()
	Sing(lyrics string)
}

type StudentInf interface {
	SayHi()
	Sing(lyrics string)
	BorrowMoney(amount float32)
}

type EmployeeInf interface {
	SayHi()
	Sing(lyrics string)
	SpendSalary(amount float32)
}

func main() {
	hh := Human{name: "human name", phone: "110", age: 22}
	ss := Student{Human: hh, school: "schoolssss", loan: 3333}
	ee := Employee{Human: hh, company: "corp", money: 9999.3}
	xx := make([]HumanInf, 3)
	xx[0], xx[1], xx[2] = &hh, &ss, &ee
	fmt.Println(&ee)
	fmt.Println(ee)
	for _, zz := range xx {
		zz.SayHi()
	}
	var inf HumanInf
	inf = &ss
	inf.Sing("aass")
	inf.SayHi()

	inf = &ee
	inf.Sing("bbb")
	inf.SayHi()

	r1 := Rectangle{height: 33, width: 44}
	c1 := Circle{radius: 32}
	fmt.Println(r1.area())
	fmt.Println(c1.area())
	fmt.Println(r1, c1)

	fmt.Println(r1.areaRef())
	fmt.Println(c1.areaRef())
	fmt.Println(r1, c1)

	p := Person{Name: "sss", Age: 3}
	fmt.Println(p)
	addAge(p)
	fmt.Println(p)
	addAgeRef(&p)
	fmt.Println(p)

	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	odd := filter(slice, isOdd)
	fmt.Println(odd)

	x := 3
	x = add1(&x)
	fmt.Println(x)

	var i int = 99
	const Pi = 3.1415926
	const pref string = "asfsdf"
	var c complex64 = 5 + 99i

	fmt.Println(i, Pi, c)

	s := "ssss"
	bc := []byte(s)
	bc[0] = 'z'
	s2 := string(bc)
	fmt.Println(s2)

	errs := errors.New("test new error")
	if errs != nil {
		fmt.Println(errs)
	}

	const (
		a = iota
		b
		cz
		d
	)
	fmt.Println(d)

	var arr [9]int
	arr[0] = 3
	fmt.Println(arr)
	arr2 := [...]int{3, 4, 5, 2}
	fmt.Println(arr2)

	var sl []byte
	sl = append(sl, 'a')
	fmt.Println("%c", sl[0])

	arr3 := [][]int{{1, 2, 3}, {3, 4, 5}}
	fmt.Println(arr3)

	url := "http://api.open-notify.org/astros.json"

	//var nums map[string]int
	nums := map[string]int{}
	nums2 := map[string]int{"aaa": 1, "bbb": 44, "zzzf": 4}
	result := make(map[string]int)

	result["ff"] = 11
	nums["ff"] = 11
	nums["ff3"] = 113
	fmt.Println(result["ff"], nums["ff"], nums2)
	delete(nums, "ff")
	fmt.Println(nums)

	zz, ok := nums["z"]
	if ok {
		fmt.Println(zz)
	} else {
		fmt.Println("no map")
	}

	zf := new([5]int)
	fmt.Println(zf)
	zf[0] = 3
	fmt.Println(zf)

	buffer := make(map[string][]interface{})
	zzzzz, ok1 := buffer["zzz"]
	if !ok1 {
		fmt.Println("no make map")
	} else {
		fmt.Println(zzzzz)
	}

	sum := 0
	for index := 0; index < 10; index++ {
		sum += index
	}
	fmt.Println("sum is equal to ", sum)

	spaceClient := http.Client{
		Timeout: time.Duration(15 * time.Second),
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "spacecount-tutorial")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	fmt.Printf("HTTP: %s\n", res.Status)

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	body1 := string(body)
	fmt.Println(body1)

	people1 := Peoples{}
	jsonErr := json.Unmarshal(body, &people1)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	fmt.Printf("%v\n", people1.People[0])

	str := `{"return": [  
                    {  
                        "field": [".*"],  
                        "start": "1473841133",  
                        "token": "token1",  
                        "expire": 1473884333,  
                        "user": "xiaochuan",  
                        "eauth": "ss"  
                    },  
                    {  
                        "field": [".*"],  
                        "start": 1473841133,  
                        "token": "token2",  
                        "expire": 1473884333,  
                        "user": "xiaochuan",  
                        "eauth": "sr"  
                    }  
                        ]  
                }`
	t_struct := T{}
	err1 := json.Unmarshal([]byte(str), &t_struct)
	if err1 != nil {
		fmt.Println("error is %v\n", err1)
	} else {
		fmt.Printf("%v\n", t_struct)
	}

}
