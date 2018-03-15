package main

import (
	"fmt"
	"strconv"
)

type Human struct {
	name  string
	age   int
	phone string
}

type Student struct {
	Human  //匿名字段
	school string
	loan   float32
}

type Employee struct {
	Human   //匿名字段
	company string
	money   float32
}

//Human实现SayHi方法
func (h Human) SayHi() {
	fmt.Printf("Hi, I am %s you can call me on %s\n", h.name, h.phone)
}

//Human实现Sing方法
func (h Human) Sing(lyrics string) {
	fmt.Println("La la la la...", lyrics)
}

//Employee重载Human的SayHi方法
func (e Employee) SayHi() {
	fmt.Printf("Hi, I am %s, I work at %s. Call me on %s\n", e.name,
		e.company, e.phone)
}

// Interface Men被Human,Student和Employee实现
// 因为这三个类型都实现了这两个方法
type Men interface {
	SayHi()
	Sing(lyrics string)
}
type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	return "name: " + p.Name + " - age: " + strconv.Itoa(p.Age) + " years"
}

func main() {

	type Element interface{}
	type List []Element
	list := make(List, 3)
	list1 := make([]interface{}, 3)
	list[0] = 1
	list[1] = "hello"
	list[2] = Person{"aaaa", 33}
	list1[0] = 2
	list1[1] = "world"
	list1[2] = Person{"bbbb", 3444}
	fmt.Println(list, list1)
	for index, el := range list {
		if value, ok := el.(int); ok {
			fmt.Printf("list[%d} is an int, is %d\n", index, value)
		} else if value, ok := el.(string); ok {
			fmt.Printf("list[%d} is an string, is %s\n", index, value)
		} else if value, ok := el.(Person); ok {
			fmt.Printf("list[%d} is an person, is %s\n", index, value)
		}
	}
	for index, el := range list1 {
		switch value := el.(type) {
		case int:
			fmt.Printf("list1[%d} is an int, is %d\n", index, value)
		case string:
			fmt.Printf("list1[%d} is an string, is %s\n", index, value)
		case Person:
			fmt.Printf("list1[%d} is an person, is %s\n", index, value)

		}
	}

	var a interface{}
	var ii int = 5
	s := "sssss"
	a = ii
	fmt.Println(a)
	a = s
	fmt.Println(a)

	mike := Student{Human{"Mike", 25, "222-222-XXX"}, "MIT", 0.00}
	paul := Student{Human{"Paul", 26, "111-222-XXX"}, "Harvard", 100}
	sam := Employee{Human{"Sam", 36, "444-222-XXX"}, "Golang Inc.", 1000}
	tom := Employee{Human{"Tom", 37, "222-444-XXX"}, "Things Ltd.", 5000}

	//定义Men类型的变量i
	var i Men

	//i能存储Student
	i = mike
	fmt.Println("This is Mike, a Student:")
	i.SayHi()
	i.Sing("November rain")

	//i也能存储Employee
	i = tom
	fmt.Println("This is tom, an Employee:")
	i.SayHi()
	i.Sing("Born to be wild")

	//定义了slice Men
	fmt.Println("Let's use a slice of Men and see what happens")
	x := make([]Men, 3)
	//这三个都是不同类型的元素，但是他们实现了interface同一个接口
	x[0], x[1], x[2] = paul, sam, mike

	for _, value := range x {
		value.SayHi()
	}

}
