package main

import "testing"
import "bytes"
// import "encoding/hex"

func TestPad(t *testing.T) {
    in := []byte("YELLOW SUBMA")
    out := append(in, []byte{4, 4, 4, 4}...)
    if !bytes.Equal(out, pad(in)) {
        t.Errorf("Failed to correctly implement PKCS7 padding, expected %b, got %b", out, pad(in))
    }
}

// func TestDecrypt(t *testing.T) {
//     key := []byte("example key 1234")
//     correct := []byte("exampleplaintext")
//     ciphertext, _ := hex.DecodeString("f363f3ccdcb12bb883abf484ba77d9cd7d32b5baecb3d4b1b3e0e4beffdb3ded")
//     clear := decrypt(key, ciphertext[16:], ciphertext[:16])
//     if !bytes.Equal(correct, clear) {
//         t.Errorf("Failed to properly decrypt message", correct, clear);
//     }
// }

// func TestEncrypt(t *testing.T) {
//     hexDump, _ := hex.DecodeString("f363f3ccdcb12bb883abf484ba77d9cd7d32b5baecb3d4b1b3e0e4beffdb3ded")
//     iv := hexDump[:16]
//     correct := hexDump[16:]

//     key := []byte("example key 1234")
//     clear := []byte("exampleplaintext")
//     ciphertext := encrypt(key, clear, iv)
//     if !bytes.Equal(correct, ciphertext) {
//         t.Errorf("Failed to properly encrypt message", correct, ciphertext)
//     }
// }

func TestBoth(t *testing.T) {
    iv := []byte{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
    clear := []byte("I am a message that needs to be transmitted far and wide to secret people");
    key := []byte("s3cr3t")

    ciphertext := encrypt(key, clear, iv)
    out := decrypt(key, ciphertext, iv)

    if !bytes.Equal(out, clear) {
        t.Errorf("Failed to properly encrypt then decrypt message, %b, %b", clear, out)
    }
}