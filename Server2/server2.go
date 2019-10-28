package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"./vdcs"
)

var pendingEval map[string]vdcs.GarbledMessage
var completedEval map[string]vdcs.ResEval

func main() {
	server()
}

func evalCircuit(ID string) {
	completedEval[ID] = vdcs.ResEval{
		ComID: vdcs.ComID{
			CID: ID,
		},
	}
	completedEval[ID] = vdcs.Evaluate(pendingEval[ID])
	delete(pendingEval, ID)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var x vdcs.GarbledMessage
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
		var x vdcs.ComID
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
			responseJSON, err := json.Marshal(value)
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
