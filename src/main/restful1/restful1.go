package main

import (
	"fmt"

	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s\n", ps.ByName("name"))
}
func Adduser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	fmt.Fprintf(w, "Adduser, %s\n", ps.ByName("name"))
}
func main() {

	router := httprouter.New()
	router.GET("/Index", Index)
	router.GET("/hello/:name", Hello)

	router.POST("/adduser/:uid", Adduser)

	http.ListenAndServe(":8080", router)
}
