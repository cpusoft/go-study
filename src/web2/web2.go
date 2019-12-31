package main

import (
	"html/template"
	"os"
)

type Person struct {
	UserName string
}

func main() {
	t := template.New("filename")
	t, _ = t.Parse("hello, {{.UserName}}!")
	p := Person{UserName: "name"}
	t.Execute(os.Stdout, p)
}
