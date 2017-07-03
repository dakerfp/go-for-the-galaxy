package main

import (
	"encoding/gob"
	"flag"
	"io"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/dakerfp/go-for-the-galaxy"
)

var (
	serverFlag = flag.Bool("server", false, "decide if it will run as server or as client")
	clientFlag = flag.Bool("client", false, "decide if it will run as server or as client")
	addrFlag   = flag.String("addr", ":7771", "the server address")
	seedFlag   = flag.Int64("seed", 0, "random seed")
)

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	startGameServer(*addrFlag)
}

type GameRoom struct {
	GameId int
}

func (room *GameRoom) Serve(conns ...io.ReadWriteCloser) error {
	defer func() {
		for _, conn := range conns {
			defer conn.Close()
		}
	}()

	model := galaxy.NewRandomMap(60, 50)
	encs := make([]*gob.Encoder, len(conns))
	cmdQueue := make(chan galaxy.Command)
	for i, conn := range conns {
		enc := gob.NewEncoder(conn)
		encs[i] = enc
		if err := enc.Encode(galaxy.Player(i + 1)); err != nil {
			return err
		}

		go readCommands(conn, cmdQueue)
	}

	return model.Run(cmdQueue, func(game *galaxy.Game) error {
		for _, enc := range encs {
			if err := enc.Encode(*game); err != nil {
				return err
			}
		}
		return nil
	})
}

func readCommands(r io.Reader, cmds chan galaxy.Command) {
	defer close(cmds)

	dec := gob.NewDecoder(r)
	for {
		var cmd galaxy.Command
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

func startGameServer(addr string) {
	ln, err := net.Listen("tcp", addr)
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
		room := GameRoom{gameId}
		go room.Serve(conn1, conn2)
	}
}
