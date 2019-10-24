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

type circuit struct {
	O    []bool `json:"o"`
	Feed string `json:"feed"`
	CID  string `json:"key"`
}

func main() {
	_myStringMatch_string_8_string_8Ch1 := make(chan circuit)
	go comm1("myStringMatch_string_8_string_8", 1, _myStringMatch_string_8_string_8Ch1)
	_myStringMatch_string_8_string_8Ch0 := make(chan circuit)
	go comm0("myStringMatch_string_8_string_8", 0, _myStringMatch_string_8_string_8Ch0)
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

func comm0(cir string, cID int, chVDCSCommCircRes chan<- circuit) {
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
	k = getFromServerGarble(k.CID)
	chVDCSCommCircRes <- k
}

func eval0(s1 string, s2 string, cID int, chVDCSEvalCircRes <-chan circuit) bool {
	_inWire0 := []byte(s1)

	_inWire1 := []byte(s2)

	//generate input wires for given inputs
	k := <-chVDCSEvalCircRes
	sendToServerEval(k, _inWire0, _inWire1)
	var res []byte = getFromServerEval(k.CID)
	fmt.Println(res)
	return strings.Contains(string(_inWire0), string(_inWire1))
}

func comm1(cir string, cID int, chVDCSCommCircRes chan<- circuit) {
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
	k = getFromServerGarble(k.CID)
	chVDCSCommCircRes <- k
}

func eval1(s3 string, s2 string, cID int, chVDCSEvalCircRes <-chan circuit) bool {
	_inWire0 := []byte(s3)

	_inWire1 := []byte(s2)

	//generate input wires for given inputs
	k := <-chVDCSEvalCircRes
	sendToServerEval(k, _inWire0, _inWire1)
	var res []byte = getFromServerEval(k.CID)
	fmt.Println(res)
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

func sendToServerEval(k circuit, inWire0 []byte, inWire1 []byte) {
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

func getFromServerGarble(id string) (k circuit) {
	iDJSON, err := json.Marshal(circuit{CID: id})
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
	iDJSON, err := json.Marshal(circuit{CID: id})
	req, err := http.NewRequest("GET", "http://localhost:8081/get", bytes.NewBuffer(iDJSON))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()
	return body
}
