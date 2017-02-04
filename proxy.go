package main

import (
	"encoding/gob"
	"io"
)

type ProxyGame struct {
	RW   io.ReadWriteCloser
	game Game
}

func (p *ProxyGame) Player() (Player, error) {
	dec := gob.NewDecoder(p.RW)
	var player Player
	err := dec.Decode(&player)
	if err != nil {
		return 0, err
	}
	return player, nil
}

func (p *ProxyGame) Run(cmdQueue chan Command, draw func(*Game) error) error {
	dec := gob.NewDecoder(p.RW)

	go func() { // Send commands
		enc := gob.NewEncoder(p.RW)
		for cmd := range cmdQueue {
			if err := enc.Encode(cmd); err != nil {
				return
			}
		}
	}()

	for {
		err := dec.Decode(&p.game)
		switch err {
		case nil:
			if err = draw(&p.game); err != nil {
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

func (p *ProxyGame) Probe(x, y int) *Planet {
	return p.game.Probe(x, y)
}
