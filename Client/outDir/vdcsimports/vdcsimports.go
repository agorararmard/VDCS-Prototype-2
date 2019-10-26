package vdcsimports

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type ComID struct {
	CID string `json:"key"`
}
type Circuit struct {
	O    []bool `json:"o"`
	Feed string `json:"feed"`
	ComID
	R string `json:"randomness"`
}
type GarbledCircuit struct {
	GarbledValues []byte `json:"garbledValues"`
	InWire0       []byte `json:"inWire0"`
	InWire1       []byte `json:"inWire1"`
	ComID
}
type ResEval struct {
	Res []byte `json:"res"`
	ComID
}

func comm(cir string, cID int, chVDCSCommCircRes chan<- GarbledCircuit) {
	file, _ := ioutil.ReadFile(cir + ".json")
	k := Circuit{}
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

func sendToServerGarble(k Circuit) bool {
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
	var k ResEval
	err = json.Unmarshal(body, &k)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()
	return k.Res
}
