package main

import (
	"errors"
	"math/rand"
)

type t_color struct {
	R int
	G int
	B int
}

type t_world struct {
	data [][] int
}

type t_unit struct {
	data []float32
	color t_color

	life uint
	gain uint

	x int
	y int
}

func Color (index int, value int) t_color {
	if index > 2 {
		panic(errors.New("unInvalidColorIndex"))
	}

	tpl := [3]int {0, 0, 0}
	tpl[index] = value

	return t_color{tpl[0], tpl[1], tpl[2]}
}

func Triangles (columns, rows, x, y int) []float32 {
	var (
		triangles = []float32{
			-1 , 1 * (2 / float32(rows) - 1), 0,
			1 * (2 / float32(columns) - 1), -1, 0,
			-1, -1, 0,

			-1 , 1 * (2 / float32(rows) - 1), 0,
			1 * (2 / float32(columns) - 1), 1 * (2 / float32(rows) - 1), 0,
			1 * (2 / float32(columns) - 1), -1, 0,
		}
	)

	offset := []float32 {
		0.0,

		2 / float32(columns) * float32(x),
		2 / float32(rows) * float32(y),
	}

	for i := 0; i < len(triangles); i++ {
		triangles[i] += offset[(i + 1) % 3]
	}

	return triangles
}

func Units (columns, rows, population int) [][]*t_unit {
	units := make([][]*t_unit, columns)

	for x := range units {
		units[x] = make([]*t_unit, columns)

		for y :=  range units[x] {
			units[x][y] = &t_unit{Triangles(columns, rows, x, y),
				t_color{255, 255, 255 }, 0, 0, x, y}
		}
	}

	for population > 0 {
		units[rand.Intn(columns)][rand.Intn(rows)].gain = 1
		population--
	}

	return units
}

func (self *t_unit) Neighbors(units [][]*t_unit) int {
	count := 0

	for x := self.x - 1; x <= self.x + 1; x++ {
		if x >= 0 && x <= len(units) - 1 {

			for y := self.y - 1; y <= self.y + 1; y++ {
				if y >= 0 && y <= len(units[x]) -1 {

					if x != self.x || y != self.y {
						if units[x][y].life > 0 {
							count++
						}
					}
				}
			}
		}
	}

	return  count
}

func (self *t_unit) Refresh() {
	self.life = self.gain
}