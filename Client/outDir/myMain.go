package main

import (
	"fmt"
)
import (
"./vdcs"
"strings"
)


func main() {
_myStringMatch_string_8_string_8Ch1:= make(chan vdcs.GarbledCircuit)
go vdcs.Comm("myStringMatch_string_8_string_8",1,_myStringMatch_string_8_string_8Ch1)
_myStringMatch_string_8_string_8Ch0:= make(chan vdcs.GarbledCircuit)
go vdcs.Comm("myStringMatch_string_8_string_8",0,_myStringMatch_string_8_string_8Ch0)
	var s1 string = "Hello World!"
	var s2 string = "lo Wo"

	//VDCS
	if eval0(s1, s2, 0,_myStringMatch_string_8_string_8Ch0) == true {
		fmt.Println("Yes they are!")
	} else {
		fmt.Println("No they aren't!")
	}
	var s3 string = "Hello Earth!"
	//VDCS
	if eval1(s3, s2, 1,_myStringMatch_string_8_string_8Ch1) == true {
		fmt.Println("Yes they are!")
	} else {
		fmt.Println("No they aren't!")
	}
}

func myStringMatch(a, b string) (true bool) {
	return
}

func eval0(s1 string,s2 string,cID int, chVDCSEvalCircRes <-chan vdcs.GarbledCircuit) (bool){
_inWire0:=[]byte(s1)

_inWire1:=[]byte(s2)

	//generate input wires for given inputs
k := <-chVDCSEvalCircRes
k.InWire0 = _inWire0
k.InWire1 = _inWire1
vdcs.SendToServerEval(k)
var res []byte = vdcs.GetFromServerEval(k.CID)
fmt.Println(string(res))
return strings.Contains(string(_inWire0), string(_inWire1))
}

func eval1(s3 string,s2 string,cID int, chVDCSEvalCircRes <-chan vdcs.GarbledCircuit) (bool){
_inWire0:=[]byte(s3)

_inWire1:=[]byte(s2)

	//generate input wires for given inputs
k := <-chVDCSEvalCircRes
k.InWire0 = _inWire0
k.InWire1 = _inWire1
vdcs.SendToServerEval(k)
var res []byte = vdcs.GetFromServerEval(k.CID)
fmt.Println(string(res))
return strings.Contains(string(_inWire0), string(_inWire1))
}
