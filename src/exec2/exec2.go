package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/cpusoft/goutil/executil"
)

func ExecCommandStdoutPipe(commandName string, params []string, fmtShow bool) (contentArray []string, err error) {

	var line string
	contentArray = make([]string, 0)
	cmd := exec.Command(commandName, params...)
	//显示运行的命令
	if fmtShow {
		fmt.Printf("exec:%s\n", strings.Join(cmd.Args[:], " "))
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		//belogs.Error("execCommand(): commandName:", commandName, "  err:", err)
		if fmtShow {
			fmt.Fprintln(os.Stderr, "error=>", err.Error())
		}
		return contentArray, err
	}

	cmd.Start()
	reader := bufio.NewReader(stdout)
	for {
		tmp, _, err2 := reader.ReadLine()
		line = string(tmp)
		//line, err2 := reader.ReadString('\n') //[]byte(osutil.GetNewLineSep())[0])
		if err2 != nil || io.EOF == err2 {
			//belogs.Error("execCommand(): ReadString(): line: ", line, "  err2:", err2)
			break
		}
		if fmtShow {
			fmt.Println(line)
		}
		//belogs.Debug("execCommand(): line:", line)
		contentArray = append(contentArray, line)
	}

	cmd.Wait()
	return contentArray, nil
}

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
	grepCmd := `grep "RP program" /root/rpki/rpstir2-rp/log/rpstir2-rp.*`
	// 2. 调用 bash -c 执行完整命令（-c 表示从字符串读取命令）
	// params 格式：[]string{"-c", 完整命令字符串}
	out, err := executil.ExecCommandCombinedOutput("bash", []string{"-c", grepCmd})
	fmt.Println("out", out)
	fmt.Println(err)
}
