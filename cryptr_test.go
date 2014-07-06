package main

import "testing"
import "bytes"
import "log"

// import "encoding/hex"

func TestPad(t *testing.T) {
	in := []byte("YELLOW SUBMA")
	out := append(in, []byte{4, 4, 4, 4}...)
	if !bytes.Equal(out, pad(in)) {
		t.Errorf("Failed to correctly implement PKCS7 padding, expected %b, got %b", out, pad(in))
	}
}

func TestEncryptAndDecrypt(t *testing.T) {
	iv := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	clear := []byte("I am a message that needs to be transmitted far and wide to secret people")
	key := []byte("s3cr3t")

	ciphertext := encrypt(key, clear, iv)
	out := decrypt(key, ciphertext, iv)

	if !bytes.Equal(out, clear) {
		t.Errorf("Failed to properly encrypt then decrypt message, %b, %b", clear, out)
	}
}

func TestHash(t *testing.T) {
	message := []byte("i am a message")
	correct := fromBase64([]byte("pkId5Gox5klK09SiHVZRytUd0Qs/VfIqMYl/TohrX/g="))

	out := hash(message)
	if !bytes.Equal(correct, out) {
		t.Errorf("Failed to properly hash message, %b, %b", correct, out)
	}
}

func TestSignAndVerify(t *testing.T) {
	message := []byte("I am a very very very very very very very very long message")
	signature := sign(message)

	log.Println("SIGNATURE IS ", len(signature), signature)
	if !verify(message, signature, &PrivateKey.PublicKey) {
		t.Errorf("Failed to verify signed hash")
	}
}
