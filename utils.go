package main

import (
	"log"
	"strconv"
	"strings"
)

func StringToArray(s string) []int {
	var res []int

	cache := strings.Split(s, " ")
	for i := 0; i < len(cache); i++ {
		num, err := strconv.Atoi(cache[i])
		if err != nil {
			log.Fatalln("Problem converting string to Int")
		}
		res = append(res, num)
	}

	return res
}
