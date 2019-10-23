package main

import (
	"fmt"
)
import (
"sync"
)
var wg = sync.WaitGroup{}


func main() {
wg.Add(2)

go comm1("myStringMatch_string_8_string_4",1)
go comm0("myStringMatch_string_8_string_4",0)
	var s1 string = "Hello World!"
	var s2 string = "lo Wo"
	//VDCS
	if eval0(s1, s2, 0) == true {
		fmt.Println("Yes they are!")
	} else {
		fmt.Println("No they aren't!")
	}
	var s3 string = "Hello World!"
	var s4 string = "lo Wo"
	//VDCS
	if eval1(s3, s4, 1) == true {
		fmt.Println("Yes they are!")
	} else {
		fmt.Println("No they aren't!")
	}
}

func comm0(cir string,cID int) {
fmt.Println(cir)
fmt.Println(cID)
//get the circuit in JSON format
//Generate input wires
//post to server
//Wait for response
wg.Done()}
func eval0(s1 string,s2 string,cID int) (bool){
wg.Wait()
cir := "You did it!"
fmt.Println(cir)
//generate input wires for given inputs
//fetch the garbled circuit with the cID
//post to server
//Wait for response
 return true
}
func comm1(cir string,cID int) {
fmt.Println(cir)
fmt.Println(cID)
//get the circuit in JSON format
//Generate input wires
//post to server
//Wait for response
wg.Done()}
func eval1(s3 string,s4 string,cID int) (bool){
wg.Wait()
cir := "You did it!"
fmt.Println(cir)
//generate input wires for given inputs
//fetch the garbled circuit with the cID
//post to server
//Wait for response
 return true
}