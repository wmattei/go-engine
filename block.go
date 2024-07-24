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
	Normal  minemath.Vec3
}

type Block struct {
	Type  BlockType
	Faces [6]Face
	Color Color

	NeedsCulling bool
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

func (b *Block) CullFaces(c *Chunk, position [3]int) {
	posX, posY, posZ := position[0], position[1], position[2]

	neighbors := [6][3]int{
		{posX + 1, posY, posZ}, // right
		{posX - 1, posY, posZ}, // left
		{posX, posY + 1, posZ}, // top
		{posX, posY - 1, posZ}, // bottom
		{posX, posY, posZ + 1}, // front
		{posX, posY, posZ - 1}, // back
	}

	for i, n := range neighbors {
		neighborBlock := c.At(int(n[0]), int(n[1]), int(n[2]))
		b.Faces[i].Visible = neighborBlock != nil && neighborBlock.Type == Air
	}

	if posX == 0 {
		b.checkNeighborChunkFace(c.LeftNeighbor(), 15, posY, posZ, Left)
	}
	if posX == 15 {
		b.checkNeighborChunkFace(c.RightNeighbor(), 0, posY, posZ, Right)
	}
	if posZ == 0 {
		b.checkNeighborChunkFace(c.BackNeighbor(), posX, posY, 15, Back)
	}
	if posZ == 15 {
		b.checkNeighborChunkFace(c.FrontNeighbor(), posX, posY, 0, Front)
	}
}

func (b *Block) checkNeighborChunkFace(neighborChunk *Chunk, x, y, z int, face Direction) {
	if neighborChunk == nil {
		b.Faces[face].Visible = false
	} else {
		neighborBlock := neighborChunk.At(x, y, z)
		b.Faces[face].Visible = !neighborBlock.IsSolid()
	}
}

func (f *Face) GetVerticesAndIndices(x, y, z int, direction Direction, indexOffset uint32, lightDir minemath.Vec3) ([]float32, []uint32) {
	var faceVertices []float32
	var faceIndices []uint32

	// fmt.Println(x, y, z, string(direction))

	var clr *Color
	alpha := 1.0

	if f.Texture.Color != nil {
		clr = f.Texture.Color
	} else {
		clr = &Color{R: 255, G: 255, B: 255}
		alpha = 0.0
	}

	color := clr.ToVec4()

	intensity := minemath.CalculateLightIntensity(f.Normal, lightDir)
	color[0] = color[0] * intensity
	color[1] = color[1] * intensity
	color[2] = color[2] * intensity

	index := f.Texture.Index
	// if x == 0 {
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
