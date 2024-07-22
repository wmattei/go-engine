package main

import (
	minemath "github.com/wmattei/minceraft/math"
)

var (
	normalRight  = minemath.Vec3{1, 0, 0}
	normalLeft   = minemath.Vec3{-1, 0, 0}
	normalTop    = minemath.Vec3{0, 1, 0}
	normalBottom = minemath.Vec3{0, -1, 0}
	normalFront  = minemath.Vec3{0, 0, 1}
	normalBack   = minemath.Vec3{0, 0, -1}
)

type Face struct {
	Visible bool
	Texture *Texture
	Normal  *minemath.Vec3
}

type Block struct {
	Type     BlockType
	Faces    [6]*Face
	Color    *Color
	Position *minemath.Vec3
}

type Direction int

const (
	Right Direction = iota
	Left
	Top
	Bottom
	Front
	Back
)

type BlockType string

const (
	Air   BlockType = "air"
	Grass BlockType = "grass"
)

func (b *Block) Update(w *Chunk) {}

func (b *Block) CullFaces(c *Chunk) {
	position := b.Position
	neighbors := [6][3]float32{
		{position.X() + 1, position.Y(), position.Z()}, // right
		{position.X() - 1, position.Y(), position.Z()}, // left
		{position.X(), position.Y() + 1, position.Z()}, // top
		{position.X(), position.Y() - 1, position.Z()}, // bottom
		{position.X(), position.Y(), position.Z() + 1}, // front
		{position.X(), position.Y(), position.Z() - 1}, // back
	}
	for i, n := range neighbors {
		neighborBlock := c.At(int(n[0]), int(n[1]), int(n[2]))
		b.Faces[i].Visible = neighborBlock != nil && neighborBlock.Type == Air
	}

	if position.X() == 0 {
		leftNeighborChunk := c.LeftNeighbor()
		if leftNeighborChunk == nil {
			b.Faces[Left].Visible = false
		} else {
			neighborBlock := leftNeighborChunk.GetBlock(int(15), int(position.Y()), int(position.Z()))
			if !neighborBlock.IsSolid() {
				b.Faces[Left].Visible = true
			}
		}
	}

	if position.X() == 15 {
		rightNeighborChunk := c.RightNeighbor()
		if rightNeighborChunk == nil {
			b.Faces[Right].Visible = false
		} else {
			neighborBlock := rightNeighborChunk.GetBlock(int(0), int(position.Y()), int(position.Z()))
			if !neighborBlock.IsSolid() {
				b.Faces[Right].Visible = true
			}
		}

	}

	if position.Z() == 0 {
		backNeighborChunk := c.BackNeighbor()
		if backNeighborChunk == nil {
			b.Faces[Back].Visible = false
		} else {
			neighborBlock := backNeighborChunk.GetBlock(int(position.X()), int(position.Y()), int(15))
			if !neighborBlock.IsSolid() {
				b.Faces[Back].Visible = true
			}
		}
	}

	if position.Z() == 15 {
		frontNeighborChunk := c.FrontNeighbor()
		if frontNeighborChunk == nil {
			b.Faces[Front].Visible = false
		} else {
			neighborBlock := frontNeighborChunk.GetBlock(int(position.X()), int(position.Y()), int(0))
			if !neighborBlock.IsSolid() {
				b.Faces[Front].Visible = true
			}
		}

	}

}

