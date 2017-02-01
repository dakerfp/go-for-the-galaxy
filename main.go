package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/nsf/termbox-go"
)

var (
	serverFlag = flag.Bool("server", false, "decide if it will run as server or as client")
	clientFlag = flag.Bool("client", false, "decide if it will run as server or as client")
	portFlag   = flag.Uint("port", 7771, "port")
)

func main() {
	flag.Parse()

	if *serverFlag {
		startGameServer(*portFlag)
		return
	}

	var model GameInterface
	if *clientFlag {
		conn, err := net.Dial("tcp", fmt.Sprintf(":%d", *portFlag))
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		model = &ProxyGame{RW: conn}
	} else {
		model = NewDefaultGame()
	}

	// Initializing termbox
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputMouse)

	// Setup & Run game
	player, err := model.Player()
	if err != nil {
		panic(err)
	}
	cmdQueue := make(chan Command)
	go termboxInput(player, model, cmdQueue)
	model.Run(cmdQueue, termboxDraw)
}
