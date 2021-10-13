package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	belogs "github.com/beego/beego/v2/core/logs"
	"github.com/cpusoft/go-json-rest/rest"
	"github.com/cpusoft/goutil/httpserver"
	"github.com/cpusoft/goutil/osutil"
)

type HttpResponse struct {
	Result string `json:"result"`
	Msg    string
}
type Country struct {
	Code string
	Name string
}

func main() {
	router, err := rest.MakeRouter(
		rest.Post("/uploadfile/#userid", uploadFile),
	)

	if err != nil {
		belogs.Error("main(): failed: err:", err)
		return
	}
	// if have http port, then sart http server, default is off
	httpserver.ListenAndServe(":9181", &router)
}

// vc fetch all rtr data from rp
func uploadFile(w rest.ResponseWriter, req *rest.Request) {

	belogs.Info("uploadFile(): start ")
	userId := req.PathParam("userid")
	fmt.Println("userId:", userId)

	receiveFiles, err := httpserver.ReceiveFiles("./", req)
	fmt.Println("receiveFiles:", receiveFiles, "   err:", err)
	w.WriteJson(HttpResponse{Result: "ok", Msg: "ok"})
}

// return: map[fileFormName]=fileName, such as map["file1"]="aabbccdd.txt"
func ReceiveFiles(receiveDir string, r *rest.Request) (receiveFiles map[string]string, err error) {
	//belogs.Debug("ReceiveFiles(): receiveDir:", receiveDir)
	defer r.Body.Close()

	reader, err := r.MultipartReader()
	if err != nil {
		belogs.Error("ReceiveFiles(): err:", err)
		return nil, err
	}
	receiveFiles = make(map[string]string)
	for {
		part, err := reader.NextPart()
		if err == io.EOF || part == nil {
			break
		}
		if !strings.HasSuffix(receiveDir, string(os.PathSeparator)) {
			receiveDir = receiveDir + string(os.PathSeparator)
		}
		file := receiveDir + osutil.Base(part.FileName())
		form := strings.TrimSpace(part.FormName())
		belogs.Debug("ReceiveFiles():FileName:", part.FileName(), "   FormName:", part.FormName()+"   file:", file)
		if part.FileName() == "" { // this is FormData
			data, _ := ioutil.ReadAll(part)
			ioutil.WriteFile(file, data, 0644)
		} else { // This is FileData
			dst, _ := os.Create(file)
			defer dst.Close()
			io.Copy(dst, part)
		}
		receiveFiles[form] = file
	}
	return receiveFiles, nil
}
