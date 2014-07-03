package main 

import (
    "net"
    "log"
    "fmt"
    "os"
    "bufio"
)

const (

    // RemoteHost the remote host we are connecting to
    RemoteHost = "localhost"

    // RemotePort the remote port we are connecting to
    RemotePort = "1337"

    // KEY is a SUPER SECRET PRE-SHARED KEY. `echo "1337h4x" | md5`
)

var KEY = []byte("3a989dba6fe6c87f")

func main() {

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

        count, err := conn.Read(data)
        if err != nil {
            log.Fatalln("Failed to read from connection", err);
        }
        data = decrypt(KEY, data[:count], IV)
        fmt.Printf("%s (%d): %s", conn.RemoteAddr().String(), len(data), data)
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
        data = encrypt(KEY, data, IV)

        conn.Write(data);
    }
}