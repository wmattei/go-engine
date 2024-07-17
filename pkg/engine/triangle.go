package engine

import minemath "github.com/wmattei/minceraft/math"

type Vertex struct {
	Position *minemath.Vec3
	Color    Color
}

type Triangle struct {
	Vertices [3]Vertex
}
