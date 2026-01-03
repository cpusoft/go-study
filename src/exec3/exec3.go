package main

import (
	"fmt"
	"strings"

	"github.com/cpusoft/goutil/executil"
)

func main() {
	/*
	   //p := `10.1.135.22 -p 1-50000`
	   p := `-sV --script mysql-brute -p13308 10.1.135.22 --script-args userdb=./users.txt,passdb=./passwords.txt`
	   params := strings.Split(p, " ")
	   out, err := ExecCommandStdoutPipe("nmap", params, true)
	   fmt.Println("out", out)
	   fmt.Println(err)
	*/
	// 1. 拼接完整的 grep 命令（包含通配符，由 Shell 解析）
	keywords := []string{"RP program", "The", "time"}
	logs := `/root/rpki/rpstir2-rp/log/rpstir2-rp.*`
	var cmds strings.Builder
	cmds.WriteString("grep \"" + keywords[0] + "\" " + logs)
	for i := 1; i < len(keywords); i++ {
		cmds.WriteString(" | grep \"" + keywords[i] + "\"")
	}
	fmt.Println("cmds:", cmds.String())
	out, err := executil.ExecCommandCombinedOutput("bash", []string{"-c", cmds.String()})
	if err != nil {
		fmt.Println("ExecCommandReturnContent fail:", err, " cmds:", cmds)
		return
	}
	fmt.Println("ok cmds:", cmds, "\r\n\r\n", "  out:\r\n", out)

	fmt.Println(err)
}
