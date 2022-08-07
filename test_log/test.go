package main

import "fmt"

var b, a int

func aaa() {

	fmt.Println("jaja")

}

func main() {
	for i := 0; i <= 3; i++ {
		go aaa()

	}
	for {
	}
}
