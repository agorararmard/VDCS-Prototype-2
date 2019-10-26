package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

var pendingGarble = make(map[string]Circuit)
var completedGarble = make(map[string]GarbledCircuit)

func main() {
	server()
}

func garbleCircuit(ID string) {
	completedGarble[ID] = GarbledCircuit{
		GarbledValues: []byte("Hello World"),
		ComID:         ComID{CID: ID},
	}
	delete(pendingGarble, ID)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var x Circuit
		jsn, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Fatal("Error reading", err)
		}

		err = json.Unmarshal(jsn, &x)

		if err != nil {
			log.Fatal("bad decode", err)
		}

		pendingGarble[x.CID] = x
		go garbleCircuit(x.CID)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var x ComID
		jsn, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal("Error reading", err)
		}
		err = json.Unmarshal(jsn, &x)
		if err != nil {
			log.Fatal("bad decode", err)
		}

		for _, ok := pendingGarble[x.CID]; ok; {
		}

		value, ok := completedGarble[x.CID]
		if ok {
			response := value
			responseJSON, err := json.Marshal(response)
			if err != nil {
				fmt.Fprintf(w, "error %s", err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(responseJSON)
			delete(completedGarble, x.CID)
		}
	}
}

func server() {
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/get", getHandler)
	http.ListenAndServe(":8080", nil)
}
