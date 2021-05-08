// Package goshared ...
package goshared

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
)

func Ternary(b bool, t, f interface{}) interface{} {
	if b {
		return t
	}
	return f
}

func ParseQuery(str string) (url.Values, error) {
	URLStr, e := url.QueryUnescape(str)
	if e != nil {
		return nil, e
	}
	u, e := url.Parse(URLStr)
	if e != nil {
		return nil, e
	}
	return u.Query(), nil
}

func FlatPath(path string) []string {
	var r []string
	// First element
	r = append(r, path[0:1])
	// Middle elements
	tmpPath := strings.Replace(path, "/", " ", 1)
	for {
		offset := strings.Index(tmpPath, "/")
		if offset != -1 {
			r = append(r, path[:offset])
			tmpPath = strings.Replace(tmpPath, "/", " ", 1)
		} else {
			break
		}
	}
	// Latest element
	r = append(r, path)
	return r
}

func AES256Encrypt(stringToEncrypt string, keyString string) (encryptedString string) {
	//Since the key is in string, we need to convert decode it to bytes
	key, _ := hex.DecodeString(keyString)
	if len(key) != 32 { // AES256 key length is 32
		panic("AES256 encript: key string error")
	}
	plaintext := []byte(stringToEncrypt)
	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

func AES256Decrypt(encryptedString string, keyString string) (decryptedString string) {
	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)
	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	//Get the nonce size
	nonceSize := aesGCM.NonceSize()
	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return fmt.Sprintf("%s", plaintext)
}

// IOMV Solve problem that os.Rename() give error "invalid cross-device link"
// in Docker container with Volumes(container FS and volume FS are different).
// Note that you need have permissions on both FS.
func IOMV(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("failed removing original file: %s", err)
	}
	return nil
}

// RepeatedlyDo do some operation at least once
// @Param op represent operation function which has 'func() error' signature
// @Param rt represent repeated times
func RepeatedlyDo(op func() error, rt uint) error {
	var count uint = 0
	var e error
	for e = op(); e != nil && count < rt; count++ {
		e = op()
		if e == nil {
			return nil
		}
	}
	return e
}
