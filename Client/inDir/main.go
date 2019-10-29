package main

import (
	"fmt"
)

func main() {
	//var s1 string = "Hello World!"
	//var s2 string = "lo Wo"
	var i string = "1"
	var j string = "1"
	//VDCS
	if myEqual(i, j) == true {
		fmt.Println("Yes they are!")
	} else {
		fmt.Println("No they aren't!")
	}
	//var s3 string = "Hello Earth!"
	var z string = "0"
	//VDCS
	if myEqual(j, z) == true {
		fmt.Println("Yes they are!")
	} else {
		fmt.Println("No they aren't!")
	}
}

func myEqual(a, b string) bool {
	return a == b
}
