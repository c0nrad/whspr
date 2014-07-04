package main 

import (
    "crypto/aes"
    "crypto/cipher"
    "log"
    "bytes"
    "encoding/base64"
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

    IV = ciphertext[len(ciphertext)-16:len(ciphertext)]

    ciphertext = toBase64(ciphertext)
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
    ciphertext = fromBase64(ciphertext)

    block, err := aes.NewCipher(pad(key))
    if err != nil {
        log.Fatalln("Failed to build aes cipher block", err);
    }

    if len(ciphertext) % aes.BlockSize != 0 {
        log.Fatalln("decrypt cipher text not correct size", len(ciphertext), ciphertext, aes.BlockSize);
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

func toBase64(data []byte) []byte {
    out := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
    base64.StdEncoding.Encode(out, data)
    return out
}

func fromBase64(data []byte) []byte {
    out := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
    n, err := base64.StdEncoding.Decode(out, data)

    if err != nil {
        log.Fatalln("Error decoding base64 in decrypt", err, n)
    }

    out = bytes.TrimRight(out, "\x00")
    return out
}