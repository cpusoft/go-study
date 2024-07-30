package main

type CallInterface interface {
	call()
}

var Inf CallInterface

func main() {
	Inf.call()
}
