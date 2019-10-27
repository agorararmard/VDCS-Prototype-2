package vdcs

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type Wire struct {
	WireID    string `json:"WireID"`
	WireLabel []byte `json:"WireLabel"`
}
type Gate struct {
	GateID     string   `json:"GateID"`
	GateInputs []string `json:"GateInputs"`
}
type CircuitGate struct {
	Gate
	TruthTable []bool `json:"TruthTable"`
}
type GarbledGate struct {
	Gate
	GarbledValues []byte `json:"GarbledValues"`
}

type ComID struct {
	CID string `json:"ComID"`
}
type Circuit struct {
	InputGates  []CircuitGate `json:"InputGates"`
	MiddleGates []CircuitGate `json:"MiddleGates"`
	OutputGates []CircuitGate `json:"OutputGates"`
}
type CircuitMessage struct {
	Circuit
	ComID
	R []byte `json:"Randomness"`
}
type GarbledCircuit struct {
	InputGates  []GarbledGate `json:"InputGates"`
	MiddleGates []GarbledGate `json:"MiddleGates"`
	OutputGates []GarbledGate `json:"OutputGates"`
	ComID
}

type GarbledMessage struct {
	InputWires []Wire `json:"InputWires"`
	GarbledCircuit
	OutputWires []Wire `json:"OutputGates"`
}

type ResEval struct {
	Res []byte `json:"Result"`
	ComID
}

//basically, the channel will need to send the input/output mapping as well
func Comm(cir string, cID int, chVDCSCommCircRes chan<- GarbledMessage) {
	file, _ := ioutil.ReadFile(cir + ".json")
	k := Circuit{}
	err := json.Unmarshal([]byte(file), &k) //POSSIBLE BUG
	if err != nil {
		log.Fatal(err)
	}
	rand.Seed(int64(cID))
	mCirc := CircuitMessage{Circuit: Circuit{
		InputGates:  k.InputGates,
		MiddleGates: k.MiddleGates,
		OutputGates: k.OutputGates,
	},
		ComID: ComID{strconv.Itoa(rand.Int())},
		R:     []byte("hello world! 3aml eh?"),
	}

	SendToServerGarble(mCirc)
	//Generate input wires
	//Wait for response
	var gcm GarbledMessage = GetFromServerGarble(mCirc.CID)
	//Validate Correctness of result
	chVDCSCommCircRes <- gcm
}

func SendToServerGarble(k CircuitMessage) bool {
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

func SendToServerEval(k GarbledMessage) {
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

func GetFromServerGarble(id string) (k GarbledMessage) {
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

func GetFromServerEval(id string) []byte {
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
