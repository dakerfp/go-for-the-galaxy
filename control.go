package main

import (
	"time"
)

const animationSpeed = 10 * time.Millisecond

func gameControl(cmdQueue chan Command, game *Game, draw func(*Game) error) {
	const animationSpeed = 10 * time.Millisecond

	fallingTimer := time.NewTicker(animationSpeed)

	draw(game)

	var cmds []Command

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
