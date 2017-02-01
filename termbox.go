package main

import (
	"github.com/nsf/termbox-go"

	"strconv"
)

func termboxInput(player Player, game *Game, cmds chan Command) {
	defer close(cmds)
	var from *Planet
	for {
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC, termbox.KeyCtrlD:
				cmds <- Command{CommandType: CommandQuit}
				return
			}

		case termbox.EventMouse:
			if ev.Key == termbox.MouseRelease {
				if from == nil {
					from = game.Probe(ev.MouseX, ev.MouseY)
				} else {
					to := game.Probe(ev.MouseX, ev.MouseY)
					cmds <- Command{CommandSendFleet, from.Id, to.Id, from.Units / 2, player}
					from = nil
				}
			}
		}
	}
}

const backgroundColor = termbox.ColorBlack

var playerColors = []termbox.Attribute{
	termbox.ColorYellow,
	termbox.ColorBlue,
	termbox.ColorRed,
}

func termboxDraw(g *Game) error {
	termbox.Clear(backgroundColor, backgroundColor)
	for _, fleet := range g.Fleets {
		if fleet.Dead {
			continue
		}
		x := int(fleet.Pos.X)
		y := int(fleet.Pos.Y)
		termbox.SetCell(x, y, '*', playerColors[fleet.Player], backgroundColor)
	}
	for name, planet := range g.Planets {
		x := int(planet.Center.X)
		y := int(planet.Center.Y)
		termbox.SetCell(x, y, rune(name[0]), playerColors[planet.Player], backgroundColor)

		s := strconv.Itoa(int(planet.Units))
		for _, r := range s {
			x++
			termbox.SetCell(x, y, r, playerColors[planet.Player], backgroundColor)
		}
	}
	termbox.Flush()

	return nil
}