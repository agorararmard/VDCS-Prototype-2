package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"./vdcs"
)

var pendingGarble = make(map[string]vdcs.CircuitMessage)
var completedGarble = make(map[string]vdcs.GarbledMessage)

var mutex = sync.RWMutex{}

func main() {
	server()
}

func garbleCircuit(ID string) {

	mutex.Lock()
	completedGarble[ID] = vdcs.Garble(pendingGarble[ID])
	//fmt.Println("\n\n\nHere is a completed Garble: ", completedGarble[ID], "\n\n\n")
	mutex.Unlock()

	mutex.Lock()
	delete(pendingGarble, ID)
	mutex.Unlock()
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
		mutex.Lock()
		pendingGarble[x.CID] = x
		mutex.Unlock()
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
		mutex.RLock()
		for _, ok := pendingGarble[x.CID]; ok && (len(pendingGarble) != 0); {
			mutex.RUnlock()
			time.Sleep(10 * time.Microsecond)
			mutex.RLock()
			if _, oke := completedGarble[x.CID]; oke {
				break
			}
			fmt.Println("Trapped in Here!!")
		}
		mutex.RUnlock()

		mutex.RLock()
		value, ok := completedGarble[x.CID]
		mutex.RUnlock()
		if ok {
			responseJSON, err := json.Marshal(value)
			if err != nil {
				fmt.Fprintf(w, "error %s", err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(responseJSON)
			mutex.Lock()
			delete(completedGarble, x.CID)
			mutex.Unlock()
		}
	}
}

func server() {
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/get", getHandler)
	http.ListenAndServe(":8080", nil)
}
