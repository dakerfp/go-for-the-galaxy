package main

import (
	"github.com/nsf/termbox-go"

	"strconv"
	"time"
)

func termboxControl(game *Game, draw func(*Game) error) error {
	var cmds []Command
	var from *Planet

	fallingTimer := time.NewTicker(animationSpeed)
	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	err := draw(game)
	if err != nil {
		return err
	}

mainloop:
	for {
		select {
		case ev := <-eventQueue:
			switch ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyEsc, termbox.KeyCtrlC, termbox.KeyCtrlD:
					return nil
				}

			case termbox.EventMouse:
				if ev.Key == termbox.MouseRelease {
					if from == nil {
						from = game.Probe(ev.MouseX, ev.MouseY)
					} else {
						to := game.Probe(ev.MouseX, ev.MouseY)
						cmds = []Command{Command{from.Id, to.Id, 50, 1}}
						from = nil
					}
				}
			}

		case <-fallingTimer.C:
			game.Tick(cmds)
			cmds = nil
			// Sole player on
			if len(game.CountPlanetsByPlayer()) == 1 {
				break mainloop
			}
			err := draw(game)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

var from *Planet

func termboxInput(player Player, game *Game, cmds chan Command) {
	defer close(cmds)
	for {
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC, termbox.KeyCtrlD:
				return
			}

		case termbox.EventMouse:
			if ev.Key == termbox.MouseRelease {
				if from == nil {
					from = game.Probe(ev.MouseX, ev.MouseY)
				} else {
					to := game.Probe(ev.MouseX, ev.MouseY)
					cmds <- Command{from.Id, to.Id, from.Units / 2, player}
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
