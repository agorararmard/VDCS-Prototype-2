package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"./vdcs"
)

var pendingEval = make(map[string]vdcs.GarbledMessage)
var completedEval = make(map[string]vdcs.ResEval)
var mutexP = sync.RWMutex{}
var mutexC = sync.RWMutex{}

func main() {
	server()
}

func evalCircuit(ID string) {
	mutexC.Lock()
	completedEval[ID] = vdcs.ResEval{
		ComID: vdcs.ComID{
			CID: ID,
		},
	}
	mutexC.Unlock()

	mutexC.Lock()
	mutexP.RLock()
	completedEval[ID] = vdcs.Evaluate(pendingEval[ID])
	mutexP.RUnlock()
	mutexC.Unlock()

	mutexP.Lock()
	delete(pendingEval, ID)
	mutexP.Unlock()
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
		mutexP.Lock()
		pendingEval[x.CID] = x
		fmt.Println("Pending Evaluation: ", pendingEval[x.CID])
		mutexP.Unlock()
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
		mutexP.RLock()
		for _, ok := pendingEval[x.CID]; ok; {
		}
		mutexP.RUnlock()

		mutexC.RLock()
		value, ok := completedEval[x.CID]
		fmt.Println("Completed Execution: ", completedEval[x.CID])
		mutexC.RUnlock()

		if ok {
			responseJSON, err := json.Marshal(value)
			if err != nil {
				fmt.Fprintf(w, "error %s", err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(responseJSON)
			mutexC.Lock()
			delete(completedEval, x.CID)
			mutexC.Unlock()
		}
	}
}

func server() {
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/get", getHandler)
	http.ListenAndServe(":8081", nil)
}
