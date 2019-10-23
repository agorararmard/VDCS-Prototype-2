package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var goCnt int

var supportedFunc [1]string = [1]string{"myStringMatch"}

var mapImports map[string]bool = map[string]bool{"fmt": false, "strings": false}

const commBlock string = "func comm(cir string,cID int) {\nfmt.Println(cir)\nfmt.Println(cID)\n//get the circuit in JSON format\n//Generate input wires\n//post to server\n//Wait for response\n}"
const evalBlock string = "func eval(cID int) {\nfmt.Println(cir)\n//generate input wires for given inputs\n//fetch the garbled circuit with the cID\n//post to server\n//Wait for response\n return true\n}"

func addImports(s []string, idx int) []string {
	var concat string
	for key, val := range mapImports {
		if val == false {
			concat += "\"" + key + "\"\n"
		}
	}
	concat = "import (\n" + concat + ")\n"

	s = append(s[:idx+1], append(strings.Split(concat, "\n"), s[idx+1:]...)...)
	return s
}

func main() {

	//reading code from source
	data, err := ioutil.ReadFile("main.goc")
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
					fmt.Println("------")
					fmt.Println(strings.Split(proc[i], "\"")[1])
					fmt.Println("------")
					importIdx = i
				}
			}
		} else {
			if strings.Contains(proc[i], ")") == true {
				importIdx = i
			} else {
				fmt.Println("------")
				fmt.Println(strings.Split(proc[i], "\"")[1])
				fmt.Println("------")
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
			fmt.Println("I'm here and it's true")
			proc = addComm(proc, i+1, mainIdx)
			proc = addEval(proc, i+1)
			goCnt++
			i++
			loopLen++
		}
	}

	for _, val := range proc {
		fmt.Println(string(val))
	}
}

func addComm(s []string, idx, mainIdx int) []string {
	circ, params := extractCircuit(s[idx])
	typesA := getTypes(s, params)
	fmt.Println(typesA)
	for _, val := range typesA {
		circ += "_" + val
	}
	fmt.Println("Let's create call")
	var call string = "go comm(" + circ + strconv.Itoa(goCnt) + ")"
	fmt.Println(call)
	s = append(s[:mainIdx+1], append(strings.Split(call, "\n"), s[mainIdx+1:]...)...)
	s = append(s, strings.Split(commBlock, "\n")...)
	fmt.Println("\n\n\n\nGoing out of add comm")
	return s
}

func addEval(s []string, idx int) []string {
	s = append(s, strings.Split(evalBlock, "\n")...)
	return s
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
	fmt.Println("Returning")
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
					fmt.Println(line, val)
					fmt.Println(segLine)
					fmt.Println(segLine[1], val)

					if segLine[1] == val {
						typesA = append(typesA, segLine[2])
						if inc == 1 {
							typesA[inc] += k
						} else {
							typesA[inc] += n
						}
						fmt.Println(typesA)
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
