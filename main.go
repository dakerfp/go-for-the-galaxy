package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"

	"github.com/nsf/termbox-go"
)

var (
	serverFlag = flag.Bool("server", false, "decide if it will run as server or as client")
	clientFlag = flag.Bool("client", false, "decide if it will run as server or as client")
	portFlag   = flag.Uint("port", 7771, "port")
	seedFlag   = flag.Int64("seed", 0, "random seed")
)

func main() {
	flag.Parse()
	rand.Seed(*seedFlag)

	if *serverFlag {
		startGameServer(*portFlag)
		return
	}

	// Initializing termbox
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputMouse)

	var model GameInterface
	if *clientFlag {
		conn, err := net.Dial("tcp", fmt.Sprintf(":%d", *portFlag))
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		model = &ProxyGame{RW: conn}
	} else {
		w, h := termbox.Size()
		model = NewRandomMap(float32(w), float32(h))
	}

	// Setup & Run game
	player, err := model.Player()
	if err != nil {
		panic(err)
	}
	cmdQueue := make(chan Command)
	go termboxInput(player, model, cmdQueue)
	model.Run(cmdQueue, termboxDraw)
}
