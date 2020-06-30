package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	bufferSize = 10
)

var (
	cipher DES_8encryption
)

type DESFileEncryptor struct {
	filename string
}

func (encryptor *DESFileEncryptor) getBinaryByteArray(byteVal byte) []byte {
	var byteArray []byte
	for byteVal > 0 {
		byteArray = append([]byte{(byte)(byteVal % 2)}, byteArray...)
		byteVal /= 2
	}

	return byteArray
}

func (encryptor *DESFileEncryptor) encryptChunk(buffer *[]byte, encryptionBuffer *[][]byte) {
	for i := 0; i < len(*buffer); i++ {
		byteVal := (*buffer)[i]
		var binaryByteArray []byte = encryptor.getBinaryByteArray(byteVal)
		(*encryptionBuffer)[i] = cipher.Encrypt(binaryByteArray)
	}
}

func (encryptor *DESFileEncryptor) writeEncryptionBufferToFile()

func (encryptor *DESFileEncryptor) EncryptFile(filename string) {
	encryptor.filename = filename

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln("Problem encrypting file", err)
		return
	}
	defer file.Close()

	var buffer []byte
	buffer = make([]byte, bufferSize)
	var encryptionBuffer [][]byte
	encryptionBuffer = make([][]byte, bufferSize)
	for i := 0; i < bufferSize; i++ {
		encryptionBuffer[i] = make([]byte, 8)
	}

	for {
		bytesread, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Fatalln("Problem while reading the contents of the file.")
			}
			break
		}

		encryptor.encryptChunk(&buffer, &encryptionBuffer)
		//	Write encryptionBuffer into file.
		encryptor.writeEncryptionBufferToFile(&encryptionBuffer)

		fmt.Println("bytesread:", bytesread, "bytes to string:", string(buffer[:bytesread]))
	}
}
