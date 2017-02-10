package main

import (
	"math/rand"
)

func NewSquareGameMap() *Game {
	return &Game{
		Planets: map[string]*Planet{
			"a": &Planet{
				Center:   Vec2{5, 5},
				Player:   1,
				Units:    100,
				Size:     2,
				Capacity: 500,
			},
			"b": &Planet{
				Center:   Vec2{20, 5},
				Player:   0,
				Units:    100,
				Size:     1,
				Capacity: 500,
			},
			"c": &Planet{
				Center:   Vec2{5, 20},
				Player:   0,
				Units:    100,
				Size:     1,
				Capacity: 500,
			},
			"d": &Planet{
				Center:   Vec2{20, 20},
				Player:   2,
				Units:    100,
				Size:     2,
				Capacity: 500,
			},
		},
	}
}

func NewRandomMap(width, height float32) *Game {
	g := &Game{Planets: map[string]*Planet{}}

	letters := "cdefghijklmno"

	for _, letter := range letters {
		x := rand.Float32() * width
		y := rand.Float32() * height
		g.Planets[string(letter)] = &Planet{
			Center:   Vec2{x, y},
			Player:   0,
			Units:    rand.Float32() * 100,
			Size:     rand.Float32() * 2,
			Capacity: 500 + 500*rand.Float32(),
		}
	}

	x := rand.Float32() * width
	y := rand.Float32() * height
	g.Planets["a"] = &Planet{
		Center:   Vec2{x, y},
		Player:   1,
		Units:    100,
		Size:     2,
		Capacity: 500 + 500*rand.Float32(),
	}

	x = rand.Float32() * width
	y = rand.Float32() * height
	g.Planets["b"] = &Planet{
		Center:   Vec2{x, y},
		Player:   2,
		Units:    100,
		Size:     2,
		Capacity: 500 + 500*rand.Float32(),
	}

	return g
}
