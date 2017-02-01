package main

import (
	"encoding/gob"
	"fmt"
	"github.com/nsf/termbox-go"
	"io"
	"log"
	"net"
)

func readCommands(r io.Reader, cmds chan Command) {
	dec := gob.NewDecoder(r)
	for {
		var cmd Command
		err := dec.Decode(&cmd)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			return
		}
		cmds <- cmd
	}
}

func serveGames(port uint) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return
	}

	for gameId := 0; ; gameId++ {
		log.Println("Waiting player 1")
		conn1, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		log.Println("Waiting player 2")
		conn2, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go serveGame(conn1, conn2, gameId)
	}
}

func serveGame(conn1 io.ReadWriteCloser, conn2 io.ReadWriteCloser, gameId int) error {
	defer log.Println("Ending game ", gameId)
	defer conn1.Close()
	defer conn2.Close()

	log.Println("Starting game ", gameId)

	model := NewRandomMap(60, 50)
	enc1 := gob.NewEncoder(conn1)
	enc2 := gob.NewEncoder(conn2)

	if err := enc1.Encode(Player(1)); err != nil {
		return err
	}
	if err := enc2.Encode(Player(2)); err != nil {
		return err
	}

	draw := func(game *Game) error {
		if err := enc1.Encode(*game); err != nil {
			return err
		}
		if err := enc2.Encode(*game); err != nil {
			return err
		}
		return nil
	}

	cmdQueue := make(chan Command)

	go readCommands(conn1, cmdQueue)
	go readCommands(conn2, cmdQueue)

	model.Run(cmdQueue, draw)

	return nil
}

func clientGame(port uint, draw func(*Game) error) error {
	err := termbox.Init()
	if err != nil {
		return err
	}

	defer termbox.Close()
	termbox.SetInputMode(termbox.InputMouse)

	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	var player Player
	var game Game
	dec := gob.NewDecoder(conn)
	err = dec.Decode(&player)
	if err != nil {
		return err
	}

	err = dec.Decode(&game)
	if err != nil {
		return err
	}

	err = draw(&game)
	if err != nil {
		return err
	}

	cmdQueue := make(chan Command)
	go termboxInput(player, &game, cmdQueue)
	go netSendCommand(conn, cmdQueue)

	for {
		err := dec.Decode(&game)
		switch err {
		case nil:
			err = draw(&game)
			if err != nil {
				return err
			}
		case io.EOF:
			return nil
		default:
			return err
		}
	}

	return nil
}

type ProxyGame struct {
	conns []net.Conn
}

func (p *ProxyGame) Tick(cmds []Command) {}
