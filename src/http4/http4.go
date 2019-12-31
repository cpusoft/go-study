package main

import (
	_ "encoding/json"
	_ "fmt"
	"log"
	_ "net"
	"net/http"
	_ "strings"
	"sync"

	"github.com/ant0ine/go-json-rest/rest"
)

func main() {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/countries", GetAllCountries),
		rest.Post("/countries", PostCountry),
		rest.Get("/countries/:code", GetCountry),
		rest.Delete("/countries/:code", DeleteCountry),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}

type Country struct {
	Code string
	Name string
}

var store = map[string]*Country{}
var lock = sync.RWMutex{}

func GetCountry(w rest.ResponseWriter, r *rest.Request) {
	code := r.PathParam("code")

	lock.RLock()
	var country *Country
	if store[code] != nil {
		country = &Country{}
		*country = *store[code]
	}
	lock.RUnlock()

	if country == nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(country)
}

func GetAllCountries(w rest.ResponseWriter, r *rest.Request) {
	lock.RLock()
	countries := make([]Country, len(store))
	i := 0
	for _, country := range store {
		countries[i] = *country
		i++
	}
	lock.RUnlock()
	w.WriteJson(&countries)
}

func PostCountry(w rest.ResponseWriter, r *rest.Request) {
	country := Country{}
	err := r.DecodeJsonPayload(&country)
	log.Println(country)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println(err)

	if country.Code == "" {
		rest.Error(w, "country code required", 400)
		return
	}
	log.Println(country.Code)

	if country.Name == "" {
		rest.Error(w, "country name required", 400)
		return
	}
	log.Println(country.Name)

	lock.Lock()
	store[country.Code] = &country
	log.Println(store)

	lock.Unlock()
	log.Println("WriteJson")
	w.WriteJson(store)
}

func DeleteCountry(w rest.ResponseWriter, r *rest.Request) {
	code := r.PathParam("code")
	lock.Lock()
	delete(store, code)
	lock.Unlock()
	w.WriteHeader(http.StatusOK)
}
