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
	p10                 []int
	initialPermutation  []int
	inversePermutation  []int
	expansionPermuation []int
	s0                  [][]int //	s0 matrix
	s1                  [][]int //	s1-matrix
	key                 string
	key1                string
	key2                string
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
func (cipher *DES_8encryption) init(configFile string) {
	log.Println("Initializing cipher")

	var s string
	file, err := os.Open(configFile)
	if err != nil {
		log.Fatalln("Error opening config file", configFile, "Error: ", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	s = scanner.Text()
	cipher.key = strings.Split(s, ":")[1]

	scanner.Scan()
	s = strings.Split(s, ":")[1]
	fmt.Println("p10 split: ", s)
	cipher.p10 = StringToArray(s)
}
