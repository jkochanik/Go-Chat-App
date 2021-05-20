package main

import (
	"bufio"
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/sirupsen/logrus"
	//"encoding/json"
	"net"
	"os"
	"time"
)



type message struct {
	Name string
	Body string
}

// Application constants, defining host, port, and protocol.
const (
	connHost        = "raspberrypi.local"
	connPort        = "8080"
	connType        = "tcp"
)




func main() {
	fmt.Println("Connecting to", connType, "server", connHost+":"+connPort)
	conn, err := net.Dial(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}


	//makes channel for messages
	//mess := make(chan string)

	err1 := ui.Init()
	if err1 != nil {
		logrus.Fatalf("failed to initializing termui: %v, err")
	}
	defer ui.Close()

	//Creates prompt block
	t := widgets.NewParagraph()
	t.SetRect(0, 33, 30, 37)
	t.Text = "Hello this is a test"
	ui.Render(t)

	//Creates chat block
	l := widgets.NewList()
	l.Title = "chat"
	l.Rows = []string{
		"Chat Below",
	}
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false
	l.SetRect(0, 0, 30, 30)

	ui.Render(l);
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	input := ""

	//go checkForMessage(conn,mess)



	//Listens for inputs
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "<Escape>":
				ui.Close()
			default:
				if e.Type == ui.KeyboardEvent && len(e.ID) == 1 {
					input += e.ID
					t.Text = input
					ui.Render(t)
				} else if e.ID == "<Backspace>" {
					input = input[0 : len(input)-1]
					t.Text = input
					ui.Render(t)
				} else if e.ID == "<Space>" {
					input += " "
					t.Text = input
					ui.Render(t)
				} else if e.ID == "<Enter>" {
					conn.Write([]byte(input))
					//message, _ := bufio.NewReader(conn).ReadString('\n')
					//l.Rows = append(l.Rows, message)
					input = ""
					t.Text = input
					ui.Render(t)
				}
			}

		case <-ticker:
			//if mess != nil {
			//	l.Rows = append(l.Rows, <- mess)
			//}
		}
		ui.Render(l)
	}
}


func checkForMessage(Connection net.Conn, mess chan string ) {
	message, err :=  bufio.NewReader(Connection).ReadString('\n')
	if err != nil {
		ui.Close()
		fmt.Println("Lost Connection with Server")
		Connection.Close()
		return
	}
	mess <- message
}