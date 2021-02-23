package main

import (
	"hacpaka/gogol"
	"math/rand"
	"runtime"
)

const (
	width = 1920
	height = 1080
)

type Unit struct {
	Life uint
	Gain uint
}

func (u *Unit) Refresh() {
	u.Life = u.Gain
}

type World struct {
	units [][]*Unit
}

func (w *World) Init (columns, rows, population uint) {
	w.units = make([][]*Unit, columns)

	for x := range w.units {
		w.units[x] = make([]*Unit, rows)

		for y :=  range w.units[x] {
			w.units[x][y] = &Unit{0, 0}
		}
	}

	for population > 0 {
		w.units[rand.Intn(int(columns))][rand.Intn(int(rows))].Gain = 1
		population--
	}
}

func main() {
	runtime.LockOSThread()

	world := new(World)
	//world.Init()

	action := func(units [][]*gogol.Point) error {
		for x := range units {
			for y := range units[x] {
				units[x][y].Refresh()

				nb := units[x][y].Neighbors(units)

				if units[x][y].Life > 0 {
					if nb < 2 || nb > 3 {
						units[x][y].Gain = 0
					}

					if nb == 2 || nb == 3 {
						units[x][y].Gain = 1
					}
				}else {
					if nb == 3 {
						units[x][y].Gain = 1
					}
				}
			}
		}

		return nil
	}

	err := new(gogol.Engine).Init(action, width, height)
	if err != nil {
		panic(err)
	}
}
