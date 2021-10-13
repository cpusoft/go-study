package main

//	"os"
//	"syscall"

//	belogs "github.com/beego/beego/v2/core/logs"

//var globalFile *os.File

func main() {
	/*
		logFile, err := os.OpenFile("./log/fatal.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
		if err != nil {
			belogs.Info("服务启动出错", "打开异常日志文件失败", err)
			return
		}

		globalFile = belogs.GetLogger().
		// 将进程标准出错重定向至文件，进程崩溃时运行时将向该文件记录协程调用栈信息
		err = syscall.Dup2(int(logFile.Fd()), int(os.Stderr.Fd()))
		if err != nil {
			return
		}
		return
	*/
}
