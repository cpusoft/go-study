package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"time"
)

type Person struct {
	Name string
	Age  int
}
type Interface interface {
	sort.Interface
	Push(x interface{})
	Pop() interface{}
}

func say(s string) {
	for i := 0; i < 5; i++ {
		runtime.Gosched()
		fmt.Println(s)
	}
}

func sum(a []int, c chan int) {
	total := 0
	for _, v := range a {
		total += v
	}
	c <- total
}
func sumbuf(a []int, c chan int) {
	total := 0
	for _, v := range a {

		total += v
		fmt.Println(v, total)
	}
	c <- total
}

func fib(n int, c chan int) {
	x, y := 1, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	close(c)
}

func personSum(ps []Person, p chan Person) {
	one := Person{"", 0}
	for _, v := range ps {
		one.Age += v.Age
		one.Name += (v.Name + ";")
	}
	p <- one
}

func main() {
	//golang中格式不是‘yyyy-MM-dd HH:mm:ss’，而是采用golang诞生时间作为格式
	ts := time.Now().Format("2006-01-02")
	fmt.Println(ts)

	var ps [20]Person
	for i := 0; i < len(ps); i++ {
		ps[i].Age = rand.Int()
		f := rand.Float64()
		ps[i].Name = strconv.FormatFloat(f, 'f', 6, 64)
		fmt.Println(ps[i])
	}
	fmt.Println(ps)
	pm := make(chan Person, 3)
	go personSum(ps[len(ps)/2:], pm)
	go personSum(ps[:len(ps)/2], pm)
	pSum := <-pm
	fmt.Println(pSum)

	cfb := make(chan int, 10)
	go fib(cap(cfb), cfb)
	for i := range cfb {
		fmt.Println(i)
	}

	a := []int{3, 4, 5, 6, 2, 4, 3, 3, 2, 22, 5, 43, 3, 4, 2, -45, -6, -7, -4, -5, -6, -7, -8, -3, -3, -4, -4, -33}
	cbuf := make(chan int, 9)
	go sumbuf(a[:len(a)/2], cbuf)
	go sumbuf(a[len(a)/2:], cbuf)

	xx, yy := <-cbuf, <-cbuf
	fmt.Println("xx:", xx, "yy:", yy)

	ci := make(chan int)
	cs := make(chan string)
	cf := make(chan interface{})
	fmt.Println(ci, cs, cf)

	go sum(a, ci)
	x := <-ci
	fmt.Println(x)
	go say("world")
	say("hello")

	//	p := Person{"aaa", 33}
	var p float64 = 33.4
	//t := reflect.TypeOf(p)
	v := reflect.ValueOf(p)
	//fmt.Println("TypeOf  :", t.Elem().Field(0).Tag, "  ", t.Elem().Field(0))
	fmt.Println("ValueOf   type:", v.Type(), "   kind:", v.Kind(), "  value:", v.String())

	vv := reflect.ValueOf(&p)
	pp := vv.Elem()
	pp.SetFloat(999.3)
	fmt.Println(p)

}
