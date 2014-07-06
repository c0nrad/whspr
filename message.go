package main

// Message struct being sent over the wire
type Message struct {
	IV        []byte
	Data      []byte
	Signature []byte
	Name      []byte
}

func NewMessage(data []byte) Message {
	iv := IV[:]

	data = data
	ciphertext := encrypt([]byte(KEY), data, iv)
	signature := sign(ciphertext)
	name := Name

	m := Message{toBase64(iv), toBase64(ciphertext), toBase64(signature), toBase64(name)}
	return m
}

func ParseMessage(m Message) []byte {
	m.IV = fromBase64(m.IV, true)
	m.Data = fromBase64(m.Data, true)
	m.Signature = fromBase64(m.Signature, true)
	m.Name = fromBase64(m.Name, false)

	valid := verify(m.Data, m.Signature, getPublicKey(m.Name))
	if valid {
		data := decrypt([]byte(KEY), m.Data, m.IV)
		out := append(m.Name, ": "...)
		out = append(out, data...)
		return out
	} else {
		return []byte("")
	}
}
