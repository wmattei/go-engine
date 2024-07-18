package main

import minemath "github.com/wmattei/minceraft/math"

var (
	BLACK       = Color{0, 0, 0}
	WHITE       = Color{255, 255, 255}
	RED         = Color{255, 0, 0}
	GREEN       = Color{0, 255, 0}
	BLUE        = Color{0, 0, 255}
	LIGHT_BLUE  = Color{173, 216, 230}
	LIGHT_GREEN = Color{144, 238, 144}
)

type Color struct {
	R, G, B int
}

func (c Color) ToVec4() minemath.Vec4 {
	return minemath.Vec4{float32(c.R) / 255, float32(c.G) / 255, float32(c.B) / 255, 1.0}
}
