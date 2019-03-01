package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func sayHelloWorld(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form, "  path", r.URL.Path, r.URL.Scheme, r.Referer(), r.UserAgent(), r.RemoteAddr)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("value:", strings.Join(v, ""))
	}
	fmt.Fprint(w, "hello world")
}

func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("template/login.gtpl")
		if t == nil {
			log.Println("ParseFiles fail")
		}
		log.Println(t)
		mapp := make(map[string]string)
		mapp["username"] = "username input"
		mapp["password"] = "password input"
		err := t.Execute(w, mapp)
		log.Println(err)
	} else {
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
		fmt.Println("[0]", r.Form["username"][0])
		if len(r.Form["username"][0]) == 0 {

		}
	}
}

func main() {
	http.HandleFunc("/", sayHelloWorld)
	http.HandleFunc("/login", login)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
