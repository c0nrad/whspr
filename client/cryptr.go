package main 

import (
    "crypto/aes"
    "crypto/cipher"
    "log"
    "bytes"
)

var IV = []byte{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}

//AES128-CBC
func encrypt(key, message, iv []byte) []byte {
    block, err := aes.NewCipher(pad(key))
    if err != nil {
        log.Fatalln("Failed to build aes cipher block", err);
    }
    message = pad(message)

    mode := cipher.NewCBCEncrypter(block, iv)

    ciphertext := make([]byte, len(message))
    mode.CryptBlocks(ciphertext, message)

    return ciphertext
}


// PKCS#7
func pad(message []byte) []byte {
    count := len(message)
    diff := 0
    for (count + diff) % aes.BlockSize != 0 {
        diff++;
    }

    out := append(message, bytes.Repeat([]byte{byte(diff)}, diff)...)
    return out
}

func decrypt(key, ciphertext, iv []byte) []byte {
    block, err := aes.NewCipher(pad(key))
    if err != nil {
        log.Fatalln("Failed to build aes cipher block", err);
    }

    if len(ciphertext) % aes.BlockSize != 0 {
        log.Fatalln("decrypt cipher text not correct size", len(ciphertext), aes.BlockSize);
    }

    mode := cipher.NewCBCDecrypter(block, iv)

    plain := make([]byte, len(ciphertext))
    mode.CryptBlocks(plain, ciphertext)

    return unpad(plain)
}

func unpad(message []byte) []byte {
    count := message[len(message)-1]
    if count >= 16 {
        return message;
    }

    i := byte(len(message)) - count
    for i < byte(len(message)) {
        if message[i] != count {
            return message;
        }
        i++;
    }
    return message[:len(message)-int(count)]
}