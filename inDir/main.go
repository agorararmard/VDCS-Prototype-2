package main

import (
	"fmt"
)

To making it actually do func main() {
	var s1 string = "Hello World!"
	var s2 string = "lo Wo"
	//VDCS
	if myStringMatch(s1, s2) == true {
		fmt.Println("Yes they are!")
	} else {
		fmt.Println("No they aren't!")
	}
	var s3 string = "Hello Earth!"
	//VDCS
	if myStringMatch(s3, s2) == true {
		fmt.Println("Yes they are!")
	} else {
		fmt.Println("No they aren't!")
	}
}

func myStringMatch(a, b string) bool {
	return true
}
