package main

import (
	"fmt"
)

func main() {
	_myStringMatch_string_8_string_8Ch1 := make(chan int)
	go comm1("myStringMatch_string_8_string_8", 1, _myStringMatch_string_8_string_8Ch1)
	_myStringMatch_string_8_string_8Ch0 := make(chan int)
	go comm0("myStringMatch_string_8_string_8", 0, _myStringMatch_string_8_string_8Ch0)
	var s1 string = "Hello World!"
	var s2 string = "lo Wo"
	//VDCS
	if eval0(s1, s2, 0, _myStringMatch_string_8_string_8Ch0) == true {
		fmt.Println("Yes they are!")
	} else {
		fmt.Println("No they aren't!")
	}
	var s3 string = "Hello World!"
	var s4 string = "lo Wo"
	//VDCS
	if eval1(s3, s4, 1, _myStringMatch_string_8_string_8Ch1) == true {
		fmt.Println("Yes they are!")
	} else {
		fmt.Println("No they aren't!")
	}
}

func myStringMatch(a, b string) bool {
	return true
}

func comm0(cir string, cID int, chVDCSCommCircRes chan<- int) {
	fmt.Println(cir)
	fmt.Println(cID)
	//get the circuit in JSON format
	//Generate input wires
	//post to server
	//Wait for response
	chVDCSCommCircRes <- 32
}
func eval0(s1 string, s2 string, cID int, chVDCSEvalCircRes <-chan int) bool {
	_inWire0 := []byte(s1)

	_inWire1 := []byte(s2)

	i := <-chVDCSEvalCircRes
	fmt.Println(i)
	cir := "You did it!"
	fmt.Println(cir)
	fmt.Println(_inWire0)
	fmt.Println(_inWire1)
	//generate input wires for given inputs
	//fetch the garbled circuit with the cID
	//post to server
	//Wait for response
	return true
}
func comm1(cir string, cID int, chVDCSCommCircRes chan<- int) {
	fmt.Println(cir)
	fmt.Println(cID)
	//get the circuit in JSON format
	//Generate input wires
	//post to server
	//Wait for response
	chVDCSCommCircRes <- 32
}
func eval1(s3 string, s4 string, cID int, chVDCSEvalCircRes <-chan int) bool {
	_inWire0 := []byte(s3)

	_inWire1 := []byte(s4)

	i := <-chVDCSEvalCircRes
	fmt.Println(i)
	cir := "You did it!"
	fmt.Println(cir)
	fmt.Println(_inWire0)
	fmt.Println(_inWire1)
	//generate input wires for given inputs
	//fetch the garbled circuit with the cID
	//post to server
	//Wait for response
	return true
}
