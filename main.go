package main

import "fmt"

var cipher DES_8encryption

func main() {
	cipher.init("des_input.txt")
	fmt.Println("key:", cipher.key)
	fmt.Println("p10: ", cipher.p10)
}
