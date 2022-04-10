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
	infos, err := dmi.BIOS()
	for i := range infos {
		fmt.Println(infos[i])
	}

	fmt.Println("BaseBoard")
	infos, err = dmi.BaseBoard()
	for i := range infos {
		fmt.Println(infos[i])
	}

	fmt.Println("Chassis")
	infos, err = dmi.Chassis()
	for i := range infos {
		fmt.Println(infos[i])
	}

	fmt.Println("MemoryArray")
	infos, err = dmi.MemoryArray()
	for i := range infos {
		fmt.Println(infos[i])
	}

	fmt.Println("MemoryDevice")
	infos, err = dmi.MemoryDevice()
	for i := range infos {
		fmt.Println(infos[i])
	}

	fmt.Println("Onboard")
	infos, err = dmi.Onboard()
	for i := range infos {
		fmt.Println(infos[i])
	}

	fmt.Println("PortConnector")
	infos, err = dmi.PortConnector()
	for i := range infos {
		fmt.Println(infos[i])
	}

	fmt.Println("Processor")
	infos, err = dmi.Processor()
	for i := range infos {
		fmt.Println(infos[i])
	}

	fmt.Println("ProcessorCache")
	infos, err = dmi.ProcessorCache()
	for i := range infos {
		fmt.Println(infos[i])
	}

	fmt.Println("Slot")
	infos, err = dmi.Slot()
	for i := range infos {
		fmt.Println(infos[i])
	}

	fmt.Println("System")
	infos, err = dmi.System()
	for i := range infos {
		fmt.Println(infos[i])
	}

	checkError(err)

}
