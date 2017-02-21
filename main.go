package main

import (
	"flag"
	"math/rand"
	"net"

	"github.com/nsf/termbox-go"
)

var (
	serverFlag = flag.Bool("server", false, "decide if it will run as server or as client")
	clientFlag = flag.Bool("client", false, "decide if it will run as server or as client")
	addrFlag   = flag.String("addr", ":7771", "the server address")
	seedFlag   = flag.Int64("seed", 0, "random seed")
)

func main() {
	flag.Parse()
	rand.Seed(*seedFlag)

	if *serverFlag {
		startGameServer(*addrFlag)
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
		conn, err := net.Dial("tcp", *addrFlag)
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
