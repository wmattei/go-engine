package main

import (
	"math"
	"math/rand/v2"
)

type Noise struct{}

const permutationSize = 256

var permutation = [permutationSize * 2]int{}

func init() {
	// Initialize the permutation array
	p := rand.Perm(permutationSize)
	for i := 0; i < permutationSize; i++ {
		permutation[i] = p[i]
		permutation[i+permutationSize] = p[i]
	}
}

func fade(t float64) float64 {
	return t * t * t * (t*(t*6-15) + 10)
}

// Linear interpolation function
func lerp(t, a, b float64) float64 {
	return a + t*(b-a)
}

func grad(hash int, x, y float64) float64 {
	h := hash & 15
	var u, v float64
	if h < 8 {
		u = x
	} else {
		u = y
	}
	if h < 4 {
		v = y
	} else if h == 12 || h == 14 {
		v = x
	} else {
		v = 0
	}

	if h&1 == 0 {
		u = -u
	}
	if h&2 == 0 {
		v = -v
	}

	return u + v
}

func PerlinNoise2D(x, y float64) float64 {
	x0 := int(math.Floor(x)) & 255
	y0 := int(math.Floor(y)) & 255
	x -= math.Floor(x)
	y -= math.Floor(y)
	u := fade(x)
	v := fade(y)

	a := permutation[x0] + y0
	aa := permutation[a]
	ab := permutation[a+1]
	b := permutation[x0+1] + y0
	ba := permutation[b]
	bb := permutation[b+1]

	// Blend the results from the four corners
	return lerp(v, lerp(u, grad(permutation[aa], x, y), grad(permutation[ba], x-1, y)),
		lerp(u, grad(permutation[ab], x, y-1), grad(permutation[bb], x-1, y-1)))
}

func (n *Noise) GetHeight(x, z int) int {
	frequency := 0.02
	amplitude := 30.0

	// Generate Perlin noise value for the given coordinates
	noiseValue := PerlinNoise2D(float64(x)*frequency, float64(z)*frequency)

	// Scale the noise value to get the height
	y := amplitude * noiseValue

	return int(y)
}
