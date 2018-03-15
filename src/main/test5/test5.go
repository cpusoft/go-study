package main

import (
	"crypto/md5"
	"fmt"

	"io"
	"math/rand"
	"strconv"
	"text/template"
	"time"
	"unicode"
)

type Person struct {
	Name string
	Age  int
}

func personSum(ps []Person, p chan Person) {
	one := Person{"", 0}
	for _, v := range ps {
		one.Age += v.Age
		one.Name += (v.Name + ";")
	}
	p <- one
	fmt.Printf("%+v    ", p)
}
func fib(n int, c chan int) {
	x, y := 1, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	close(c)
}
func randSeq(n int) string {
	//letters := []rune("abcdefghijklmnopqrstuvwxyz")
	letters := []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func fi(c, quit chan int, expire chan bool) {
	x, y := 1, 1
	for {
		select {
		case c <- x:
			fmt.Println("before x:", x, "  y:", y)
			x, y = y, x+y
			fmt.Println("after x:", x, "  y:", y)
			fmt.Println(c)
		case <-time.After(5 * time.Second):
			fmt.Println("timeout")
			expire <- true
			break

		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}
func goRun(c, quit chan int) {
	for i := 0; i < 10; i++ {
		fmt.Println(<-c)
	}
	quit <- 0
}

func main() {
	s := "sssss<script/>"
	ss1 := template.HTMLEscapeString(s)
	fmt.Println(ss1)

	ss1 = template.JSEscapeString(s)
	fmt.Println(ss1)

	h := md5.New()
	io.WriteString(h, strconv.FormatInt(int64(time.Now().Nanosecond()), 10))
	io.WriteString(h, "test")
	token := fmt.Sprintf("%x", h.Sum(nil))
	fmt.Println(token)

	h = md5.New()
	h.Write([]byte("sfdsafdsafdsafds"))
	hmd5 := h.Sum(nil)
	fmt.Printf("%x", hmd5)

	ss := make([]string, 13)
	sss := []string{}
	ss[0] = "aa"
	sss = append(sss, "bbb")
	fmt.Printf("%q\n", ss)
	fmt.Printf("%q\n", sss)

	for _, r := range "Hello 世界！" {
		// 判断字符是否为汉字
		if unicode.Is(unicode.Scripts["Han"], r) {
			fmt.Printf("%c", r) // 世界
		}
	}

	c := make(chan int)
	quit := make(chan int)
	expire := make(chan bool)
	go goRun(c, quit)
	fi(c, quit, expire)
	<-expire

	var ps [20]Person
	for i := 0; i < len(ps); i++ {
		ps[i].Age = rand.Intn(100)
		ps[i].Name = randSeq(5)
		fmt.Printf("%+v\n", ps[i])
	}
	fmt.Printf("%+v\n", ps)

	pm := make(chan Person, 3)
	go personSum(ps[len(ps)/2:], pm)
	go personSum(ps[:len(ps)/2], pm)
	pSum := <-pm
	fmt.Printf("%+v\n", pSum)

	cfb := make(chan int, 10)
	go fib(cap(cfb), cfb)
	for i := range cfb {
		fmt.Println(i)
	}
}
