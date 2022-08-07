package main

import "fmt"

var b, a int

func aaa() {

	fmt.Println("jaja")

}

func bbb() {
	go aaa()
	fmt.Println("bbb执行结束")
}
func main() {
	bbb()
	fmt.Println("main执行结束")
	for {

	}
}
