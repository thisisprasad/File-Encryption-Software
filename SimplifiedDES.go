package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type DES_8encryption struct {
	p4                  []int //	4-length permutation
	p8                  []int //	8-length permutation
	p10                 []int //	10-length permuation
	initialPermutation  []int
	inversePermutation  []int
	expansionPermuation []int
	s0                  [][]int //	s0 matrix
	s1                  [][]int //	s1-matrix
	key                 []byte
	key1                []byte
	key2                []byte
}

/**
Sequence of reading input is strict
e.g.:
key:0010010111
P10:3 5 2 7 4 10 1 9 8 6
P8:6 3 7 4 8 5 10 9
P4:2 4 3 1
IP:2 6 3 1 4 8 5 7
IP-1:4 1 3 5 7 2 8 6
E/P:4 1 2 3 2 3 4 1
S0:1 0 3 2 3 2 1 0 0 2 1 3 3 1 3 2
S1:0 1 2 3 2 0 1 3 3 0 1 0 2 1 0 3
*/
func (cipher *DES_8encryption) readFile(configFile string) {
	var s string
	file, err := os.Open(configFile)
	if err != nil {
		log.Fatalln("Error opening config file", configFile, "Error: ", err)
	}
	defer file.Close()

	//	Key
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	s = scanner.Text()
	cipher.key = []byte(strings.Split(s, ":")[1])
	//	Converting char value to integer value
	for i := 0; i < len(cipher.key); i++ {
		cipher.key[i] = cipher.key[i] - '0'
	}

	//	p10 permutation
	scanner.Scan()
	s = strings.Split(scanner.Text(), ":")[1]
	cipher.p10 = StringToIntArray(s)

	//	p8 permutation
	scanner.Scan()
	s = strings.Split(scanner.Text(), ":")[1]
	cipher.p8 = StringToIntArray(s)

	//	p4 permuation
	scanner.Scan()
	s = strings.Split(scanner.Text(), ":")[1]
	cipher.p4 = StringToIntArray(s)

	//	initial permutation
	scanner.Scan()
	s = strings.Split(scanner.Text(), ":")[1]
	cipher.initialPermutation = StringToIntArray(s)

	//	inverse permutation
	scanner.Scan()
	s = strings.Split(scanner.Text(), ":")[1]
	cipher.inversePermutation = StringToIntArray(s)

	//	expansion permuation
	scanner.Scan()
	s = strings.Split(scanner.Text(), ":")[1]
	cipher.expansionPermuation = StringToIntArray(s)

	//	s0 matrix
	scanner.Scan()
	s = strings.Split(scanner.Text(), ":")[1]
	cipher.s0 = StringTo2DIntArray(s, 4, 4)

	//	s1 matrix
	scanner.Scan()
	s = strings.Split(scanner.Text(), ":")[1]
	cipher.s1 = StringTo2DIntArray(s, 4, 4)
}

func (cipher *DES_8encryption) generateIntermediateKeys() {
	var p10key = cipher.applyPermutation(cipher.key, cipher.p10)
	var leftHalf []byte = p10key[0:5]
	var rightHalf []byte = p10key[5:]

	fmt.Println("lefthalf:", leftHalf, "righthalf:", rightHalf)
	cipher.circularLeftShift(&leftHalf, 1)
	cipher.circularLeftShift(&rightHalf, 1)
	fmt.Println("lefthalf:", leftHalf, "righthalf:", rightHalf)
}

func (cipher *DES_8encryption) Init(configFile string) {
	log.Println("Initializing cipher...")

	cipher.readFile(configFile)
	cipher.generateIntermediateKeys()

	log.Println("cipher intialization complete...")
}

func (cipher *DES_8encryption) circularLeftShift(data *[]byte, shiftBy int) {
	shiftBy %= len(*data)
	var cache []byte = make([]byte, shiftBy)
	for i := 0; i < shiftBy; i++ {
		cache[i] = (*data)[i]
	}
	for i := shiftBy; i < len(*data); i++ {
		(*data)[i-shiftBy] = (*data)[i]
	}
	pos := 0
	for i := len(*data) - shiftBy; i < len(*data); i++ {
		(*data)[i] = cache[pos]
		pos += 1
	}
}

func (cipher *DES_8encryption) applyPermutation(data []byte, permutation []int) []byte {
	var res []byte = make([]byte, len(data))
	pos := 0
	for i := 0; i < len(permutation); i++ {
		ch := data[permutation[i]-1] //	0-indexed
		res[pos] = ch
		pos += 1
	}

	return res
}

func (cipher *DES_8encryption) Encrypt(plainText string) string {
	// ipBits := cipher.applyPermutation(plainText, cipher.initialPermutation)

	return ""
}
