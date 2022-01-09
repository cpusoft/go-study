package main

import (
	"fmt"
	"strings"

	conf "github.com/cpusoft/goutil/conf"
)

func main() {
	key := "slurm::slurmServerHttpPort"
	s := conf.String(key)
	fmt.Println("k:", key, "s:", s)
	s = conf.VariableString("slurm::slurmServerHttpPort")
	fmt.Println(s)

	s = conf.VariableString("slurm::slurmPath")
	fmt.Println(s)

}
func VariableString(key string) string {
	if len(key) == 0 {
		fmt.Println(key)

		return ""
	}
	value := `${rpstir2-rp::serverHttpPort}`
	start := strings.Index(value, "${")
	end := strings.Index(value, "}")
	fmt.Println(key, value)
	fmt.Println(start)
	fmt.Println(end)
	if start >= 0 && end > 0 && start < end {
		//${rpstir2::datadir}/rsyncrepo -->rpstir2::datadir
		replaceKey := string(value[start+len("${") : end])
		if len(replaceKey) == 0 || len(conf.String(replaceKey)) == 0 {
			return value
		}
		//rpstir2::datadir -->get  "/root/rpki/data"
		replaceValue := conf.String(replaceKey)
		prefix := string(value[:start])
		suffix := ""
		if end+1 < len(value) {
			suffix = string(value[end+1:])
		}
		///root/rpki/data/rsyncrepo
		newValue := prefix + replaceValue + suffix
		return newValue
	}
	return ""

}
