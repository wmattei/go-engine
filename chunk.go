package main

import (
	"math/rand/v2"

	"github.com/go-gl/gl/v4.1-core/gl"
	minemath "github.com/wmattei/minceraft/math"
)

const WORLD_HEIGHT = 40

type Chunk struct {
	Blocks   map[[3]int]*Block
	Position [2]int

	VAO      uint32
	VBO      uint32
	EBO      uint32
	Vertices []float32
	Indices  []uint32
}

func (c *Chunk) Update() {
	for _, block := range c.Blocks {
		block.Update(c)
	}
}

func (chunk *Chunk) Initialize() {
	gl.GenVertexArrays(1, &chunk.VAO)
	gl.GenBuffers(1, &chunk.VBO)
	gl.GenBuffers(1, &chunk.EBO)
	chunk.UpdateBuffers()
}

func (chunk *Chunk) UpdateBuffers() {
	gl.BindVertexArray(chunk.VAO)

	vertices, indices := chunk.generateMeshData()
	chunk.Vertices = vertices
	chunk.Indices = indices

	gl.BindBuffer(gl.ARRAY_BUFFER, chunk.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(chunk.Vertices)*4, gl.Ptr(chunk.Vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, chunk.EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(chunk.Indices)*4, gl.Ptr(chunk.Indices), gl.STATIC_DRAW)

	// Set up vertex attributes
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 7*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 7*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(0)

}

func (chunk *Chunk) generateMeshData() ([]float32, []uint32) {
	var vertices []float32
	var indices []uint32
	indexOffset := uint32(0)

	for x := 0; x < 16; x++ {
		for y := 0; y < WORLD_HEIGHT; y++ {
			for z := 0; z < 16; z++ {
				block := chunk.At(x, y, z)
				if block != nil {
					block.CullFaces(chunk)
					for direction, face := range block.Faces {
						if !face.Visible {
							continue
						}
						faceVertices, faceIndices := face.GetVerticesAndIndices(x, y, z, direction, indexOffset)
						vertices = append(vertices, faceVertices...)
						indices = append(indices, faceIndices...)
						indexOffset += 4
					}
				}
			}
		}
	}
	return vertices, indices
}

func (chunk *Chunk) Render() {
	gl.BindVertexArray(chunk.VAO)
	gl.DrawElements(gl.TRIANGLES, int32(len(chunk.Indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.BindVertexArray(0)
}

func NewChunk(chunkX, chunkZ, size int) *Chunk {
	chunk := &Chunk{
		Position: [2]int{chunkX, chunkZ},
		Blocks:   make(map[[3]int]*Block),
	}
	for x := 0; x < size; x++ {
		for z := 0; z < size; z++ {
			for y := 0; y < WORLD_HEIGHT; y++ {
				color := randomColor()
				block := NewBlock(float32(x), float32(y), float32(z), &color)
				chunk.Blocks[[3]int{x, y, z}] = block
			}
		}
	}

	return chunk
}

func (c *Chunk) At(x, y, z int) *Block {
	return c.Blocks[[3]int{x, y, z}]
}

func (c *Chunk) GetModelMatrix() minemath.Mat4 {
	return minemath.GetTranslationMatrix(c.Position[0]*16, 0, c.Position[1]*16)
}

func randomColor() Color {
	return Color{
		R: rand.IntN(255),
		G: rand.IntN(255),
		B: rand.IntN(255),
	}
}
