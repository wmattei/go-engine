package main

import "math"

type Noise struct{}

func (n *Noise) GetHeight(x, z int) int {
	frequency := float64(0.1)
	amplitude := float64(5.0)
	y := amplitude * math.Sin(frequency*float64(x)) * math.Sin(frequency*float64(z))

	return int(y)
}
