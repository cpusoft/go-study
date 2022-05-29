package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	zonefile "github.com/bwesterb/go-zonefile"
)

// Increments the serial of a zonefile
func main() {
	file := `mydomain.com.zone`
	fmt.Println(file)
	// Load zonefile
	data, ioerr := ioutil.ReadFile(file)
	if ioerr != nil {
		fmt.Println(file, ioerr)
		os.Exit(2)
	}
	fmt.Println(len(data))

	zf, perr := zonefile.Load(data)
	if perr != nil {
		fmt.Println(file, perr.LineNo(), perr)
		os.Exit(3)
	}
	fmt.Println(zf)
	fmt.Println(len(zf.Entries()))
	// Find SOA entry

	for i, e := range zf.Entries() {
		fmt.Println(i, e)
		fmt.Println("command:", string(e.Command()))
		fmt.Println("domain:", string(e.Domain()))
		fmt.Println("class:", string(e.Class()))
		fmt.Println("type:", string(e.Type()))
		var sTTL string
		if e.TTL() == nil {
			sTTL = ""
		} else {
			sTTL = strconv.Itoa(*e.TTL())
		}
		fmt.Println("ttl:", sTTL)
		vs := e.Values()
		for j := range vs {
			fmt.Println("value: ", j, string(vs[j]))
		}
		fmt.Println("------")
	}
	/*
		ok := false
		for _, e := range zf.Entries() {
			fmt.Println(e)
			if !bytes.Equal(e.Type(), []byte("SOA")) {
				continue
			}
			vs := e.Values()
			if len(vs) != 7 {
				fmt.Println("Wrong number of parameters to SOA line")
				os.Exit(4)
			}
			serial, err := strconv.Atoi(string(vs[2]))
			if err != nil {
				fmt.Println("Could not parse serial:", err)
				os.Exit(5)
			}
			e.SetValue(2, []byte(strconv.Itoa(serial+1)))
			ok = true
			break
		}
		if !ok {
			fmt.Println("Could not find SOA entry")
			os.Exit(6)
		}

		fh, err := os.OpenFile(file, os.O_WRONLY, 0)
		if err != nil {
			fmt.Println(file, err)
			os.Exit(7)
		}

		_, err = fh.Write(zf.Save())
		if err != nil {
			fmt.Println(file, err)
			os.Exit(8)
		}
	*/
}
