package main

import (
	_ "encoding/json"
	_ "fmt"
	"log"
	"net"
	"net/http"
	_ "strings"
	_ "sync"

	"github.com/ant0ine/go-json-rest/rest"
)

func main() {
	/*
		api := rest.NewApi()
		api.Use(rest.DefaultDevStack...)
		api.SetApp(rest.AppSimple(func(w rest.ResponseWriter, r *rest.Request) {
			w.WriteJson(map[string]string{"Body": "Hello World!"})
		}))
		log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
	*/
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/lookup/#host", func(w rest.ResponseWriter, req *rest.Request) {
			log.Print(req.PathParam("host"))
			ip, err := net.LookupIP(req.PathParam("host"))
			if err != nil {
				rest.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteJson(&ip)
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}
