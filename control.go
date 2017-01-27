package main

import (
	"time"

	"github.com/nsf/termbox-go"
)

const animationSpeed = 10 * time.Millisecond

func termboxControl(game *Game, draw func(*Game)) {
	var cmds []Command
	var from *Planet

	fallingTimer := time.NewTicker(animationSpeed)
	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	draw(game)

mainloop:
	for {
		select {
		case ev := <-eventQueue:
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
			draw(game)

		default:
			draw(game)
		}
	}
}
