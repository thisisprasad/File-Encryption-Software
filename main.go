package main

import "fmt"

var (
	cipher    DES_8encryption
	encryptor DESFileEncryptor
	input     []byte = []byte{0, 1, 1, 1, 1, 0, 0, 1}
)

func main() {
	cipher.Init("des_input.txt")

	encryptor.EncryptFile("test.txt")
	encryptor.DecryptFile("test.txt.enc")

	var encryptedData []byte = cipher.Encrypt(input)
	var decryptedData []byte = cipher.Decrypt(encryptedData)
	fmt.Println("Decrypted data:", decryptedData)
}
