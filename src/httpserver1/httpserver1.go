package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"

	belogs "github.com/beego/beego/v2/core/logs"
	"github.com/cpusoft/go-json-rest/rest"
	"github.com/cpusoft/goutil/httpserver"
)

type HttpResponse struct {
	Result string `json:"result"`
	Msg    string
}

func main() {
	router, err := rest.MakeRouter(
		rest.Post("/upload", uploadFile),
	)

	if err != nil {
		belogs.Error("main(): failed: err:", err)
		return
	}
	// if have http port, then sart http server, default is off
	httpserver.ListenAndServe(":8080", &router)
}

// vc fetch all rtr data from rp
func uploadFile(w rest.ResponseWriter, req *rest.Request) {

	belogs.Info("uploadFile(): start ")

	file, handler, err := req.FormFile("file")
	if err != nil {
		belogs.Error("uploadFile(): FormFile: err:", err)
		rest.Error(w, err.Error(), http.StatusNoContent)
		return
	}
	defer file.Close()
	tmpFile, _ := ioutil.TempFile("", handler.Filename+"-*.tmp")
	f, err := os.OpenFile(tmpFile.Name(), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		belogs.Error("uploadFile(): OpenFile: err:", err)
		rest.Error(w, err.Error(), http.StatusNoContent)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	w.WriteJson(HttpResponse{Result: "ok", Msg: "ok"})
}
