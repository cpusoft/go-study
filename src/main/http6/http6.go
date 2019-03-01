package main

import (
	"encoding/json"
	"fmt"
	"github.com/cpusoft/go-json-rest/rest"
	"net"
	"net/http"
	"time"

	belogs "github.com/astaxie/beego/logs"
)

func Lookup(w rest.ResponseWriter, req *rest.Request) {
	fmt.Print(req.PathParam("host"))

	fmt.Print(req.FormValue("key1"))
	fmt.Print(req.FormValue("key2"))

	fmt.Println("Content-Type", req.Header.Get("Content-Type"))
	ip, err := net.LookupIP(req.PathParam("host"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/rpki-slurm")
	w.WriteJson(&ip)
}

/*
	"%b", "{{.BytesWritten | dashIf0}}",
	"%B", "{{.BytesWritten}}",
	"%D", "{{.ResponseTime | microseconds}}",
	"%h", "{{.ApacheRemoteAddr}}",
	"%H", "{{.R.Proto}}",
	"%l", "-",
	"%m", "{{.R.Method}}",
	"%P", "{{.Pid}}",
	"%q", "{{.ApacheQueryString}}",
	"%r", "{{.R.Method}} {{.R.URL.RequestURI}} {{.R.Proto}}",
	"%s", "{{.StatusCode}}",
	"%S", "\033[{{.StatusCode | statusCodeColor}}m{{.StatusCode}}",
	"%t", "{{if .StartTime}}{{.StartTime.Format \"02/Jan/2006:15:04:05 -0700\"}}{{end}}",
	"%T", "{{if .ResponseTime}}{{.ResponseTime.Seconds | printf \"%.3f\"}}{{end}}",
	"%u", "{{.RemoteUser | dashIfEmptyStr}}",
	"%{User-Agent}i", "{{.R.UserAgent | dashIfEmptyStr}}",
	"%{Referer}i", "{{.R.Referer | dashIfEmptyStr}}",
*/
func ListenAndServer(port string, router *rest.App) {
	//logfile, _ := os.OpenFile(, os.O_RDWR|os.O_CREATE, 0666)
	//logger := log.New(logfile, "\r\n", log.Ldate|log.Ltime|log.Llongfile)

	ts := time.Now().Format("2006-01-02")

	logConfig := make(map[string]interface{})
	logConfig["filename"] = `E:\logs\access.log` + "." + ts
	logConfig["level"] = 6
	logConfigStr, _ := json.Marshal(logConfig)
	//fmt.Println("log:logConfigStr", string(logConfigStr))
	belogs.SetLogger(belogs.AdapterFile, string(logConfigStr))

	api := rest.NewApi()
	MyAccessProdStack := rest.AccessProdStack
	MyAccessProdStack[0] = &rest.AccessLogApacheMiddleware{
		Logger: belogs.GetLogger("access"),
		Format: rest.CombinedLogFormat,
	}
	/*
		var MyAccessProdStack = []rest.Middleware{
			&rest.AccessLogApacheMiddleware{
				Logger: belogs.GetLogger("access"),
				Format: rest.CombinedLogFormat,
			},
			&rest.TimerMiddleware{},
			&rest.RecorderMiddleware{},
			&rest.RecoverMiddleware{},
			&rest.GzipMiddleware{},
		}
	*/
	api.Use(MyAccessProdStack...)
	api.SetApp(*router)
	fmt.Println(http.ListenAndServe(port, api.MakeHandler()))
}
func main() {

	router, err := rest.MakeRouter(
		rest.Get("/lookup/#host", Lookup),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	ListenAndServer(":8080", &router)

}
