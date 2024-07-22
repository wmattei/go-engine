package main

import (
	minemath "github.com/wmattei/minceraft/math"
)

type Face struct {
	Visible bool
	Texture *Texture
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
		b.Faces[i].Visible = c.At(int(n[0]), int(n[1]), int(n[2])) == nil
	}

	if position.X() == 15 && c.HasRightNeighbors() {
		b.Faces[Right].Visible = false
	}

	if position.X() == 0 && c.HasLeftNeighbors() {
		b.Faces[Left].Visible = false
	}

	if position.Z() == 15 && c.HasFrontNeighbors() {
		b.Faces[Front].Visible = false
	}

	if position.Z() == 0 && c.HasBackNeighbors() {
		b.Faces[Back].Visible = false
	}
}

func (f *Face) GetVerticesAndIndices(x, y, z int, direction Direction, indexOffset uint32) ([]float32, []uint32) {
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

	color := clr.ToVec4()

	switch direction {
	case Right:
		faceVertices = []float32{
			float32(x + 1), float32(y), float32(z), color[0], color[1], color[2], color[3], float32(alpha), 1.0, 0.0, float32(f.Texture.Index), // Bottom-right
			float32(x + 1), float32(y + 1), float32(z), color[0], color[1], color[2], color[3], float32(alpha), 1.0, 1.0, float32(f.Texture.Index), // Top-right
			float32(x + 1), float32(y + 1), float32(z + 1), color[0], color[1], color[2], color[3], float32(alpha), 0.0, 1.0, float32(f.Texture.Index), // Top-left
			float32(x + 1), float32(y), float32(z + 1), color[0], color[1], color[2], color[3], float32(alpha), 0.0, 0.0, float32(f.Texture.Index), // Bottom-left
		}
	case Left:
		faceVertices = []float32{
			float32(x), float32(y), float32(z + 1), color[0], color[1], color[2], color[3], float32(alpha), 1.0, 0.0, float32(f.Texture.Index), // Bottom-right
			float32(x), float32(y + 1), float32(z + 1), color[0], color[1], color[2], color[3], float32(alpha), 1.0, 1.0, float32(f.Texture.Index), // Top-right
			float32(x), float32(y + 1), float32(z), color[0], color[1], color[2], color[3], float32(alpha), 0.0, 1.0, float32(f.Texture.Index), // Top-left
			float32(x), float32(y), float32(z), color[0], color[1], color[2], color[3], float32(alpha), 0.0, 0.0, float32(f.Texture.Index), // Bottom-left
		}
	case Top:
		faceVertices = []float32{
			float32(x), float32(y + 1), float32(z), color[0], color[1], color[2], color[3], float32(alpha), 0.0, 1.0, float32(f.Texture.Index), // Bottom-left
			float32(x + 1), float32(y + 1), float32(z), color[0], color[1], color[2], color[3], float32(alpha), 1.0, 1.0, float32(f.Texture.Index), // Bottom-right
			float32(x + 1), float32(y + 1), float32(z + 1), color[0], color[1], color[2], color[3], float32(alpha), 1.0, 0.0, float32(f.Texture.Index), // Top-right
			float32(x), float32(y + 1), float32(z + 1), color[0], color[1], color[2], color[3], float32(alpha), 0.0, 0.0, float32(f.Texture.Index), // Top-left
		}
	case Bottom:
		faceVertices = []float32{
			float32(x), float32(y), float32(z), color[0], color[1], color[2], color[3], float32(alpha), 0.0, 0.0, float32(f.Texture.Index), // Bottom-left
			float32(x + 1), float32(y), float32(z), color[0], color[1], color[2], color[3], float32(alpha), 1.0, 0.0, float32(f.Texture.Index), // Bottom-right
			float32(x + 1), float32(y), float32(z + 1), color[0], color[1], color[2], color[3], float32(alpha), 1.0, 1.0, float32(f.Texture.Index), // Top-right
			float32(x), float32(y), float32(z + 1), color[0], color[1], color[2], color[3], float32(alpha), 0.0, 1.0, float32(f.Texture.Index), // Top-left
		}
	case Front:
		faceVertices = []float32{
			float32(x), float32(y), float32(z + 1), color[0], color[1], color[2], color[3], float32(alpha), 0.0, 0.0, float32(f.Texture.Index), // Bottom-left
			float32(x + 1), float32(y), float32(z + 1), color[0], color[1], color[2], color[3], float32(alpha), 1.0, 0.0, float32(f.Texture.Index), // Bottom-right
			float32(x + 1), float32(y + 1), float32(z + 1), color[0], color[1], color[2], color[3], float32(alpha), 1.0, 1.0, float32(f.Texture.Index), // Top-right
			float32(x), float32(y + 1), float32(z + 1), color[0], color[1], color[2], color[3], float32(alpha), 0.0, 1.0, float32(f.Texture.Index), // Top-left
		}
	case Back:
		faceVertices = []float32{
			float32(x), float32(y), float32(z), color[0], color[1], color[2], color[3], float32(alpha), 0.0, 0.0, float32(f.Texture.Index), // Bottom-left
			float32(x + 1), float32(y), float32(z), color[0], color[1], color[2], color[3], float32(alpha), 1.0, 0.0, float32(f.Texture.Index), // Bottom-right
			float32(x + 1), float32(y + 1), float32(z), color[0], color[1], color[2], color[3], float32(alpha), 1.0, 1.0, float32(f.Texture.Index), // Top-right
			float32(x), float32(y + 1), float32(z), color[0], color[1], color[2], color[3], float32(alpha), 0.0, 1.0, float32(f.Texture.Index), // Top-left
		}
	}

	faceIndices = []uint32{
		indexOffset, indexOffset + 1, indexOffset + 2,
		indexOffset, indexOffset + 2, indexOffset + 3,
	}

	return faceVertices, faceIndices
}
