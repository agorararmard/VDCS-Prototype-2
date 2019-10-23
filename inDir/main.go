package main

import (
	"fmt"
)

func main() {
	var s1 string = "Hello World!"
	var s2 string = "lo Wo"
	//VDCS
	if myStringMatch(s1, s2) == true {
		fmt.Println("Yes they are!")
	} else {
		fmt.Println("No they aren't!")
	}
	var s3 string = "Hello World!"
	var s4 string = "lo Wo"
	//VDCS
	if myStringMatch(s3, s4) == true {
		fmt.Println("Yes they are!")
	} else {
		fmt.Println("No they aren't!")
	}
}
