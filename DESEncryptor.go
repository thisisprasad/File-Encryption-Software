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

type DESFileEncryptor struct {
	filename           string
	encryptionFilename string
	cipher             DES_8encryption
}

func (encryptor *DESFileEncryptor) getBinaryByteArray(byteVal byte) []byte {
	var byteArray []byte
	for byteVal > 0 {
		byteArray = append([]byte{(byte)(byteVal % 2)}, byteArray...)
		byteVal /= 2
	}
	for (int)(len(byteArray)) < 8 {
		byteArray = append([]byte{0}, byteArray...)
	}

	return byteArray
}

func (encryptor *DESFileEncryptor) convertBinaryByteArrayToByte(byteArray *[]byte) byte {
	var res byte = 0
	for i := len(*byteArray) - 1; i >= 0; i-- {
		if (*byteArray)[i] == (byte)(1) {
			res |= (1 << i)
		} else {
			mask := ^(1 << i)
			res &= (byte)(mask)
		}
	}

	return res
}

func (encryptor *DESFileEncryptor) encryptChunk(buffer *[]byte, encryptionBuffer *[][]byte) {
	for i := 0; i < len(*buffer); i++ {
		byteVal := (*buffer)[i]
		var binaryByteArray []byte = encryptor.getBinaryByteArray(byteVal)
		(*encryptionBuffer)[i] = encryptor.cipher.Encrypt(binaryByteArray)
	}
}

func (encryptor *DESFileEncryptor) writeEncryptionBufferToFile(encryptionBuffer *[][]byte, filename string) {
	// var byteVal byte
	fmt.Println("enccryption buffer length:", len(*encryptionBuffer))
	var byteArray []byte = make([]byte, len(*encryptionBuffer))
	var cache []byte = make([]byte, 8)
	for i := 0; i < len(*encryptionBuffer); i++ {
		copy(cache[:], (*encryptionBuffer)[i][:])
		byteVal := encryptor.convertBinaryByteArrayToByte(&cache)
		byteArray[i] = byteVal
	}

	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		_, ferr := os.Create(filename)
		if ferr != nil {
			log.Fatalln("Problem encrypting file", ferr)
		}
	}
	//	write byte array to buffer
	permissions := 0644
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, (os.FileMode)(permissions))
	if err != nil {
		log.Fatalln("Problem encrypting file", err)
	}
	defer file.Close()
	_, err := file.WriteLine
}

func (encryptor *DESFileEncryptor) EncryptFile(filename string) {
	encryptor.filename = filename
	encryptor.encryptionFilename = filename + ".enc"
	encryptor.cipher.Init("des_input.txt")

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
		encryptor.writeEncryptionBufferToFile(&encryptionBuffer, encryptor.encryptionFilename)

		fmt.Println("bytesread:", bytesread, "bytes to string:", string(buffer[:bytesread]))
	}
}
