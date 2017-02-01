package main

import (
	"encoding/gob"
	"fmt"
	"github.com/nsf/termbox-go"
	"io"
	"log"
	"net"
	"time"
)

func tcpListenCommands(conn net.Conn, cmds chan Command) {
	dec := gob.NewDecoder(conn)
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

func tcpControl(conn1 net.Conn, conn2 net.Conn, game *Game, draw func(*Game) error) {
	const animationSpeed = 10 * time.Millisecond
	var cmds []Command

	fallingTimer := time.NewTicker(animationSpeed)
	cmdQueue := make(chan Command)
	go tcpListenCommands(conn1, cmdQueue)
	go tcpListenCommands(conn2, cmdQueue)

	draw(game)

mainloop:
	for {
		select {
		case cmd := <-cmdQueue:
			cmds = append(cmds, cmd)

		case <-fallingTimer.C:
			game.Tick(cmds)
			cmds = nil
			// Sole player on
			if len(game.CountPlanetsByPlayer()) == 1 {
				break mainloop
			}
			draw(game)
		}
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

func serveGame(conn1 net.Conn, conn2 net.Conn, gameId int) {
	defer log.Println("Ending game ", gameId)
	log.Println("Starting game ", gameId)

	// model := NewSquareGameMap()
	model := NewRandomMap(60, 50)

	enc1 := gob.NewEncoder(conn1)
	enc2 := gob.NewEncoder(conn2)

	if err := enc1.Encode(Player(1)); err != nil {
		panic(err)
	}
	if err := enc2.Encode(Player(2)); err != nil {
		panic(err)
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

	log.Println("Starting game")
	tcpControl(conn1, conn2, model, draw)

	if err := conn1.Close(); err != nil {
		log.Println(err)
	}
	if err := conn2.Close(); err != nil {
		log.Println(err)
	}
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

	cmdQueue := make(chan Command)
	go func() {
		defer conn.Close()
		enc := gob.NewEncoder(conn)
		for cmd := range cmdQueue {
			if err := enc.Encode(cmd); err != nil {
				return
			}
		}
	}()

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

	go termboxInput(player, &game, cmdQueue)

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
