package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type Recurlyservers struct {
	XMLName     xml.Name `xml:"servers"`
	Version     string   `xml:"version,attr"`
	Svs         []server `xml:"server"`
	Description string   `xml:",innerxml"`
}

type server struct {
	XMLName    xml.Name `xml:"server"`
	ServerName string   `xml:"serverName"`
	ServerIP   string   `xml:"serverIP"`
}

type Servers struct {
	XMLName xml.Name `xml:"servers"`
	Version string   `xml:"version,attr"`
	Svs     []server `xml:"server"`
}

func main() {
	file, err := os.Open("server.xml")
	if err != nil {
		fmt.Printf("open file err: %v\n", err)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("ReadAll file err: %v\n", err)
		return
	}

	v := Recurlyservers{}
	err = xml.Unmarshal(data, &v)
	if err != nil {
		fmt.Printf("ReadAll file err: %v\n", err)
		return
	}
	fmt.Println(v)

	vv := &Servers{Version: "1"}
	vv.Svs = append(vv.Svs, server{ServerName: "Shanghai", ServerIP: "127.0.0.1"})
	vv.Svs = append(vv.Svs, server{ServerName: "Beijing", ServerIP: "127.0.0.2"})
	output, err := xml.MarshalIndent(vv, " ", "	")
	os.Stdout.Write([]byte(xml.Header))
	os.Stdout.Write(output)
}
