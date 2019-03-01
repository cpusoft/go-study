package main

import (
	_ "encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "strings"
	_ "sync"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	"gopkg.in/tylerb/graceful.v1"
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

	go func() {
		api1 := rest.NewApi()
		api1.Use(rest.DefaultProdStack...)
		router1, err := rest.MakeRouter(
			rest.Get("/lookup/#host", Lookup),
			rest.Post("/countries", PostCountry),
		)
		if err != nil {
			log.Fatal(err)
		}
		api1.SetApp(router1)
		http.ListenAndServe(":8080", api1.MakeHandler())
	}()
	go func() {
		api := rest.NewApi()
		api.Use(rest.DefaultDevStack...)
		router, err := rest.MakeRouter(
			rest.Get("/message", func(w rest.ResponseWriter, req *rest.Request) {
				for cpt := 1; cpt <= 10; cpt++ {

					// wait 1 second
					time.Sleep(time.Duration(1) * time.Second)

					w.WriteJson(map[string]string{
						"Message": fmt.Sprintf("%d seconds", cpt),
					})
					w.(http.ResponseWriter).Write([]byte("\n"))

					// Flush the buffer to client
					w.(http.Flusher).Flush()
				}
			}),
		)
		if err != nil {
			log.Fatal(err)
		}
		api.SetApp(router)

		server := &graceful.Server{
			Timeout: 10 * time.Second,
			Server: &http.Server{
				Addr:    ":8082",
				Handler: api.MakeHandler(),
			},
		}

		log.Fatal(server.ListenAndServe())
	}()
	go func() {
		api2 := rest.NewApi()
		api2.Use(rest.DefaultDevStack...)
		router2, err := rest.MakeRouter(
			rest.Get("/admin/fetchallroa", FetchAllRoa),
		)
		if err != nil {
			log.Fatal(err)
		}
		api2.SetApp(router2)
		http.ListenAndServe(":8081", api2.MakeHandler())
	}()
	go func() {
		//go run "D:\Program Files\Go\src\crypto\tls\generate_cert.go" --host localhost
		certFile := `E:\Go\go-study\src\main\http5\cert_gobuild.pem`
		keyFile := `E:\Go\go-study\src\main\http5\key_gobuild.pem`

		api3 := rest.NewApi()
		//statusMw = &rest.StatusMiddleware{}
		api3.Use(statusMw)
		api3.Use(rest.DefaultDevStack...)
		router3, err := rest.MakeRouter(
			rest.Get("/https/status", Status),
		)
		if err != nil {
			log.Fatal(err)
		}
		api3.SetApp(router3)
		http.ListenAndServeTLS(":8443", certFile, keyFile, api3.MakeHandler())
	}()

	select {}

}

func Lookup(w rest.ResponseWriter, req *rest.Request) {
	log.Print(req.PathParam("host"))

	log.Print(req.FormValue("key1"))
	log.Print(req.FormValue("key2"))

	log.Println("Content-Type", req.Header.Get("Content-Type"))
	ip, err := net.LookupIP(req.PathParam("host"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/rpki-slurm")
	w.WriteJson(&ip)
}
func FetchAllRoa(w rest.ResponseWriter, req *rest.Request) {
	w.WriteJson("ok")
}

var statusMw = &rest.StatusMiddleware{}

func Status(w rest.ResponseWriter, req *rest.Request) {
	w.WriteJson(statusMw.GetStatus())
}

type Country struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

func PostCountry(w rest.ResponseWriter, r *rest.Request) {
	//country := Country{}
	/*
		e := r.ParseForm()
		log.Println(e)
		log.Println(r.Form)
		log.Println(r.Form["code"])
		log.Println(r.Form["name"])
		log.Println(r.Form["address"])
		w.WriteJson("ok")
	*/
	country := make([]Country, 0)
	err := r.DecodeJsonPayload(&country)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Print(country)
	w.WriteJson(&country)

}
