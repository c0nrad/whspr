package main

import (
    tm "github.com/buger/goterm"
)

var Messages [][]byte
var WindowHeight int 
var WindowWidth int
var DebugLine []byte

func initDisplay() {
    WindowHeight = tm.Height()
    WindowWidth = tm.Width()
    addMessage([]byte("WELCOME TO WHSPR"))

    render()
}

func getInputBoxY() int {
    return WindowHeight - 2
}

func getDebugLineX() int {
    return WindowHeight 
}

func getDisplayBoxHeight() int {
    return getInputBoxY() - 1
}

func setDebugLine(text []byte, color int) {
    DebugLine = text
    tm.MoveCursor(getDebugLineX(), 1)
    tm.Printf("%s", tm.Color(string(text), color))
    tm.MoveCursor(2,2)
    //tm.Flush()
}

func addMessage(text []byte) {
    Messages = append(Messages, text)
}

func getMessages() [][]byte {
    count := getDisplayBoxHeight() - 2
    
    if len(Messages) < count {
        return Messages
    }

    return Messages[len(Messages) - count:]
}

func render() {
    tm.Clear()

    tm.MoveCursor(1, 1)
    box := tm.NewBox(WindowWidth, getDisplayBoxHeight() , 0)
    messages := getMessages()
    for _, m := range messages {
        box.Write(m)
        box.Write([]byte{ byte('\n') })
    }
    tm.Println(tm.Color(box.String(), tm.GREEN))


    setDebugLine([]byte("Rendered..."), tm.WHITE)

    //tm.Flush()
    tm.MoveCursor(getInputBoxY() - 1, 3)
    tm.Flush()
    tm.MoveCursor(getInputBoxY() - 1, 3)
}