package vdcs

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	cryptoRand "crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

type Wire struct {
	WireID    string `json:"WireID"`
	WireLabel []byte `json:"WireLabel"`
}
type Gate struct {
	GateID     string   `json:"GateID"`
	GateInputs []string `json:"GateInputs"`
}
type CircuitGate struct {
	Gate
	TruthTable []bool `json:"TruthTable"`
}
type GarbledGate struct {
	Gate
	GarbledValues [][]byte `json:"GarbledValues"`
}

type ComID struct {
	CID string `json:"ComID"`
}
type Circuit struct {
	InputGates  []CircuitGate `json:"InputGates"`
	MiddleGates []CircuitGate `json:"MiddleGates"`
	OutputGates []CircuitGate `json:"OutputGates"`
}
type Randomness struct {
	Rin       int64 `json:"Rin"`
	Rout      int64 `json:"Rout"`
	Rgc       int64 `json:"Rgc"`
	LblLength int   `json:"LblLength"`
}
type CircuitMessage struct {
	Circuit
	ComID
	Randomness
}
type GarbledCircuit struct {
	InputGates  []GarbledGate `json:"InputGates"`
	MiddleGates []GarbledGate `json:"MiddleGates"`
	OutputGates []GarbledGate `json:"OutputGates"`
	ComID
}

type GarbledMessage struct {
	InputWires []Wire `json:"InputWires"`
	GarbledCircuit
	OutputWires []Wire `json:"OutputWires"`
}

type ResEval struct {
	Res [][]byte `json:"Result"`
	ComID
}

//basically, the channel will need to send the input/output mapping as well
func Comm(cir string, cID int, chVDCSCommCircRes chan<- GarbledMessage) {
	file, _ := ioutil.ReadFile(cir + ".json")
	k := Circuit{}
	err := json.Unmarshal([]byte(file), &k) //POSSIBLE BUG
	if err != nil {
		log.Fatal(err)
	}
	rand.Seed(int64(cID))
	mCirc := CircuitMessage{Circuit: Circuit{
		InputGates:  k.InputGates,
		MiddleGates: k.MiddleGates,
		OutputGates: k.OutputGates,
	},
		ComID: ComID{strconv.Itoa(rand.Int())},
		Randomness: Randomness{Rin: rand.Int63(),
			Rout:      rand.Int63(),
			Rgc:       rand.Int63(),
			LblLength: 16, //Should be rand.Int()%16 + 16
		},
	}
	//fmt.Println(mCirc)

	for !SendToServerGarble(mCirc) {

	}

	//Generate input wires
	//Wait for response
	inputSize := len(mCirc.InputGates) * 2
	outputSize := len(mCirc.OutputGates)

	arrIn := YaoGarbledCkt_in(mCirc.Rin, mCirc.LblLength, inputSize)
	arrOut := YaoGarbledCkt_out(mCirc.Rout, mCirc.LblLength, outputSize)
	var gcm GarbledMessage
	var oke bool
	for gcm, oke = GetFromServerGarble(mCirc.CID); !oke; {

	}

	//gcm = Garble(mCirc)
	//Validate Correctness of result
	//fmt.Println(gcm)
	//fmt.Println("\nHere:\n", arrIn, "\nThere\n", arrOut)

	for k, val := range gcm.InputWires {
		if bytes.Compare(arrIn[k], val.WireLabel) != 0 {
			fmt.Println("I was cheated on this: ", arrIn[k], val.WireLabel)
			panic("The server has cheated me") //redo the process, by recovering from panic by recalling comm
		}
	}
	for k, val := range gcm.OutputWires {
		if bytes.Compare(arrOut[k], val.WireLabel) != 0 {

			fmt.Println("I was cheated on this: ", arrOut[k], val.WireLabel)
			panic("The server has cheated me") //redo the process, by recovering from panic by recalling comm
		}
	}
	//Send Circuit to channel
	chVDCSCommCircRes <- gcm
}

func SendToServerGarble(k CircuitMessage) bool {
	circuitJSON, err := json.Marshal(k)
	req, err := http.NewRequest("POST", "http://localhost:8080/post", bytes.NewBuffer(circuitJSON))
	if err != nil {
		fmt.Println("generating request failed")
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	resp.Body.Close()
	if err != nil {
		//log.Fatal(err)
		return false
	}
	return true
}

func GetFromServerGarble(id string) (k GarbledMessage, ok bool) {
	ok = false //assume failure
	iDJSON, err := json.Marshal(ComID{CID: id})
	req, err := http.NewRequest("GET", "http://localhost:8080/get", bytes.NewBuffer(iDJSON))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &k)
	if err != nil {
		return
	}
	resp.Body.Close()
	if k.CID != id {
		panic("The server sent me the wrong circuit") //replace with a request repeat.
	}
	ok = true
	return
}

