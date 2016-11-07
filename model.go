package main

import (
	"math"
	"time"
)

type Player int

type Vec2 struct {
	X float32
	Y float32
}

func add(a Vec2, b Vec2) Vec2 {
	return Vec2{a.X + b.X, a.Y + b.Y}
}

func sub(a Vec2, b Vec2) Vec2 {
	return Vec2{a.X - b.X, a.Y - b.Y}
}

func mult(a Vec2, k float32) Vec2 {
	return Vec2{a.X * k, a.Y * k}
}

func dist(a Vec2, b Vec2) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return math.Sqrt(float64(dx*dx + dy*dy))
}

type Planet struct {
	Id     string
	Center Vec2
	Size   float32
	Units  float32
	Player Player
}

type Fleet struct {
	Player Player
	To     string
	Units  float32
	Pos    Vec2
	Vel    Vec2
	Dest   Vec2
	Dead   bool
}

type Game struct {
	Planets      map[string]*Planet
	Fleets       []Fleet
	Winner       Player
	fallingTimer *time.Ticker
}

type Command struct {
	From  string
	To    string
	Units float32
}

func (g *Game) Tick(cmd *Command) {
	if cmd != nil {
		from, ok := g.Planets[cmd.From]
		to, ok2 := g.Planets[cmd.To]
		if ok && ok2 {
			if cmd.Units < from.Units {
				fleet := Fleet{
					Player: from.Player,
					To:     cmd.To,
					Pos:    from.Center,
					Units:  cmd.Units,
					Vel:    mult(sub(to.Center, from.Center), 0.01),
					Dest:   to.Center,
				}
				g.Fleets = append(g.Fleets, fleet)
				from.Units -= cmd.Units
				g.Planets[cmd.From] = from
			}
		}
	}

	oldFleets := g.Fleets
	newFleets := g.Fleets[:0]
	for _, fleet := range oldFleets {
		fleet.Pos = add(fleet.Pos, fleet.Vel)
		if dist(fleet.Pos, fleet.Dest) > 1 {
			newFleets = append(newFleets, fleet)
			continue
		}
		planet, ok := g.Planets[fleet.To]
		if !ok {
			panic(cmd.To)
		}
		if fleet.Player == planet.Player {
			planet.Units += fleet.Units
		} else {
			planet.Units -= fleet.Units
			if planet.Units < 0 {
				planet.Units = -planet.Units
				planet.Player = fleet.Player
			}
		}
	}
	g.Fleets = newFleets

	for _, planet := range g.Planets {
		if planet.Player != 0 {
			planet.Units += planet.Size * 0.05 // Speed multiplier
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
