package main

import (
	"encoding/json"
	"flag"
	"io"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/dakerfp/go-for-the-galaxy"
	"github.com/nsf/termbox-go"
)

var (
	addrFlag = flag.String("addr", ":7771", "the server address")
)

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	// Initializing termbox
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputMouse)

	conn, err := net.Dial("tcp", *addrFlag)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	model := &ProxyGame{RW: conn}

	// Setup & Run game
	player, err := model.Player()
	if err != nil {
		panic(err)
	}

	cmdQueue := make(chan galaxy.Command)
	draw := make(chan galaxy.Game)
	go termboxInput(player, model, cmdQueue)
	go model.Run(cmdQueue, draw)
	for g := range draw {
		if err := termboxDraw(g); err != nil {
			panic(err)
		}
	}
}

type ProxyGame struct {
	RW   io.ReadWriteCloser
	game galaxy.Game
}

func (p *ProxyGame) Player() (player galaxy.Player, err error) {
	dec := json.NewDecoder(p.RW)
	err = dec.Decode(&player)
	return
}

func (p *ProxyGame) Run(cmdQueue <-chan galaxy.Command, draw chan<- galaxy.Game) error {
	defer close(draw)

	dec := json.NewDecoder(p.RW)

	go func() { // Send commands
		enc := json.NewEncoder(p.RW)
		for cmd := range cmdQueue {
			if err := enc.Encode(cmd); err != nil {
				return
			}
		}
	}()

	for {
		err := dec.Decode(&p.game)
		switch err {
		case nil:
			draw <- p.game
		case io.EOF:
			return nil
		default:
			return err
		}
	}
}

func (p *ProxyGame) Probe(x, y int) *galaxy.Planet {
	return p.game.Probe(x, y)
}

func termboxInput(player galaxy.Player, game galaxy.GameInterface, cmds chan galaxy.Command) {
	defer close(cmds)
	var from *galaxy.Planet
	fraction := float32(0.5)
	var createLink bool
	var destroyLink bool
	for {
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC, termbox.KeyCtrlD:
				cmds <- galaxy.Command{CommandType: galaxy.CommandQuit}
				return
			}

			// Use 0 - 9 keys to define % of ships cast from each planet
			switch ev.Ch {
			case '0':
				fraction = 1.0
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				fraction = float32(ev.Ch-'0') / 10.0
			case 'l':
				createLink = true
				destroyLink = false
			case 'q':
				createLink = false
				destroyLink = false
			case 'd':
				createLink = false
				destroyLink = true
			}

		case termbox.EventMouse:
			if ev.Key == termbox.MouseRelease {
				if from == nil {
					from = game.Probe(ev.MouseX, ev.MouseY)
					break
				}

				to := game.Probe(ev.MouseX, ev.MouseY)
				switch {
				case createLink:
					cmds <- galaxy.Command{galaxy.CommandCreateLink, from.Id, to.Id, from.Size * 0.05, player}
					createLink = false
				case destroyLink:
					cmds <- galaxy.Command{galaxy.CommandDestroyLink, from.Id, to.Id, 0, player}
					destroyLink = false
				default:
					cmds <- galaxy.Command{galaxy.CommandSendFleet, from.Id, to.Id, from.Units * fraction, player}
				}
				from = nil
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

func termboxDraw(g galaxy.Game) error {
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
