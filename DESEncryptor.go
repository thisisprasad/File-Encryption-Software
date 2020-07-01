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
	filename                string
	encryptionFilename      string
	cipher                  DES_8encryption
	decryptionFileConnector *os.File
	decryptionFilename      string
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

func (encryptor *DESFileEncryptor) encryptChunk(buffer *[]byte, bufferDataSize int, encryptionBuffer *[][]byte) {
	for i := 0; i < bufferDataSize; i++ {
		byteVal := (*buffer)[i]
		var binaryByteArray []byte = encryptor.getBinaryByteArray(byteVal)
		(*encryptionBuffer)[i] = encryptor.cipher.Encrypt(binaryByteArray)
	}
}

func (engine *DESFileEncryptor) decryptChunk(buffer *[]byte, bufferDataSize int, decryptionBuffer *[][]byte) {
	for i := 0; i < bufferDataSize; i++ {
		byteVal := (*buffer)[i]
		var binaryByteArray []byte = engine.getBinaryByteArray(byteVal)
		(*decryptionBuffer)[i] = engine.cipher.Decrypt(binaryByteArray)
	}
}

func (encryptor *DESFileEncryptor) writeEncryptionBufferToFile(encryptionBuffer *[][]byte,
	encryptionDataSize int,
	filename string) {
	// var byteVal byte
	fmt.Println("enccryption buffer length:", encryptionDataSize)
	var byteArray []byte = make([]byte, encryptionDataSize)
	var cache []byte = make([]byte, 8)
	for i := 0; i < encryptionDataSize; i++ {
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
	_, err = file.WriteString((string)(byteArray))

	file.Close()
}

func (engine *DESFileEncryptor) writeDecryptionBufferToFile(decryptionBuffer *[][]byte,
	decryptionDataSize int,
	filename string) {

	fmt.Println("Decryption buffer length:", decryptionDataSize)
	var byteArray []byte = make([]byte, decryptionDataSize)
	var cache []byte = make([]byte, 8)
	for i := 0; i < decryptionDataSize; i++ {
		copy(cache[:], (*decryptionBuffer)[i][:])
		byteVal := engine.convertBinaryByteArrayToByte(&cache)
		byteArray[i] = byteVal
	}

	_, err := engine.decryptionFileConnector.WriteString((string)(byteArray))
	if err != nil {
		log.Fatalln("Problem decrypting file", err)
	}
}

/**
Reads in a chunk of byte-data from file.
Encrypts it.
Writes the encrypted chunk of data in the encryption file
*/
func (encryptor *DESFileEncryptor) run(filename string) {

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln("Problem encrypting file", err)
		return
	}
	defer file.Close()

	//	Create encryption file
	_, err = os.Create(encryptor.encryptionFilename)
	if err != nil {
		log.Fatalln("Problem encrypting file", err)
	}

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

		encryptor.encryptChunk(&buffer, bytesread, &encryptionBuffer)
		//	Write encryptionBuffer into file.
		encryptor.writeEncryptionBufferToFile(&encryptionBuffer, bytesread, encryptor.encryptionFilename)

		fmt.Println("bytesread:", bytesread, "bytes to string:", string(buffer[:bytesread]))
	}
}

/**
Decrypt the bytes of the encrypted file
*/
func (engine *DESFileEncryptor) runDecryption(filename string) {
	var buffer []byte = make([]byte, bufferSize)
	var decryptionBuffer [][]byte = make([][]byte, bufferSize)
	for i := 0; i < bufferSize; i++ {
		decryptionBuffer[i] = make([]byte, 8)
	}

	fmt.Println("Encryption filename:", engine.encryptionFilename)
	file, ferr := os.Open(engine.encryptionFilename)
	if ferr != nil {
		log.Fatalln("Problem opening encrypted file", ferr)
	}
	for {
		bytesread, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Fatalln("Problem reading encrypted file", err)
			}
			break
		}
		engine.decryptChunk(&buffer, bytesread, &decryptionBuffer)
		//	Write decryptionBuffer into file
		engine.writeDecryptionBufferToFile(&decryptionBuffer, bytesread, engine.decryptionFilename)

		fmt.Println("Decryption, bytesread:", bytesread, " bytes to string:", string(buffer[:bytesread]))
	}
}

/**
public File-encryption API
*/
func (encryptor *DESFileEncryptor) EncryptFile(filename string) {
	log.Println("File-Encryption procedure started...")

	encryptor.filename = filename
	encryptor.encryptionFilename = filename + ".enc"
	encryptor.cipher.Init("des_input.txt")
	encryptor.run(filename)

	log.Println("File-Encryption procedure complete...")
}

/**
Public Decryption API
*/
func (engine *DESFileEncryptor) DecryptFile(filename string) {
	log.Println("Decryption procedure started...")

	engine.decryptionFilename = "dec." + engine.filename
	var err error
	engine.decryptionFileConnector, err = os.Create(engine.decryptionFilename)
	if err != nil {
		log.Fatalln("Problem decrypting the file", err)
	}
	permissions := 0644
	engine.decryptionFileConnector, err =
		os.OpenFile(engine.decryptionFilename, os.O_APPEND|os.O_WRONLY, (os.FileMode)(permissions))
	// file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, (os.FileMode)(permissions))
	engine.runDecryption(filename)
	engine.decryptionFileConnector.Close()

	log.Println("Decryption procedure complete...")
}
