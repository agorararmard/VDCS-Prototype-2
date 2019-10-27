package main

import (
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
	"os"
	"strconv"

	//"crypto"
	"math/rand"

	"./vdcs"
)

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

func Garble(circ vdcs.CircuitMessage) vdcs.GarbledMessage {

	inputSize := len(circ.InputGates) * 2
	outputSize := len(circ.OutputGates)
	arrIn := YaoGarbledCkt_in(circ.Rin, circ.LblLength, inputSize)
	arrOut := YaoGarbledCkt_out(circ.Rout, circ.LblLength, outputSize)

	inWires := make(map[string][]vdcs.Wire)
	outWires := make(map[string][]vdcs.Wire)

	rand.Seed(circ.Rgc)

	var gc vdcs.GarbledCircuit
	inputWiresGC := []vdcs.Wire{}
	outputWiresGC := []vdcs.Wire{}

	gc.CID = circ.CID

	// Input Gates Garbling
	var wInCnt int = 0
	for k, val := range circ.InputGates {
		gc.InputGates = append(gc.InputGates, vdcs.GarbledGate{
			Gate: vdcs.Gate{
				GateID: val.GateID,
			},
		})

		gc.InputGates[k].GateInputs = val.GateInputs

		inCnt := int(math.Log2(float64(len(val.TruthTable))))

		//fmt.Printf("%v, %T\n", val.GateID, val.GateID)

		inWires[val.GateID] = []vdcs.Wire{}

		for i := 0; i < inCnt; i++ {
			inWires[val.GateID] = append(inWires[val.GateID], vdcs.Wire{
				WireLabel: arrIn[wInCnt],
			}, vdcs.Wire{
				WireLabel: arrIn[wInCnt+1],
			})
			inputWiresGC = append(inputWiresGC, vdcs.Wire{
				WireLabel: arrIn[wInCnt],
			}, vdcs.Wire{
				WireLabel: arrIn[wInCnt+1],
			})
			wInCnt += 2
		}
		outWires[val.GateID] = []vdcs.Wire{}
		outWire := GenNRandNumbers(2, circ.LblLength, 0, false)
		outWires[val.GateID] = append(outWires[val.GateID], vdcs.Wire{
			WireLabel: outWire[0],
		}, vdcs.Wire{
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
			//			encKey = append(encKey, hash)
			c, err := aes.NewCipher(encKey)
			if err != nil {
				fmt.Println(err)
			}
			gcm, err := cipher.NewGCM(c)
			if err != nil {
				fmt.Println(err)
			}
			nonce := make([]byte, gcm.NonceSize())
			if _, err = io.ReadFull(cryptoRand.Reader, nonce); err != nil {
				fmt.Println(err)
			}
			ciphertext := gcm.Seal(nonce, nonce, outKey.WireLabel, nil)
			//fmt.Println(ciphertext)

			gc.InputGates[k].GarbledValues[key] = ciphertext
		}
		//fmt.Println("\nwe got'em inWires \n")

	}

	//Middle Gates Garbling
	for k, val := range circ.MiddleGates {
		gc.MiddleGates = append(gc.MiddleGates, vdcs.GarbledGate{
			Gate: vdcs.Gate{
				GateID: val.GateID,
			},
		})

		gc.MiddleGates[k].GateInputs = val.GateInputs

		inCnt := int(math.Log2(float64(len(val.TruthTable))))

		//fmt.Printf("%v, %T\n", val.GateID, val.GateID)
		inWires[val.GateID] = []vdcs.Wire{}

		for _, j := range val.GateInputs {
			inWires[val.GateID] = append(inWires[val.GateID], outWires[j][0])
			inWires[val.GateID] = append(inWires[val.GateID], outWires[j][1])
			//wInCnt++
		}
		outWires[val.GateID] = []vdcs.Wire{}
		outWire := GenNRandNumbers(2, circ.LblLength, 0, false)
		outWires[val.GateID] = append(outWires[val.GateID], vdcs.Wire{
			WireLabel: outWire[0],
		}, vdcs.Wire{
			WireLabel: outWire[1],
		})

		gc.MiddleGates[k].GarbledValues = make([][]byte, len(val.TruthTable))
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
			//			encKey = append(encKey, hash)
			c, err := aes.NewCipher(encKey)
			if err != nil {
				fmt.Println(err)
			}
			gcm, err := cipher.NewGCM(c)
			if err != nil {
				fmt.Println(err)
			}
			nonce := make([]byte, gcm.NonceSize())
			if _, err = io.ReadFull(cryptoRand.Reader, nonce); err != nil {
				fmt.Println(err)
			}
			ciphertext := gcm.Seal(nonce, nonce, outKey.WireLabel, nil)
			//fmt.Println(ciphertext)

			gc.MiddleGates[k].GarbledValues[key] = ciphertext
		}

	}

	//Output Gates Garbling
	wOutCnt := 0
	for k, val := range circ.OutputGates {
		gc.OutputGates = append(gc.OutputGates, vdcs.GarbledGate{
			Gate: vdcs.Gate{
				GateID: val.GateID,
			},
		})

		gc.OutputGates[k].GateInputs = val.GateInputs

		inCnt := int(math.Log2(float64(len(val.TruthTable))))

		//fmt.Printf("%v, %T\n", val.GateID, val.GateID)

		inWires[val.GateID] = []vdcs.Wire{}
		for _, j := range val.GateInputs {
			inWires[val.GateID] = append(inWires[val.GateID], outWires[j][0])
			inWires[val.GateID] = append(inWires[val.GateID], outWires[j][1])

			//wInCnt++
		}

		outWires[val.GateID] = []vdcs.Wire{}

		outWires[val.GateID] = append(outWires[val.GateID], vdcs.Wire{
			WireLabel: arrOut[wOutCnt],
		}, vdcs.Wire{
			WireLabel: arrOut[wOutCnt+1],
		})

		outputWiresGC = append(outputWiresGC, vdcs.Wire{
			WireLabel: arrIn[wOutCnt],
		}, vdcs.Wire{
			WireLabel: arrIn[wOutCnt+1],
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
			//			encKey = append(encKey, hash)
			c, err := aes.NewCipher(encKey)
			if err != nil {
				fmt.Println(err)
			}
			gcm, err := cipher.NewGCM(c)
			if err != nil {
				fmt.Println(err)
			}
			nonce := make([]byte, gcm.NonceSize())
			if _, err = io.ReadFull(cryptoRand.Reader, nonce); err != nil {
				fmt.Println(err)
			}
			ciphertext := gcm.Seal(nonce, nonce, outKey.WireLabel, nil)
			//fmt.Println(ciphertext)

			gc.OutputGates[k].GarbledValues[key] = ciphertext
		}

	}

	//fmt.Println(arrIn)
	//fmt.Println(arrOut)
	//fmt.Println("Input Wires GC:", inWires)
	//fmt.Println("Output Wires GC:", outWires)
	//fmt.Println("GC: ", gc)
	gm := vdcs.GarbledMessage{
		InputWires:     inputWiresGC,
		GarbledCircuit: gc,
		OutputWires:    outputWiresGC,
	}
	return gm
}

func main() {
	file, _ := ioutil.ReadFile("cir.json")
	k := vdcs.Circuit{}
	err := json.Unmarshal([]byte(file), &k) //POSSIBLE BUG
	//fmt.Println("Here is k Input Gates:", k.InputGates)
	//fmt.Println("Here is k Middle Gates:", k.MiddleGates)
	//fmt.Println("Here is k Output Gates:", k.OutputGates)
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(int64(1))

	mCirc := vdcs.CircuitMessage{
		Circuit: vdcs.Circuit{
			InputGates:  k.InputGates,
			MiddleGates: k.MiddleGates,
			OutputGates: k.OutputGates,
		},
		ComID: vdcs.ComID{CID: strconv.Itoa(rand.Int())},
		Randomness: vdcs.Randomness{
			Rin:       int64(rand.Int()),
			Rout:      int64(rand.Int()),
			Rgc:       int64(rand.Int()),
			LblLength: 16,
		},
	}

	gCirMes := Garble(mCirc)
	fmt.Println("println: ", gCirMes)
}
