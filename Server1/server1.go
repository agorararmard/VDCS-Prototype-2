package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"./vdcs"
)

var pendingGarble = make(map[string]vdcs.CircuitMessage)
var completedGarble = make(map[string]vdcs.GarbledMessage)

func main() {
	server()
}

func garbleCircuit(ID string) {
	completedGarble[ID] = vdcs.GarbledMessage{
		GarbledCircuit: vdcs.GarbledCircuit{
			ComID: vdcs.ComID{
				CID: ID,
			},
		},
	}
	delete(pendingGarble, ID)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var x vdcs.CircuitMessage
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
		var x vdcs.ComID
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
			responseJSON, err := json.Marshal(value)
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
