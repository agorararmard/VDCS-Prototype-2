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

var pendingEval map[string]GarbledCircuit
var completedEval map[string]ResEval

func main() {
	server()
}

func evalCircuit(ID string) {
	completedEval[ID] = ResEval{
		Res:   []byte("You did it!"),
		comID: comID{CID: ID},
	}
	delete(pendingEval, ID)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var x GarbledCircuit
		jsn, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal("Error reading", err)
		}
		err = json.Unmarshal(jsn, &x)
		if err != nil {
			log.Fatal("bad decode", err)
		}
		pendingEval[x.CID] = x
		go evalCircuit(x.CID)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var x comID
		jsn, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal("Error reading", err)
		}
		err = json.Unmarshal(jsn, &x)
		if err != nil {
			log.Fatal("bad decode", err)
		}
		for _, ok := pendingEval[x.CID]; ok; {
		}
		value, ok := completedEval[x.CID]

		if ok {
			response := value
			responseJSON, err := json.Marshal(response)
			if err != nil {
				fmt.Fprintf(w, "error %s", err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(responseJSON)
			delete(completedEval, x.CID)
		}
	}
}

func server() {
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/get", getHandler)
	http.ListenAndServe(":8081", nil)
}