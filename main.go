package main

import (
	"flag"
	"log"

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
		serveGames(*portFlag)
		return
	} else if *clientFlag {
		err := clientGame(*portFlag, termboxDraw)
		if err != nil {
			log.Println(err)
		}
		return
	}

	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	defer termbox.Close()
	termbox.SetInputMode(termbox.InputMouse)

	model := NewDefaultGame()

	termboxControl(model, termboxDraw)
}
