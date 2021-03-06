package galaxy

import (
	"math"
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

func size(v Vec2) float64 {
	dx := v.X
	dy := v.Y
	return math.Sqrt(float64(dx*dx + dy*dy))
}

func norm(v Vec2) Vec2 {
	s := float32(size(v))
	if s <= 0 {
		return Vec2{}
	}
	return Vec2{v.X / s, v.Y / s}
}

type Planet struct {
	Id       string
	Center   Vec2
	Size     float32
	Units    float32
	Capacity float32
	Player   Player
	Links    []Command
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
	Planets map[string]*Planet
	Fleets  []Fleet
	Winner  Player
}

type CommandType int

const (
	CommandSendFleet CommandType = iota
	CommandCreateLink
	CommandDestroyLink
	CommandQuit
)

type Command struct {
	CommandType
	From   string
	To     string
	Units  float32
	Player Player
}

func (g *Game) SendFleet(cmd Command) bool {
	to, ok := g.Planets[cmd.To]
	if !ok {
		return false
	}

	from, ok := g.Planets[cmd.From]
	if cmd.Player != from.Player {
		return false
	}

	if cmd.Units > from.Units {
		return false
	}

	fleet := Fleet{
		Player: from.Player,
		To:     cmd.To,
		Pos:    from.Center,
		Units:  cmd.Units,
		Vel:    mult(norm(sub(to.Center, from.Center)), 0.1),
		Dest:   to.Center,
	}
	g.Fleets = append(g.Fleets, fleet)
	from.Units -= cmd.Units
	g.Planets[cmd.From] = from
	return true
}

func (g *Game) CreateLink(cmd Command) {
	_, ok := g.Planets[cmd.To]
	if !ok {
		return
	}

	from, ok := g.Planets[cmd.From]
	if !ok {
		return
	}

	if cmd.Player != from.Player {
		return
	}

	cmd.CommandType = CommandSendFleet
	from.Links = append(from.Links, cmd)
}

func (g *Game) DestroyLink(cmd Command) {
	_, ok := g.Planets[cmd.To]
	if !ok {
		return
	}

	from, ok := g.Planets[cmd.From]
	if !ok {
		return
	}

	if cmd.Player != from.Player {
		return
	}

	var newLinks []Command
	for _, link := range from.Links {
		if link.From != cmd.From ||
			link.To != cmd.To {
			newLinks = append(newLinks, link)
		}
	}
	from.Links = newLinks
}

func (g *Game) Tick(cmds []Command) {
	for _, cmd := range cmds {
		switch cmd.CommandType {
		case CommandSendFleet:
			g.SendFleet(cmd)
		case CommandCreateLink:
			g.CreateLink(cmd)
		case CommandDestroyLink:
			g.DestroyLink(cmd)
		}
	}

	for _, planet := range g.Planets {
		var newLinks []Command
		for _, cmd := range planet.Links {
			if g.SendFleet(cmd) {
				newLinks = append(newLinks, cmd)
			}
		}
		planet.Links = newLinks
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
			panic(fleet)
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
			planet.Units += planet.Size * 0.05  // Speed multiplier
			if planet.Units > planet.Capacity { // If planet capacity is exceeded, then it is neutralized
				planet.Player = 0
			}
		}
	}
}

func (g *Game) CountPlanetsByPlayer() map[Player]int {
	count := make(map[Player]int)
	for _, planet := range g.Planets {
		if planet.Player != 0 {
			count[planet.Player] += 1
		}
	}
	return count
}
