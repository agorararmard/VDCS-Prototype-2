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

var pendingGarble = make(map[string]vdcs.CircuitMessage)
var completedGarble = make(map[string]vdcs.GarbledMessage)

var mutexC = sync.RWMutex{}
var mutexP = sync.RWMutex{}

func main() {
	server()
}

func garbleCircuit(ID string) {
	mutexC.Lock()
	completedGarble[ID] = vdcs.GarbledMessage{
		GarbledCircuit: vdcs.GarbledCircuit{
			ComID: vdcs.ComID{
				CID: ID,
			},
		},
	}
	mutexC.Unlock()

	mutexC.Lock()
	mutexP.RLock()
	completedGarble[ID] = vdcs.Garble(pendingGarble[ID])
	mutexP.RUnlock()
	mutexC.Unlock()

	mutexP.Lock()
	delete(pendingGarble, ID)
	mutexP.Unlock()
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
		fmt.Println("CID:", x.CID)
		mutexP.Lock()
		pendingGarble[x.CID] = x
		mutexP.Unlock()
		fmt.Println(x)
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
		mutexP.RLock()
		for _, ok := pendingGarble[x.CID]; ok; {
		}
		mutexP.RUnlock()

		mutexC.RLock()
		value, ok := completedGarble[x.CID]
		mutexC.RUnlock()
		if ok {
			responseJSON, err := json.Marshal(value)
			if err != nil {
				fmt.Fprintf(w, "error %s", err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(responseJSON)
			mutexC.Lock()
			delete(completedGarble, x.CID)
			mutexC.Unlock()
		}
	}
}

func server() {
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/get", getHandler)
	http.ListenAndServe(":8080", nil)
}
