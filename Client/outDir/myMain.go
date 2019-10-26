package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

type ComID struct {
	CID string `json:"key"`
}
type circuit struct {
	O    []bool `json:"o"`
	Feed string `json:"feed"`
	CID  string `json:"key"`
	R    string `json:"randomness"`
}
type GarbledCircuit struct {
	GarbledValues []byte `json:"garbledValues"`
	InWire0       []byte `json:"inWire0"`
	InWire1       []byte `json:"inWire1"`
	ComID
}
type resEval struct {
	Res []byte `json:"res"`
	ComID
}

func main() {
	_myStringMatch_string_8_string_8Ch1 := make(chan GarbledCircuit)
	go comm("myStringMatch_string_8_string_8", 1, _myStringMatch_string_8_string_8Ch1)
	_myStringMatch_string_8_string_8Ch0 := make(chan GarbledCircuit)
	go comm("myStringMatch_string_8_string_8", 0, _myStringMatch_string_8_string_8Ch0)
	var s1 string = "Hello World!"
	var s2 string = "lo Wo"

	//VDCS
	if eval0(s1, s2, 0, _myStringMatch_string_8_string_8Ch0) == true {
		fmt.Println("Yes they are!")
	} else {
		fmt.Println("No they aren't!")
	}
	var s3 string = "Hello Earth!"
	//VDCS
	if eval1(s3, s2, 1, _myStringMatch_string_8_string_8Ch1) == true {
		fmt.Println("Yes they are!")
	} else {
		fmt.Println("No they aren't!")
	}
}

func myStringMatch(a, b string) (true bool) {
	return
}

func eval0(s1 string, s2 string, cID int, chVDCSEvalCircRes <-chan GarbledCircuit) bool {
	_inWire0 := []byte(s1)

	_inWire1 := []byte(s2)

	//generate input wires for given inputs
	k := <-chVDCSEvalCircRes
	k.InWire0 = _inWire0
	k.InWire1 = _inWire1
	sendToServerEval(k)
	var res []byte = getFromServerEval(k.CID)
	fmt.Println(string(res))
	return strings.Contains(string(_inWire0), string(_inWire1))
}

func eval1(s3 string, s2 string, cID int, chVDCSEvalCircRes <-chan GarbledCircuit) bool {
	_inWire0 := []byte(s3)

	_inWire1 := []byte(s2)

	//generate input wires for given inputs
	k := <-chVDCSEvalCircRes
	k.InWire0 = _inWire0
	k.InWire1 = _inWire1
	sendToServerEval(k)
	var res []byte = getFromServerEval(k.CID)
	fmt.Println(string(res))
	return strings.Contains(string(_inWire0), string(_inWire1))
}

func sendToServerGarble(k circuit) bool {
	circuitJSON, err := json.Marshal(k)
	req, err := http.NewRequest("POST", "http://localhost:8080/post", bytes.NewBuffer(circuitJSON))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func sendToServerEval(k GarbledCircuit) {
	circuitJSON, err := json.Marshal(k)
	req, err := http.NewRequest("POST", "http://localhost:8081/post", bytes.NewBuffer(circuitJSON))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()
}

func getFromServerGarble(id string) (k GarbledCircuit) {
	iDJSON, err := json.Marshal(ComID{CID: id})
	req, err := http.NewRequest("GET", "http://localhost:8080/get", bytes.NewBuffer(iDJSON))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &k)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()
	return
}

func getFromServerEval(id string) []byte {
	iDJSON, err := json.Marshal(ComID{CID: id})
	req, err := http.NewRequest("GET", "http://localhost:8081/get", bytes.NewBuffer(iDJSON))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	var k resEval
	err = json.Unmarshal(body, &k)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()
	return k.Res
}

func comm(cir string, cID int, chVDCSCommCircRes chan<- GarbledCircuit) {
	file, _ := ioutil.ReadFile(cir + ".json")
	k := circuit{}
	err := json.Unmarshal([]byte(file), &k)
	if err != nil {
		log.Fatal(err)
	}
	rand.Seed(int64(cID))
	k.CID = strconv.Itoa(rand.Int())
	sendToServerGarble(k)
	//Generate input wires
	//Wait for response
	var g GarbledCircuit = getFromServerGarble(k.CID)
	//Validate Correctness of result
	chVDCSCommCircRes <- g
}
