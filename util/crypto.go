package util

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func EncryptFile(file string, key []byte, extension string) {
	log.Printf("Encrypting file %s ...\n", file)

	plaintext, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	// The key should be 16 bytes (AES-128), 24 bytes (AES-192) or
	// 32 bytes (AES-256)
	// key, err := ioutil.ReadFile(fileKey)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panic(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Panic(err)
	}

	// Never use more than 2^32 random nonces with a given key
	// because of the risk of repeat.
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	// Save back to file
	fileOut := fmt.Sprintf("%s%s", file, extension)
	log.Printf("Writing file %s ...\n", fileOut)
	err = os.WriteFile(fileOut, ciphertext, 0777)
	if err != nil {
		log.Panic(err)
	}
}

func DecryptFile(file string, key []byte, extension string) {
	log.Printf("Decrypting file %s ...\n", file)

	ciphertext, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	// The key should be 16 bytes (AES-128), 24 bytes (AES-192) or
	// 32 bytes (AES-256)
	// key, err := ioutil.ReadFile(fileKey)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panic(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Panic(err)
	}

	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Panic(err)
	}
	fileOut := fmt.Sprintf("%s%s", file, extension)
	log.Printf("Writing file %s ...\n", fileOut)
	err = os.WriteFile(fileOut, plaintext, 0777)
	if err != nil {
		log.Panic(err)
	}
}

func GetConfigFromEncryptedFile(file string, key []byte) (map[string]string, error) {
	plainText, err := decryptFileToPlainText(file, key)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(string(plainText[:])))
	scanner.Split(bufio.ScanLines)
	var tmp string
	for scanner.Scan() {
		tmp = scanner.Text()
		pair := strings.SplitN(tmp, "=", 2)
		if len(pair) == 2 {
			key := pair[0]
			value := pair[1]
			m[key] = value
		}
	}
	return m, nil
}

func decryptFileToPlainText(file string, key []byte) ([]byte, error) {
	ciphertext, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
