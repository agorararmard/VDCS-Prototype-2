package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var goCnt int

var supportedFunc [1]string = [1]string{"myStringMatch"}

var mapImports map[string]bool = map[string]bool{
	"fmt":           false,
	"strings":       false,
	"net/http":      false,
	"bytes":         false,
	"encoding/json": false,
	"io/ioutil":     false,
	"log":           false,
	"math/rand":     false,
	"strconv":       false}

const circuitBlock string = "type circuit struct {\nO    []bool `json:\"o\"`\nFeed string `json:\"feed\"`\nCID  string `json:\"key\"`\n}\n"

//const commBlock string = "func comm(cir string,cID int, chVDCSCommCircRes chan<- circuit) {fmt.Println(cir)\nfmt.Println(cID)\n//get the circuit in JSON format\n//Generate input wires\n//post to server\n//Wait for response\nchVDCSCommCircRes<-32\n}"
const commBlock string = "func comm(cir string,cID int, chVDCSCommCircRes chan<- circuit) {file, _ := ioutil.ReadFile(cir + \".json\")\nk := circuit{}\nerr := json.Unmarshal([]byte(file), &k)\nif err != nil {\nlog.Fatal(err)\n}\nrand.Seed(int64(cID))\nk.CID = strconv.Itoa(rand.Int())\nsendToServerGarble(k)\n//Generate input wires\n//Wait for response\nk = getFromServerGarble(k.CID)\nchVDCSCommCircRes <- k\n}\n"
const evalBlock string = "func evalcID int, chVDCSEvalCircRes <-chan circuit) (bool){\n	//generate input wires for given inputs\nk := <-chVDCSEvalCircRes\nsendToServerEval(k, _inWire0, _inWire1)\nvar res []byte = getFromServerEval(k.CID)\nfmt.Println(res)\nreturn strings.Contains(string(_inWire0), string(_inWire1))\n}\n"

const sendToGarbleBlock string = "func sendToServerGarble(k circuit) bool {\ncircuitJSON, err := json.Marshal(k)\nreq, err := http.NewRequest(\"POST\", \"http://localhost:8080/post\", bytes.NewBuffer(circuitJSON))\nreq.Header.Set(\"Content-Type\", \"application/json\")\nclient := &http.Client{}\nresp, err := client.Do(req)\nresp.Body.Close()\nif err != nil {\nlog.Fatal(err)\nreturn false\n}\nreturn true\n}\n"
const getFromGarbleBlock string = "func getFromServerGarble(id string) (k circuit) {\niDJSON, err := json.Marshal(circuit{CID: id})\nreq, err := http.NewRequest(\"GET\", \"http://localhost:8080/get\", bytes.NewBuffer(iDJSON))\nreq.Header.Set(\"Content-Type\", \"application/json\")\nclient := &http.Client{}\nresp, err := client.Do(req)\nif err != nil {\nlog.Fatal(err)\n}\nbody, err := ioutil.ReadAll(resp.Body)\nerr = json.Unmarshal(body, &k)\nif err != nil {\nlog.Fatal(err)\n}\nresp.Body.Close()\nreturn\n}\n"

const sendToEvalBlock string = "func sendToServerEval(k circuit, inWire0 []byte, inWire1 []byte) {\ncircuitJSON, err := json.Marshal(k)\nreq, err := http.NewRequest(\"POST\", \"http://localhost:8081/post\", bytes.NewBuffer(circuitJSON))\nreq.Header.Set(\"Content-Type\", \"application/json\")\nclient := &http.Client{}\nresp, err := client.Do(req)\nif err != nil {\nlog.Fatal(err)\n}\nresp.Body.Close()\n}\n"
const getFromEvalBlock string = "func getFromServerEval(id string) []byte {\niDJSON, err := json.Marshal(circuit{CID: id})\nreq, err := http.NewRequest(\"GET\", \"http://localhost:8081/get\", bytes.NewBuffer(iDJSON))\nreq.Header.Set(\"Content-Type\", \"application/json\")\nclient := &http.Client{}\nresp, err := client.Do(req)\nif err != nil {\nlog.Fatal(err)\n}\nbody, err := ioutil.ReadAll(resp.Body)\nif err != nil {\nlog.Fatal(err)\n}\nresp.Body.Close()\nreturn body\n}\n"

