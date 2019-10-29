package main

import (
	"fmt"
)
import (
"bytes"
"./vdcs"
)


func main() {
_myEqual_string_1_string_1Ch1:= make(chan vdcs.GarbledMessage)
go vdcs.Comm("myEqual_string_1_string_1",1,_myEqual_string_1_string_1Ch1)
_myEqual_string_1_string_1Ch0:= make(chan vdcs.GarbledMessage)
go vdcs.Comm("myEqual_string_1_string_1",0,_myEqual_string_1_string_1Ch0)
	//var s1 string = "Hello World!"
	//var s2 string = "lo Wo"
	var i string = "1"
	var j string = "1"
	//VDCS
	if eval0(i, j, 0,_myEqual_string_1_string_1Ch0) == true {
		fmt.Println("Yes they are!")
	} else {
		fmt.Println("No they aren't!")
	}
	//var s3 string = "Hello Earth!"
	var z string = "0"
	//VDCS
	if eval1(j, z, 1,_myEqual_string_1_string_1Ch1) == true {
		fmt.Println("Yes they are!")
	} else {
		fmt.Println("No they aren't!")
	}
}

func myEqual(a, b string) bool {
	return a == b
}

func eval0(i string,j string,cID int, chVDCSEvalCircRes <-chan vdcs.GarbledMessage) (bool){
_inWire0:=[]byte(i)

_inWire1:=[]byte(j)

	//generate input wires for given inputs
k := <-chVDCSEvalCircRes
		myInWires := make([]vdcs.Wire, len(_inWire0)*8*2)
for idxByte := 0; idxByte < len(_inWire0); idxByte++ {
for idxBit := 0; idxBit < 8; idxBit++ {
contA := (_inWire0[idxByte] >> idxBit) & 1
myInWires[(idxBit+idxByte*8)*2] = k.InputWires[(idxBit+idxByte*8)*4+int(contA)]
contB := (_inWire1[idxByte] >> idxBit) & 1
myInWires[(idxBit+idxByte*8)*2+1] = k.InputWires[(idxBit+idxByte*8)*4+2+int(contB)]
}
}
/*myInWires := make([]vdcs.Wire, 6)
for idxBit := 0; idxBit < 3; idxBit++ {
contA := (_inWire0[0] >> idxBit) & 1
myInWires[(idxBit)*2] = k.InputWires[(idxBit)*4+int(contA)]
contB := (_inWire1[0] >> idxBit) & 1
myInWires[(idxBit)*2+1] = k.InputWires[(idxBit)*4+2+int(contB)]
}*/
k.InputWires = myInWires//flush output wires
myOutputWires := k.OutputWires
k.OutputWires = []vdcs.Wire{}
	for ok := vdcs.SendToServerEval(k); !ok; {
}
var res [][]byte
var oke bool
for res, oke = vdcs.GetFromServerEval(k.CID); !oke; {
}
//validate and decode res
if bytes.Compare(res[0], myOutputWires[0].WireLabel) == 0 {
return false
} else if bytes.Compare(res[0], myOutputWires[1].WireLabel) == 0 {
return true
} else {
panic("The server cheated me while evaluating")
}
}

func eval1(j string,z string,cID int, chVDCSEvalCircRes <-chan vdcs.GarbledMessage) (bool){
_inWire0:=[]byte(j)

_inWire1:=[]byte(z)

	//generate input wires for given inputs
k := <-chVDCSEvalCircRes
		myInWires := make([]vdcs.Wire, len(_inWire0)*8*2)
for idxByte := 0; idxByte < len(_inWire0); idxByte++ {
for idxBit := 0; idxBit < 8; idxBit++ {
contA := (_inWire0[idxByte] >> idxBit) & 1
myInWires[(idxBit+idxByte*8)*2] = k.InputWires[(idxBit+idxByte*8)*4+int(contA)]
contB := (_inWire1[idxByte] >> idxBit) & 1
myInWires[(idxBit+idxByte*8)*2+1] = k.InputWires[(idxBit+idxByte*8)*4+2+int(contB)]
}
}
/*myInWires := make([]vdcs.Wire, 6)
for idxBit := 0; idxBit < 3; idxBit++ {
contA := (_inWire0[0] >> idxBit) & 1
myInWires[(idxBit)*2] = k.InputWires[(idxBit)*4+int(contA)]
contB := (_inWire1[0] >> idxBit) & 1
myInWires[(idxBit)*2+1] = k.InputWires[(idxBit)*4+2+int(contB)]
}*/
k.InputWires = myInWires//flush output wires
myOutputWires := k.OutputWires
k.OutputWires = []vdcs.Wire{}
	for ok := vdcs.SendToServerEval(k); !ok; {
}
var res [][]byte
var oke bool
for res, oke = vdcs.GetFromServerEval(k.CID); !oke; {
}
//validate and decode res
if bytes.Compare(res[0], myOutputWires[0].WireLabel) == 0 {
return false
} else if bytes.Compare(res[0], myOutputWires[1].WireLabel) == 0 {
return true
} else {
panic("The server cheated me while evaluating")
}
}
