package main

import (
	"time"
)

const animationSpeed = 10 * time.Millisecond

type GameInterface interface {
	Player() (Player, error)
	Run(cmdQueue chan Command, draw func(*Game) error) error
	Probe(x, y int) *Planet
}

func (game *Game) Player() (Player, error) {
	return Player(1), nil // Local player
}

func (game *Game) Run(cmdQueue chan Command, draw func(*Game) error) error {
	fallingTimer := time.NewTicker(animationSpeed)

	var cmds []Command
	if err := draw(game); err != nil {
		return err
	}

	for {
		select {
		case cmd := <-cmdQueue:
			if cmd.CommandType == CommandQuit {
				break
			}
			cmds = append(cmds, cmd)

		case <-fallingTimer.C:
			game.Tick(cmds)
			cmds = nil
			// Sole player on
			if len(game.CountPlanetsByPlayer()) == 1 {
				return nil
			}
			if err := draw(game); err != nil {
				return err
			}
		}
	}
}

func (g *Game) Probe(x int, y int) *Planet {
	mouse := Vec2{float32(x), float32(y)}
	var min *Planet
	minDist := 100000.0
	for label, planet := range g.Planets {
		d := dist(mouse, planet.Center)
		if d < minDist {
			minDist = d
			min = planet
			min.Id = label // XXX
		}
	}
	return min
}
