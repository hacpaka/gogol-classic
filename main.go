package main

import (
	"log"
	"runtime"
)

const (
	duration = 10
	size = 20
	width  =  1920
	height = 1080
)

func main() {
	runtime.LockOSThread()

	columns := int(width / size)
	rows := int(height / size)

	units := Units(columns, rows, 1500)
	log.Printf("World created: %vx%v", columns, rows)

	action := func(prog uint32) error {
		for x := range units {
			for y := range units[x] {
				units[x][y].Refresh()

				nb := units[x][y].Neighbors(units)

				if units[x][y].life > 0 {
					if nb < 2 || nb > 3 {
						units[x][y].gain = 0
					}

					if nb == 2 || nb == 3 {
						units[x][y].gain = 1
					}
				}else {
					if nb == 3 {
						units[x][y].gain = 1
					}
				}
			}
		}

		for x := range units {
			for y := range units[x] {
				if units[x][y].life > 0 {
					draw(prog, units[x][y].data, units[x][y].color)
				}
			}
		}

		return nil
	}

	err := run(action, duration, width, height)
	if err != nil {
		panic(err)
	}
}
