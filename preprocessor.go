package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var goCnt int

var supportedFunc [1]string = [1]string{"myStringMatch"}

var mapImports map[string]bool = map[string]bool{"fmt": false, "sync": false}

const commBlock string = "func comm(cir string,cID int) {\nfmt.Println(cir)\nfmt.Println(cID)\n//get the circuit in JSON format\n//Generate input wires\n//post to server\n//Wait for response\nwg.Done()}"
const evalBlock string = "func evalcID int) (bool){\nwg.Wait()\ncir := \"You did it!\"\nfmt.Println(cir)\n//generate input wires for given inputs\n//fetch the garbled circuit with the cID\n//post to server\n//Wait for response\n return true\n}"

func main() {

	//reading code from source
	data, err := ioutil.ReadFile("inDir/main.go")
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
			fmt.Println(params, typesA)
			proc = addComm(proc, circ, mainIdx)
			proc = addEval(proc, i+1, params, typesA)
			goCnt++
			i++
			loopLen++
		}
	}
	proc = addWGADD(proc, mainIdx)

	for _, val := range proc {
		fmt.Println(string(val))
	}
	var myData []byte = []byte(strings.Join(proc, "\n"))
	err = ioutil.WriteFile("outDir/myMain.go", myData, 0777)
	// handle this error
	if err != nil {
		// print it out
		fmt.Println(err)
	}

}

func addImports(s []string, idx int) []string {
	var concat string
	for key, val := range mapImports {
		if val == false {
			concat += "\"" + key + "\"\n"
		}
	}
	concat = "import (\n" + concat + ")\nvar wg = sync.WaitGroup{}\n"

	s = append(s[:idx+1], append(strings.Split(concat, "\n"), s[idx+1:]...)...)
	return s
}

func addComm(s []string, circ string, mainIdx int) []string {

	var call string = "go comm" + strconv.Itoa(goCnt) + "(\"" + circ + "\"," + strconv.Itoa(goCnt) + ")"
	//fmt.Println(call)
	s = append(s[:mainIdx+1], append(strings.Split(call, "\n"), s[mainIdx+1:]...)...)

	stpIdx := strings.Index(commBlock, "comm")
	sigComm := commBlock[:stpIdx+4] + strconv.Itoa(goCnt) + commBlock[stpIdx+4:]
	s = append(s, strings.Split(sigComm, "\n")...)
	return s
}

func addEval(code []string, idx int, params, typesA []string) []string {
	fmt.Println(params, typesA)
	idx++
	code[idx] = strings.ReplaceAll(code[idx], "myStringMatch", "eval"+strconv.Itoa(goCnt))
	code[idx] = strings.Replace(code[idx], ")", ", "+strconv.Itoa(goCnt)+")", 1)
	stpIdx := strings.Index(evalBlock, "eval")
	sigEval := evalBlock[:stpIdx+4] + strconv.Itoa(goCnt) + "("
	for k, val := range params {
		sigEval += val + " " + strings.Split(typesA[k], "_")[0] + ","
	}
	sigEval += evalBlock[stpIdx+4:]
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
	k := "_4"

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

func addWGADD(code []string, mainIdx int) []string {
	call := "wg.Add(" + strconv.Itoa(goCnt) + ")\n"
	code = append(code[:mainIdx+1], append(strings.Split(call, "\n"), code[mainIdx+1:]...)...)
	return code
}