func main() {

	inputFile := os.Args[1] + ".go"
	outputFile := "./outDir/myMain.go"
	//reading code from source
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic("Error Reading file")
	}
	//splitting it into a slice of string to ease processing
	proc := strings.Split(string(data), "\n")
	//index to add imports
	var importIdx int = 1
	// Incval to increase the values of the stack according to what have been added

	for i := 0; i < len(proc); i++ {
		if importIdx != -1 {
			if strings.Contains(proc[i], "import") == true {
				if strings.Contains(proc[i], "(") == true {
					importIdx = -1
				} else {
					mapImports[strings.Split(proc[i], "\"")[1]] = true
					//fmt.Println("------")
					//fmt.Println(strings.Split(proc[i], "\"")[1])
					//fmt.Println("------")
					importIdx = i
				}
			}
		} else {
			if strings.Contains(proc[i], ")") == true {
				importIdx = i
			} else {
				//fmt.Println("------")
				//fmt.Println(strings.Split(proc[i], "\"")[1])
				//fmt.Println("------")
				mapImports[strings.Split(proc[i], "\"")[1]] = true
			}
		}

	}

	proc = addImports(proc, importIdx)

	var mainIdx int
	loopLen := len(proc)
	for i := 0; i < loopLen; i++ {
		if strings.Contains(proc[i], "func main()") == true {
			mainIdx = i
		}

		if strings.Contains(proc[i], "//VDCS") == true {
			//fmt.Println("I'm here and it's true")
			circ, params := extractCircuit(proc[i+1])
			typesA := getTypes(proc, params)
			//fmt.Println(typesA)
			for _, val := range typesA {
				circ += "_" + val
			}
			var chName string
			proc, chName = addComm(proc, circ, mainIdx)

			i += 2
			loopLen += 2

			proc = addEval(proc, i+1, params, typesA, chName)
			goCnt++
		}
	}

	proc = addServerFuncs(proc)

	/*for _, val := range proc {
		fmt.Println(string(val))
	}*/
	var myData []byte = []byte(strings.Join(proc, "\n"))
	err = ioutil.WriteFile(outputFile, myData, 0777)
	// handle this error
	if err != nil {
		// print it out
		fmt.Println(err)
	}

	/*
		cmd := exec.Command("go", "run", outputFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if runtime.GOOS == "windows" {
			cmd = exec.Command("tasklist")
		}
		err = cmd.Run()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
	*/
}

func addImports(s []string, idx int) []string {
	var concat string
	for key, val := range mapImports {
		if val == false {
			concat += "\"" + key + "\"\n"
		}
	}
	concat = "import (\n" + concat + ")\n" + circuitBlock

	s = append(s[:idx+1], append(strings.Split(concat, "\n"), s[idx+1:]...)...)
	return s
}

func addComm(s []string, circ string, mainIdx int) ([]string, string) {
	var chName string = "_" + circ + "Ch" + strconv.Itoa(goCnt)
	var call string = chName + ":= make(chan circuit)\ngo comm" + strconv.Itoa(goCnt) + "(\"" + circ + "\"," + strconv.Itoa(goCnt) + "," + chName + ")"
	//fmt.Println(call)
	s = append(s[:mainIdx+1], append(strings.Split(call, "\n"), s[mainIdx+1:]...)...)

	stpIdx := strings.Index(commBlock, "comm")
	sigComm := commBlock[:stpIdx+4] + strconv.Itoa(goCnt) + commBlock[stpIdx+4:]
	s = append(s, strings.Split(sigComm, "\n")...)
	return s, chName
}

func addEval(code []string, idx int, params, typesA []string, chName string) []string {
	code[idx] = strings.ReplaceAll(code[idx], "myStringMatch", "eval"+strconv.Itoa(goCnt))
	code[idx] = strings.Replace(code[idx], ")", ", "+strconv.Itoa(goCnt)+","+chName+")", 1)
	stpIdx := strings.Index(evalBlock, "eval")
	sigEval := evalBlock[:stpIdx+4] + strconv.Itoa(goCnt) + "("
	var inWires string = "{"
	for k, val := range params {
		sigEval += val + " " + strings.Split(typesA[k], "_")[0] + ","
		inWires += "\n_inWire" + strconv.Itoa(k) + ":=[]byte(" + val + ")\n"
	}

	sigEval += evalBlock[stpIdx+4:]
	sigEval = strings.Replace(sigEval, "{", inWires, 1)
	code = append(code, strings.Split(sigEval, "\n")...)
	return code
}

func extractCircuit(call string) (circ string, params []string) {

Loop:
	for _, i := range supportedFunc {
		if strings.Contains(call, i) == true {
			circ = i
			var tmp string = strings.Split(call, i)[1]
			tmp = strings.Split(tmp, "(")[1]
			params = append(params, strings.ReplaceAll(strings.Split(tmp, ",")[0], " ", ""))
			params = append(params, strings.ReplaceAll(strings.Split(strings.Split(tmp, ",")[1], ")")[0], " ", ""))
			break Loop
		}
	}
	return
}

func getTypes(code, params []string) (typesA []string) {

	n := "_8"
	k := "_8"

	inc := 0

	for _, val := range params {
		for _, line := range code {
			if strings.Contains(line, val) == true {
				if strings.Contains(line, "var") == true {
					segLine := strings.Split(strings.Split(line, "var")[1], " ")
					//fmt.Println(line, val)
					//fmt.Println(segLine)
					//fmt.Println(segLine[1], val)

					if segLine[1] == val {
						typesA = append(typesA, segLine[2])
						if inc == 1 {
							typesA[inc] += k
						} else {
							typesA[inc] += n
						}
						//fmt.Println(typesA)
						inc++
						break
					}
				} else if strings.Contains(line, "const") == true {
					segLine := strings.Split(strings.Split(line, "const")[1], " ")
					if segLine[1] == val {
						typesA = append(typesA, segLine[2])
						if inc == 1 {
							typesA[inc] += k
						} else {
							typesA[inc] += n
						}
						inc++
						break
					}
				} else {
					continue
				}
			}
		}
	}
	return
}

func addServerFuncs(code []string) []string {
	code = append(code, append(strings.Split(sendToGarbleBlock, "\n"), append(strings.Split(sendToEvalBlock, "\n"), append(strings.Split(getFromGarbleBlock, "\n"), strings.Split(getFromEvalBlock, "\n")...)...)...)...)
	return code
}
