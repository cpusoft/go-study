package main

import (
	"fmt"
	"os"

	"github.com/yumaojun03/dmidecode"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	dmi, err := dmidecode.New()
	checkError(err)

	fmt.Println("BIOS")
	infos1, err := dmi.BIOS()
	for i := range infos1 {
		fmt.Println(infos1[i])
	}

	fmt.Println("BaseBoard")
	infos2, err := dmi.BaseBoard()
	for i := range infos2 {
		fmt.Println(infos2[i])
	}

	fmt.Println("Chassis")
	infos3, err := dmi.Chassis()
	for i := range infos3 {
		fmt.Println(infos3[i])
	}

	fmt.Println("MemoryArray")
	infos4, err := dmi.MemoryArray()
	for i := range infos4 {
		fmt.Println(infos4[i])
	}

	fmt.Println("MemoryDevice")
	infos5, err := dmi.MemoryDevice()
	for i := range infos5 {
		fmt.Println(infos5[i])
	}

	fmt.Println("Onboard")
	infos6, err := dmi.Onboard()
	for i := range infos6 {
		fmt.Println(infos6[i])
	}

	fmt.Println("PortConnector")
	infos7, err := dmi.PortConnector()
	for i := range infos7 {
		fmt.Println(infos7[i])
	}

	fmt.Println("Processor")
	infos8, err := dmi.Processor()
	for i := range infos8 {
		fmt.Println(infos8[i])
	}

	fmt.Println("ProcessorCache")
	infos9, err := dmi.ProcessorCache()
	for i := range infos9 {
		fmt.Println(infos9[i])
	}

	fmt.Println("Slot")
	infos10, err := dmi.Slot()
	for i := range infos10 {
		fmt.Println(infos10[i])
	}

	fmt.Println("System")
	infos11, err := dmi.System()
	for i := range infos11 {
		fmt.Println(infos11[i])
	}

	checkError(err)

}
