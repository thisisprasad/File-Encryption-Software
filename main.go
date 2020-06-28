package main

var (
	cipher DES_8encryption
	input  string = "10100101"
)

func main() {
	cipher.Init("des_input.txt")

	cipher.Encrypt(input)
}
