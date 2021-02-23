package main

import (
	"math/rand"
)

type TColor struct {
	R int
	G int
	B int
}

type TUnit struct {
	data  []float32
	color TColor

	Life uint
	Gain uint

	x int
	y int
}

//func Color (index int, value int) TColor {
//	if index > 2 {
//		panic(errors.New("unInvalidColorIndex"))
//	}
//
//	tpl := [3]int {0, 0, 0}
//	tpl[index] = value
//
//	return TColor{tpl[0], tpl[1], tpl[2]}
//}

//func Triangles (columns, rows uint, x, y int) []float32 {
//	var (
//		triangles = []float32{
//			-1 , 1 * (2 / float32(rows) - 1), 0,
//			1 * (2 / float32(columns) - 1), -1, 0,
//			-1, -1, 0,
//
//			-1 , 1 * (2 / float32(rows) - 1), 0,
//			1 * (2 / float32(columns) - 1), 1 * (2 / float32(rows) - 1), 0,
//			1 * (2 / float32(columns) - 1), -1, 0,
//		}
//	)
//
//	offset := []float32 {
//		0.0,
//
//		2 / float32(columns) * float32(x),
//		2 / float32(rows) * float32(y),
//	}
//
//	for i := 0; i < len(triangles); i++ {
//		triangles[i] += offset[(i + 1) % 3]
//	}
//
//	return triangles
//}

func Units (columns, rows, population uint) [][]*TUnit {
	units := make([][]*TUnit, columns)

	for x := range units {
		units[x] = make([]*TUnit, columns)

		for y :=  range units[x] {
			units[x][y] = &TUnit{Triangles(columns, rows, x, y),
				TColor{255, 255, 255 }, 0, 0, x, y}
		}
	}

	for population > 0 {
		units[rand.Intn(int(columns))][rand.Intn(int(rows))].Gain = 1
		population--
	}

	return units
}

func (u *TUnit) Neighbors(units [][]*TUnit) int {
	count := 0

	for x := u.x - 1; x <= u.x + 1; x++ {
		if x >= 0 && x <= len(units) - 1 {

			for y := u.y - 1; y <= u.y + 1; y++ {
				if y >= 0 && y <= len(units[x]) -1 {

					if x != u.x || y != u.y {
						if units[x][y].Life > 0 {
							count++
						}
					}
				}
			}
		}
	}

	return  count
}

//func (u *TUnit) Refresh() {
//	u.Life = u.Gain
//}