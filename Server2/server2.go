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
var mutex = sync.RWMutex{}

func main() {
	server()
}

func evalCircuit(ID string) {
	mutex.Lock()
	completedEval[ID] = vdcs.ResEval{
		ComID: vdcs.ComID{
			CID: ID,
		},
	}
	mutex.Unlock()

	mutex.Lock()
	fmt.Println("Pending Eval before send: ", pendingEval[ID])
	completedEval[ID] = vdcs.Evaluate(pendingEval[ID])
	fmt.Println("Completed Eval before send: ", completedEval[ID])
	mutex.Unlock()

	mutex.Lock()
	delete(pendingEval, ID)
	mutex.Unlock()
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
		mutex.Lock()
		pendingEval[x.CID] = x
		fmt.Println("Pending Evaluation: ", pendingEval[x.CID])
		mutex.Unlock()
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
		mutex.RLock()
		for _, ok := pendingEval[x.CID]; ok; {
			fmt.Println("Trapped In Here!")
		}
		mutex.RUnlock()

		mutex.RLock()
		value, ok := completedEval[x.CID]
		fmt.Println("Completed Execution: ", completedEval[x.CID])
		mutex.RUnlock()

		if ok {
			responseJSON, err := json.Marshal(value)
			if err != nil {
				fmt.Fprintf(w, "error %s", err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(responseJSON)
			mutex.Lock()
			delete(completedEval, x.CID)
			mutex.Unlock()
		}
	}
}

func server() {
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/get", getHandler)
	http.ListenAndServe(":8081", nil)
}
