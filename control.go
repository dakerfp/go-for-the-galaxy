package main

import (
	"time"

	"github.com/nsf/termbox-go"
)

const animationSpeed = 10 * time.Millisecond

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputMouse)

	game := Game{
		Planets: map[string]*Planet{
			"a": &Planet{
				Center: Vec2{5, 5},
				Player: 1,
				Units:  100,
				Size:   2,
			},
			"b": &Planet{
				Center: Vec2{20, 5},
				Player: 0,
				Units:  100,
				Size:   1,
			},
			"c": &Planet{
				Center: Vec2{18, 15},
				Player: 0,
				Units:  100,
				Size:   1,
			},
			"d": &Planet{
				Center: Vec2{10, 24},
				Player: 0,
				Units:  50,
				Size:   1,
			},
			"e": &Planet{
				Center: Vec2{30, 15},
				Player: 2,
				Units:  50,
				Size:   2,
			},
			"f": &Planet{
				Center: Vec2{40, 2},
				Player: 0,
				Units:  50,
				Size:   1,
			},
			"g": &Planet{
				Center: Vec2{41, 40},
				Player: 0,
				Units:  50,
				Size:   1,
			},
		},
		fallingTimer: time.NewTicker(animationSpeed),
	}

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	draw(&game)

	var cmd *Command
	var from *Planet
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
						cmd = &Command{from.Id, to.Id, 50}
						from = nil
					}
				}
			}

		case <-game.fallingTimer.C:
			game.Tick(cmd)
			cmd = nil

		default:
			draw(&game)
			time.Sleep(animationSpeed)
		}
	}
}
