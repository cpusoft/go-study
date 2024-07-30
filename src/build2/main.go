package main

type CallInterface interface {
	call()
}

func main() {
	inf := CacheDb{}
	inf.call()
}
