package main

/*
int add(int a,int b){
return a+b;
}
*/
import "C"
import (
	"fmt"
)

func main() {
	fmt.Println(C.add(2, 1))
}
