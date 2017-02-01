package main

import (
	"time"
)

const animationSpeed = 10 * time.Millisecond

type GameInterface interface {
	Run(cmdQueue chan Command, draw func(*Game) error)
}

func (game *Game) Run(cmdQueue chan Command, draw func(*Game) error) {
	fallingTimer := time.NewTicker(animationSpeed)

	var cmds []Command
	draw(game)

mainloop:
	for {
		select {
		case cmd := <-cmdQueue:
			switch cmd.CommandType {
			case CommandSendFleet:
				cmds = append(cmds, cmd)
			case CommandQuit:
				break mainloop
			}

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
