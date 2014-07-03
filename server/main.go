package main 

import (
    "log"
    "net"
)

const (

    // HOST is the host we are listening on
    HOST = "localhost"

    // PORT is the portnumber we are listening on
    PORT = "1337"
)

var Connections[] net.Conn

func main() {
    l, err := net.Listen("tcp", HOST + ":" + PORT)
    if err != nil {
        log.Fatalln("Error listening: ", err.Error());
    }

    defer l.Close()

    log.Println("Listening on " + HOST + ":" + PORT);

    for {
        conn, err := l.Accept()
        if err != nil {
            log.Fatalln("Error accepting connection", err);
        }

        log.Println("Accepted a connection from", conn.RemoteAddr());
        Connections = append(Connections, conn);

        go listenOnConnection(conn)
    }
}

func echoOut(conn net.Conn, data []byte) {
    for i, c := range Connections {
        if (conn != c) {
            _, err := c.Write(data);
            if err != nil {
                log.Println("Received error from ", c, "writing message", data)
                Connections = append(Connections[:i], Connections[i+1:]...)
            }
        }
    }
}

func removeConnection(conn net.Conn) {
    index := 0
    for i, c := range Connections {
        if c == conn {
            index = i
        }
    }
    Connections = append(Connections[:index], Connections[index+1:]...)
}

func listenOnConnection(conn net.Conn) {
    data := make([]byte, 1024);
    defer conn.Close()

    for {
        for i := range data { data[i] = 0 }

        count, err := conn.Read(data)
        if err != nil {
            log.Println("Error reading from ", conn.RemoteAddr().String())
            removeConnection(conn)
            return
        }
        log.Printf("%s: (%d): %s", conn.RemoteAddr().String(), count, data);
        echoOut(conn, data[:count])
    }
}