func SendToServerEval(k GarbledMessage) bool {
	circuitJSON, err := json.Marshal(k)
	req, err := http.NewRequest("POST", "http://localhost:8081/post", bytes.NewBuffer(circuitJSON))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return true
}

func GetFromServerEval(id string) (res [][]byte, ok bool) {
	ok = false // assume failure
	iDJSON, err := json.Marshal(ComID{CID: id})
	req, err := http.NewRequest("GET", "http://localhost:8081/get", bytes.NewBuffer(iDJSON))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	var k ResEval
	err = json.Unmarshal(body, &k)
	if err != nil {
		return
	}
	resp.Body.Close()
	if k.CID != id {
		panic("The server sent me the wrong circuit") //replace with a request repeat.
	}
	res = k.Res
	fmt.Println("Result Returned", k.Res)
	ok = true
	return
}
func GenNRandNumbers(n int, length int, r int64, considerR bool) [][]byte {
	if considerR {
		rand.Seed(r)
	}
	seeds := make([][]byte, n)
	for i := 0; i < n; i++ {
		seeds[i] = make([]byte, length)
		_, err := rand.Read(seeds[i])
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
	return seeds
}

func YaoGarbledCkt_in(rIn int64, length int, inputSize int) [][]byte {
	return GenNRandNumbers(2*inputSize, length, rIn, true)
}

func YaoGarbledCkt_out(rOut int64, length int, outputSize int) [][]byte {
	// only one output bit for now
	return GenNRandNumbers(2*outputSize, length, rOut, true)
}

func EncryptAES(encKey []byte, plainText []byte) (ciphertext []byte, ok bool) {

	ok = false //assume failure
	//			encKey = append(encKey, hash)
	c, err := aes.NewCipher(encKey)
	if err != nil {
		//fmt.Println(err)
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		//fmt.Println(err)
		return
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(cryptoRand.Reader, nonce); err != nil {
		//fmt.Println(err)
		return
	}
	ciphertext = gcm.Seal(nonce, nonce, plainText, nil)
	//fmt.Println(ciphertext)
	ok = true

	return
}

func DecryptAES(encKey []byte, cipherText []byte) (plainText []byte, ok bool) {

	ok = false //assume failure

	c, err := aes.NewCipher(encKey)
	if err != nil {
		//fmt.Println(err)
		return
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		//fmt.Println(err)
		return
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		//fmt.Println(err)
		return
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plainText, err = gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		//fmt.Println(err)
		return
	}
	//fmt.Println(string(plaintext))
	ok = true
	return
}

func Garble(circ CircuitMessage) GarbledMessage {

	inputSize := len(circ.InputGates) * 2
	outputSize := len(circ.OutputGates)
	arrIn := YaoGarbledCkt_in(circ.Rin, circ.LblLength, inputSize)
	arrOut := YaoGarbledCkt_out(circ.Rout, circ.LblLength, outputSize)

	inWires := make(map[string][]Wire)
	outWires := make(map[string][]Wire)

	rand.Seed(circ.Rgc)

	var gc GarbledCircuit
	inputWiresGC := []Wire{}
	outputWiresGC := []Wire{}

	gc.CID = circ.CID

	// Input Gates Garbling
	var wInCnt int = 0
	for k, val := range circ.InputGates {
		gc.InputGates = append(gc.InputGates, GarbledGate{
			Gate: Gate{
				GateID: val.GateID,
			},
		})

		gc.InputGates[k].GateInputs = val.GateInputs

		inCnt := int(math.Log2(float64(len(val.TruthTable))))

		//fmt.Printf("%v, %T\n", val.GateID, val.GateID)

		inWires[val.GateID] = []Wire{}

		for i := 0; i < inCnt; i++ {
			inWires[val.GateID] = append(inWires[val.GateID], Wire{
				WireLabel: arrIn[wInCnt],
			}, Wire{
				WireLabel: arrIn[wInCnt+1],
			})
			inputWiresGC = append(inputWiresGC, Wire{
				WireLabel: arrIn[wInCnt],
			}, Wire{
				WireLabel: arrIn[wInCnt+1],
			})
			wInCnt += 2
		}
		outWires[val.GateID] = []Wire{}
		outWire := GenNRandNumbers(2, circ.LblLength, 0, false)
		outWires[val.GateID] = append(outWires[val.GateID], Wire{
			WireLabel: outWire[0],
		}, Wire{
			WireLabel: outWire[1],
		})
		//in1:	0	0	1	1
		//in0:	0	1	0	1
		//out:	1	0	0	1

		//fmt.Println("Here we getting inWires: \n")
		gc.InputGates[k].GarbledValues = make([][]byte, len(val.TruthTable))
		for key, value := range val.TruthTable {
			var concat []byte
			for i := 0; i < inCnt; i++ {
				idx := ((key >> i) & (1))
				concat = append(concat, inWires[val.GateID][(i*2)+idx].WireLabel...)
			}
			concat = append(concat, []byte(val.GateID)...)
			hash := sha256.Sum256(concat)

			var idxOut int
			if value {
				idxOut = 1
			}
			outKey := outWires[val.GateID][int(idxOut)]
			// generate a new aes cipher using our 32 byte long key
			encKey := make([]byte, 32)
			for jk, tmpo := range hash {
				encKey[jk] = tmpo
			}
			var ok bool
			gc.InputGates[k].GarbledValues[key], ok = EncryptAES(encKey, outKey.WireLabel)
			if !ok {
				fmt.Println("Encryption Failed")
			}
		}
		//fmt.Println("\nwe got'em inWires \n")

	}

	//Middle Gates Garbling
	for k, val := range circ.MiddleGates {
		gc.MiddleGates = append(gc.MiddleGates, GarbledGate{
			Gate: Gate{
				GateID: val.GateID,
			},
		})

		gc.MiddleGates[k].GateInputs = val.GateInputs

		inCnt := int(math.Log2(float64(len(val.TruthTable))))

		//fmt.Printf("%v, %T\n", val.GateID, val.GateID)
		inWires[val.GateID] = []Wire{}

		for _, j := range val.GateInputs {
			inWires[val.GateID] = append(inWires[val.GateID], outWires[j][0])
			inWires[val.GateID] = append(inWires[val.GateID], outWires[j][1])
			//wInCnt++
		}
		outWires[val.GateID] = []Wire{}
		outWire := GenNRandNumbers(2, circ.LblLength, 0, false)
		outWires[val.GateID] = append(outWires[val.GateID], Wire{
			WireLabel: outWire[0],
		}, Wire{
			WireLabel: outWire[1],
		})

		gc.MiddleGates[k].GarbledValues = make([][]byte, len(val.TruthTable))
		for key, value := range val.TruthTable {
			//Concatinating the wire labels with the GateID
			var concat []byte
			for i := 0; i < inCnt; i++ {
				idx := ((key >> i) & (1))
				concat = append(concat, inWires[val.GateID][(i*2)+idx].WireLabel...)
			}
			concat = append(concat, []byte(val.GateID)...)

			//Hashing the value
			hash := sha256.Sum256(concat)

			//Determining the value of the output wire
			var idxOut int
			if value {
				idxOut = 1
			}
			outKey := outWires[val.GateID][int(idxOut)]

			// generate a new aes cipher using our 32 byte long key
			encKey := make([]byte, 32)
			for jk, tmpo := range hash {
				encKey[jk] = tmpo
			}
			var ok bool
			gc.MiddleGates[k].GarbledValues[key], ok = EncryptAES(encKey, outKey.WireLabel)
			if !ok {
				fmt.Println("Encryption Failed")
			}
		}

	}

	//Output Gates Garbling
	wOutCnt := 0
	for k, val := range circ.OutputGates {
		gc.OutputGates = append(gc.OutputGates, GarbledGate{
			Gate: Gate{
				GateID: val.GateID,
			},
		})

		gc.OutputGates[k].GateInputs = val.GateInputs

		inCnt := int(math.Log2(float64(len(val.TruthTable))))

		//fmt.Printf("%v, %T\n", val.GateID, val.GateID)

		inWires[val.GateID] = []Wire{}
		for _, j := range val.GateInputs {
			inWires[val.GateID] = append(inWires[val.GateID], outWires[j][0])
			inWires[val.GateID] = append(inWires[val.GateID], outWires[j][1])

			//wInCnt++
		}

		outWires[val.GateID] = []Wire{}

		outWires[val.GateID] = append(outWires[val.GateID], Wire{
			WireLabel: arrOut[wOutCnt],
		}, Wire{
			WireLabel: arrOut[wOutCnt+1],
		})

		outputWiresGC = append(outputWiresGC, Wire{
			WireLabel: arrOut[wOutCnt],
		}, Wire{
			WireLabel: arrOut[wOutCnt+1],
		})
		wOutCnt += 2

		gc.OutputGates[k].GarbledValues = make([][]byte, len(val.TruthTable))
		for key, value := range val.TruthTable {
			var concat []byte
			for i := 0; i < inCnt; i++ {
				idx := ((key >> i) & (1))
				concat = append(concat, inWires[val.GateID][(i*2)+idx].WireLabel...)
			}
			concat = append(concat, []byte(val.GateID)...)
			hash := sha256.Sum256(concat)

			var idxOut int
			if value {
				idxOut = 1
			}
			outKey := outWires[val.GateID][int(idxOut)]
			// generate a new aes cipher using our 32 byte long key
			encKey := make([]byte, 32)
			for jk, tmpo := range hash {
				encKey[jk] = tmpo
			}
			var ok bool
			gc.OutputGates[k].GarbledValues[key], ok = EncryptAES(encKey, outKey.WireLabel)
			if !ok {
				fmt.Println("Encryption Failed")
			}
		}

	}

	//fmt.Println(arrIn)
	//fmt.Println(arrOut)
	//fmt.Println("Input Wires GC:", inWires)
	//fmt.Println("Output Wires GC:", outWires)
	//fmt.Println("GC: ", gc)
	gm := GarbledMessage{
		InputWires:     inputWiresGC,
		GarbledCircuit: gc,
		OutputWires:    outputWiresGC,
	}
	return gm
}

func Evaluate(gc GarbledMessage) (result ResEval) {

	result.CID = gc.CID
	outWires := make(map[string]Wire)
	var wInCnt int

	for _, val := range gc.InputGates {

		inCnt := int(math.Log2(float64(len(val.GarbledValues))))
		var concat []byte
		for i := 0; i < inCnt; i++ {
			concat = append(concat, gc.InputWires[wInCnt].WireLabel...)
			wInCnt++
		}
		concat = append(concat, []byte(val.GateID)...)
		hash := sha256.Sum256(concat)
		encKey := make([]byte, 32)
		for jk, tmpo := range hash {
			encKey[jk] = tmpo
		}
		outWires[val.GateID] = Wire{}
		for _, gValue := range val.GarbledValues {
			tmpWireLabel, ok := DecryptAES(encKey, gValue)
			if ok {
				outWires[val.GateID] = Wire{
					WireLabel: tmpWireLabel,
				}
				break
			}
		}

		if (bytes.Compare(outWires[val.GateID].WireLabel, Wire{}.WireLabel)) == 0 {
			fmt.Println("Fail Evaluation Input Gate")
		} /*else {
			fmt.Println("\n\nYaaay\nGate ", val.GateID, " Now has an output wire: \n", outWires[val.GateID].WireLabel, "\n\n")
		}*/
	}
	for _, val := range gc.MiddleGates {

		//inCnt := len(val.GateInputs)
		var concat []byte
		for _, preGate := range val.GateInputs {
			concat = append(concat, outWires[preGate].WireLabel...)
			//wInCnt++
		}
		concat = append(concat, []byte(val.GateID)...)
		hash := sha256.Sum256(concat)
		encKey := make([]byte, 32)
		for jk, tmpo := range hash {
			encKey[jk] = tmpo
		}
		outWires[val.GateID] = Wire{}
		for _, gValue := range val.GarbledValues {
			tmpWireLabel, ok := DecryptAES(encKey, gValue)
			if ok {
				outWires[val.GateID] = Wire{
					WireLabel: tmpWireLabel,
				}
				break
			}
		}
		if (bytes.Compare(outWires[val.GateID].WireLabel, Wire{}.WireLabel)) == 0 {
			fmt.Println("Fail Evaluation Middle Gate")
		} /*else {
			fmt.Println("\n\nYaaay\nGate ", val.GateID, " Now has an output wire: \n", outWires[val.GateID].WireLabel, "\n\n")
		}*/
	}

	for _, val := range gc.OutputGates {

		//inCnt := len(val.GateInputs)
		var concat []byte
		for _, preGate := range val.GateInputs {
			concat = append(concat, outWires[preGate].WireLabel...)
			//wInCnt++
		}
		concat = append(concat, []byte(val.GateID)...)
		hash := sha256.Sum256(concat)
		encKey := make([]byte, 32)
		for jk, tmpo := range hash {
			encKey[jk] = tmpo
		}
		outWires[val.GateID] = Wire{}
		for _, gValue := range val.GarbledValues {
			tmpWireLabel, ok := DecryptAES(encKey, gValue)
			if ok {
				//fmt.Println("\nI found my way out\n")
				outWires[val.GateID] = Wire{
					WireLabel: tmpWireLabel,
				}
				result.Res = append(result.Res, tmpWireLabel)
				break
			} /*else {
				fmt.Println("\nStill Trying to Find my way out\n")
			}*/
		}
		if (bytes.Compare(outWires[val.GateID].WireLabel, Wire{}.WireLabel)) == 0 {
			fmt.Println("Fail Evaluation Output Gate")
		} /*else {
			fmt.Println("\n\nYaaay\nGate ", val.GateID, " Now has an output wire: \n", outWires[val.GateID].WireLabel, "\n\n")
		}*/
	}

	return
}
