package engine

import minemath "github.com/wmattei/minceraft/math"

type RenderObject struct {
	VAO         uint32
	VBO         uint32
	EBO         uint32
	Triangles   []*Triangle
	ModelMatrix *minemath.Mat4
}
