package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"net"
	"os"
)

func main() {
	parseTrustedCerts()

	conn, err := net.Dial("tcp", RemoteHost+":"+RemotePort)
	if err != nil {
		log.Fatalln("Failed to connect to remote host", err)
	}

	defer conn.Close()

	go connToStdout(conn)
	go stdinToConn(conn)

	select {}
}

func connToStdout(conn net.Conn) {

	for {
		data := make([]byte, 1024)

		n, err := conn.Read(data)
		if err != nil {
			log.Fatalln("Failed to read from connection", err)
		}

		var m Message
		err = json.Unmarshal(data[:n], &m)
		if err != nil {
			log.Fatalln("Failed to unmarshal json", err)
		}

		data = ParseMessage(m)
		addMessage(data)
		render()
	}
}

func stdinToConn(conn net.Conn) {
	in := bufio.NewReader(os.Stdin)
	for {
		data := make([]byte, 1024)

		data, err := in.ReadBytes('\n')
		if err != nil {
			log.Fatalln("Error reading from Stdin", err)
		}

		data = bytes.TrimRight(data, "\n")

		if len(data) == 0 {
			render()
			continue
		}

		m := NewMessage(data)
		out, err := json.Marshal(m)

		if err != nil {
			log.Fatalln("Error marshaling json", err)
		}

		conn.Write(out)
		render()
	}
}
