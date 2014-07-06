package main

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"log"
)

var IV []byte

func init() {
	IV = make([]byte, 16)

	_, err := rand.Read(IV)
	if err != nil {
		log.Fatalln("Error randomizing IV:", err)
	}
}

//AES128-CBC
func encrypt(key, message, iv []byte) []byte {
	block, err := aes.NewCipher(pad(key))
	if err != nil {
		log.Fatalln("Failed to build aes cipher block", err)
	}
	message = pad(message)

	mode := cipher.NewCBCEncrypter(block, iv)

	ciphertext := make([]byte, len(message))
	mode.CryptBlocks(ciphertext, message)

	IV = ciphertext[len(ciphertext)-16 : len(ciphertext)]

	return ciphertext
}

// PKCS#7
func pad(message []byte) []byte {
	count := len(message)
	diff := 0
	for (count+diff)%aes.BlockSize != 0 {
		diff++
	}

	out := append(message, bytes.Repeat([]byte{byte(diff)}, diff)...)
	return out
}

func decrypt(key, ciphertext, iv []byte) []byte {

	block, err := aes.NewCipher(pad(key))
	if err != nil {
		log.Fatalln("Failed to build aes cipher block", err)
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		log.Fatalln("decrypt cipher text not correct size", len(ciphertext), ciphertext, aes.BlockSize)
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	plain := make([]byte, len(ciphertext))
	mode.CryptBlocks(plain, ciphertext)

	return unpad(plain)
}

func unpad(message []byte) []byte {
	count := message[len(message)-1]
	if count >= 16 {
		return message
	}

	i := byte(len(message)) - count
	for i < byte(len(message)) {
		if message[i] != count {
			return message
		}
		i++
	}
	return message[:len(message)-int(count)]
}

func toBase64(data []byte) []byte {
	out := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(out, data)
	return out
}

func fromBase64(data []byte, round bool) []byte {
	out := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(out, data)

	if err != nil {
		log.Fatalln("Error decoding base64 in decrypt", err, n)
	}

	if round {
		out = roundBytes(out, 16)
	}

	return out
}

func roundBytes(data []byte, size int) []byte {
	offset := len(data) % size
	return data[0 : len(data)-offset]
}

// SHA256
func hash(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	raw := hash.Sum(nil) //so much, so raw
	return raw
}

var PrivateKey = parsePrivateKey("./certs/key.pem")
var HashType = crypto.SHA256

func sign(data []byte) (sig []byte) {

	hashed := hash(data)
	signature, err := rsa.SignPKCS1v15(rand.Reader, PrivateKey, HashType, hashed)
	if err != nil {
		log.Fatalln("Error signing hash", err)
	}
	return signature
}

func verify(data, signature []byte, pubKey *rsa.PublicKey) bool {
	hashed := hash(data)
	err := rsa.VerifyPKCS1v15(pubKey, HashType, hashed, signature)
	if err != nil {
		log.Println("Invalid signature", err)
		return false
	}

	return true
}
