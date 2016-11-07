package main

import (
	"encoding/json"
	"net/http"
	"io"
	"fmt"
	"log"
	"time"

	"github.com/nsf/termbox-go"
	"golang.org/x/net/websocket"
)

const animationSpeed = 10 * time.Millisecond

type GameListener chan *Game

func main() {
	// err := termbox.Init()
	// if err != nil {
	// 	panic(err)
	// }
	// defer termbox.Close()
	// termbox.SetInputMode(termbox.InputMouse)

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

	commandEvent := make(chan *Command)
	newListeners := make(chan GameListener, 10)

	http.Handle("/room", websocket.Handler(func(ws *websocket.Conn) {
		listener := make(GameListener)
		newListeners <- listener

		go func() {
			cmd := &Command{}
			dec := json.NewDecoder(ws)
			for {
				err := dec.Decode(cmd)
				if err == io.EOF {
					return
				} else if err != nil {
					log.Println(err)
					return
				}
				commandEvent <- cmd
			}
		}()

		enc := json.NewEncoder(ws)

		for game := range listener {
			if err := enc.Encode(&game); err != nil {
				log.Println(err)
				return
			}
		}
	}))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `<html>
<head>
</head>
<body>
<h1>Game</h1>
<p id="text"></p>
<script>
var socket = new WebSocket("ws://localhost:8888/room");
var text = document.getElementById("text");
socket.onmessage = function (e) {
  console.log(e.data);
  text.innerText = e.data + "\n";
};
</script>
</body>
</html>
`)
	})

	go func() {
	    err := http.ListenAndServe(":8888", nil)
	    if err != nil {
	        panic("ListenAndServe: " + err.Error())
	    }
	}()

	eventQueue := make(chan termbox.Event)
	// go func() {
	// 	for {
	// 		eventQueue <- termbox.PollEvent()
	// 	}
	// }()

	// draw(&game)

	var cmd *Command
	var from *Planet

	var listeners []GameListener
	for {
		select {
		case listener := <- newListeners:
			listeners = append(listeners, listener)

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

		case cmd = <-commandEvent:
			continue

		case <-game.fallingTimer.C:
			game.Tick(cmd)
			cmd = nil

		default:
			wsDraw(listeners, &game) // draw(&game)
			time.Sleep(animationSpeed)
		}
	}
}
