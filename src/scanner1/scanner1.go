package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	expect "github.com/Netflix/go-expect"
)

func main() {
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	fmt.Println("你输入的是：", input.Text())

	c, err := expect.NewConsole(expect.WithStdout(os.Stdout))
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	cmd := exec.Command("dir")
	cmd.Stdin = c.Tty()
	cmd.Stdout = c.Tty()
	cmd.Stderr = c.Tty()

	go func() {
		c.ExpectEOF()
	}()

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second)
	c.Send("Hello\r")
	time.Sleep(time.Second)
	c.Send("world\r")

	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
