package main 

import (
    "net"
    "log"
    // "fmt"
    "os"
    "bufio"
    "encoding/json"
    "bytes"
)

const (

    // RemoteHost the remote host we are connecting to
    RemoteHost = "localhost"

    // RemotePort the remote port we are connecting to
    RemotePort = "1337"

)

// Message struct being sent over the wire 
type Message struct {
    IV []byte
    Data []byte
}

// KEY is a SUPER SECRET PRE-SHARED KEY. `echo "1337h4x" | md5`
var KEY = []byte("3a989dba6fe6c87f")

func main() {
    initDisplay()

    conn, err := net.Dial("tcp", RemoteHost + ":" + RemotePort);
    if err != nil {
        log.Fatalln("Failed to connect to remote host", err);
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
            log.Fatalln("Failed to read from connection", err);
        }

        var m Message
        err = json.Unmarshal(data[:n], &m)
        if err != nil {
            log.Fatalln("Failed to unmarshal json", err)
        }
        m.IV = fromBase64(m.IV)
        IV = m.IV

        data = decrypt(KEY, m.Data, m.IV)
        addMessage(data)
        render()
    }
}

func stdinToConn(conn net.Conn) {
    in := bufio.NewReader(os.Stdin)
    for {
        data := make([]byte, 1024)

        data, err := in.ReadBytes('\n');
        if err != nil {
            log.Fatalln("Error reading from Stdin", err);
        }

        data = bytes.TrimRight(append(Username, data...), "\n")

        iv := IV[:]
        data = encrypt(KEY, data, iv)
        m := Message{toBase64(iv), data}

        out, err := json.Marshal(m)

        if err != nil {
            log.Fatalln("Error marshaling json", err)
        }

        conn.Write(out);
        render()
    }
}