func (f *Face) GetVerticesAndIndices(x, y, z int, direction Direction, indexOffset uint32, lightDir minemath.Vec3) ([]float32, []uint32) {
	var faceVertices []float32
	var faceIndices []uint32

	var clr *Color
	alpha := 1.0

	if f.Texture.Color != nil {
		clr = f.Texture.Color
	} else {
		clr = &Color{R: 255, G: 255, B: 255}
		alpha = 0.0
	}

	intensity := minemath.CalculateLightIntensity(*f.Normal, lightDir)

	color := clr.ToVec4()

	color[0] = color[0] * intensity
	color[1] = color[1] * intensity
	color[2] = color[2] * intensity

	index := f.Texture.Index
	// if direction == Back {
	// 	index = 0
	// }

	switch direction {
	case Right:
		faceVertices = []float32{
			float32(x + 1), float32(y), float32(z), color[0], color[1], color[2], color[3], float32(alpha), 0.0, 1.0, float32(index),
			float32(x + 1), float32(y + 1), float32(z), color[0], color[1], color[2], color[3], float32(alpha), 0.0, 0.0, float32(index),
			float32(x + 1), float32(y + 1), float32(z + 1), color[0], color[1], color[2], color[3], float32(alpha), 1.0, 0.0, float32(index),
			float32(x + 1), float32(y), float32(z + 1), color[0], color[1], color[2], color[3], float32(alpha), 1.0, 1.0, float32(index),
		}
	case Left:
		faceVertices = []float32{
			float32(x), float32(y), float32(z + 1), color[0], color[1], color[2], color[3], float32(alpha), 0.0, 1.0, float32(index),
			float32(x), float32(y + 1), float32(z + 1), color[0], color[1], color[2], color[3], float32(alpha), 0.0, 0.0, float32(index),
			float32(x), float32(y + 1), float32(z), color[0], color[1], color[2], color[3], float32(alpha), 1.0, 0.0, float32(index),
			float32(x), float32(y), float32(z), color[0], color[1], color[2], color[3], float32(alpha), 1.0, 1.0, float32(index),
		}
	case Top:
		faceVertices = []float32{
			float32(x), float32(y + 1), float32(z), color[0], color[1], color[2], color[3], float32(alpha), 0.0, 1.0, float32(index),
			float32(x + 1), float32(y + 1), float32(z), color[0], color[1], color[2], color[3], float32(alpha), 0.0, 0.0, float32(index),
			float32(x + 1), float32(y + 1), float32(z + 1), color[0], color[1], color[2], color[3], float32(alpha), 1.0, 0.0, float32(index),
			float32(x), float32(y + 1), float32(z + 1), color[0], color[1], color[2], color[3], float32(alpha), 1.0, 1.0, float32(index),
		}
	case Bottom:
		faceVertices = []float32{
			float32(x), float32(y), float32(z), color[0], color[1], color[2], color[3], float32(alpha), 0.0, 1.0, float32(index),
			float32(x + 1), float32(y), float32(z), color[0], color[1], color[2], color[3], float32(alpha), 0.0, 0.0, float32(index),
			float32(x + 1), float32(y), float32(z + 1), color[0], color[1], color[2], color[3], float32(alpha), 1.0, 0.0, float32(index),
			float32(x), float32(y), float32(z + 1), color[0], color[1], color[2], color[3], float32(alpha), 1.0, 1.0, float32(index),
		}
	case Front:
		faceVertices = []float32{
			float32(x), float32(y), float32(z + 1), color[0], color[1], color[2], color[3], float32(alpha), 1.0, 1.0, float32(index),
			float32(x + 1), float32(y), float32(z + 1), color[0], color[1], color[2], color[3], float32(alpha), 0.0, 1.0, float32(index),
			float32(x + 1), float32(y + 1), float32(z + 1), color[0], color[1], color[2], color[3], float32(alpha), 0.0, 0.0, float32(index),
			float32(x), float32(y + 1), float32(z + 1), color[0], color[1], color[2], color[3], float32(alpha), 1.0, 0.0, float32(index),
		}
	case Back:
		faceVertices = []float32{
			float32(x), float32(y), float32(z), color[0], color[1], color[2], color[3], float32(alpha), 1.0, 1.0, float32(index),
			float32(x + 1), float32(y), float32(z), color[0], color[1], color[2], color[3], float32(alpha), 0.0, 1.0, float32(index),
			float32(x + 1), float32(y + 1), float32(z), color[0], color[1], color[2], color[3], float32(alpha), 0.0, 0.0, float32(index),
			float32(x), float32(y + 1), float32(z), color[0], color[1], color[2], color[3], float32(alpha), 1.0, 0.0, float32(index),
		}
	}

	faceIndices = []uint32{
		indexOffset, indexOffset + 1, indexOffset + 2,
		indexOffset, indexOffset + 2, indexOffset + 3,
	}

	return faceVertices, faceIndices
}

func (b *Block) IsSolid() bool {
	return b.Type != Air
}
