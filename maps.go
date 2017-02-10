package main

import (
	"math/rand"
)

func NewDefaultGame() *Game {
	return &Game{
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
	}
}

func NewSquareGameMap() *Game {
	return &Game{
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
				Center: Vec2{5, 20},
				Player: 0,
				Units:  100,
				Size:   1,
			},
			"d": &Planet{
				Center: Vec2{20, 20},
				Player: 2,
				Units:  100,
				Size:   2,
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
			Center: Vec2{x, y},
			Player: 0,
			Units:  rand.Float32() * 100,
			Size:   rand.Float32() * 2,
		}
	}

	x := rand.Float32() * width
	y := rand.Float32() * height
	g.Planets["a"] = &Planet{
		Center: Vec2{x, y},
		Player: 1,
		Units:  100,
		Size:   2,
	}

	x = rand.Float32() * width
	y = rand.Float32() * height
	g.Planets["b"] = &Planet{
		Center: Vec2{x, y},
		Player: 2,
		Units:  100,
		Size:   2,
	}

	return g
}
