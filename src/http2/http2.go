package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	_ "log"
	"net/http"
	_ "strings"

	"github.com/cpusoft/goutil/convert"
)

func login(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("template/login/login.gtpl")
	fmt.Println(t)
	if t == nil {
		fmt.Println("ParseFiles fail")
	}
	err := t.Execute(w, nil)
	fmt.Println(err)
}
func login2(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("template/login/login2.html")
	fmt.Println(t)
	if t == nil {
		fmt.Println("ParseFiles fail")
	}
	err := t.Execute(w, nil)
	fmt.Println(err)
}
func loginQuery(w http.ResponseWriter, r *http.Request) {

	mp := make(map[string]string)
	mp["username"] = "input username"
	mp["password"] = "input password"
	fmt.Println(mp)
	bb, err := json.Marshal(mp)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bb))
	w.Write(bb)
}

func loginSubmit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	u := r.Form["u"]
	p := r.Form["p"]

	mp := make(map[string]string)
	mp["result"] = "ok"
	mp[convert.ToString(u)] = convert.ToString(p)
	fmt.Println(mp)
	bb, _ := json.Marshal(mp)
	fmt.Println(string(bb))
	w.Write(bb)
}
func postCountry(w http.ResponseWriter, r *http.Request) {
	//country := Country{}
	e := r.ParseForm()
	fmt.Println(e)
	fmt.Println(r.Form)
	fmt.Println(r.Form["code"])
	fmt.Println(r.Form["name"])
	fmt.Println(r.Form["address"])
	w.Write([]byte("ok"))
	/*
		country := make([]Country, 0)
		err := r.DecodeJsonPayload(&country)
		if err != nil {
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Print(country)
		w.WriteJson(&country)
	*/
}

func main() {
	http.Handle("/static/", http.FileServer(http.Dir("template")))

	http.HandleFunc("/login", login)
	http.HandleFunc("/login2", login2)
	http.HandleFunc("/loginQuery", loginQuery)
	http.HandleFunc("/loginSubmit", loginSubmit)
	http.HandleFunc("/countries", postCountry)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("ListenAndServe:", err)
	}

}
