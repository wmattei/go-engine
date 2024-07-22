package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	minemath "github.com/wmattei/minceraft/math"
	"github.com/wmattei/minceraft/pkg/engine"
)

const WORLD_HEIGHT = 164
const SEA_LEVEL = 64

type Chunk struct {
	Blocks   map[[3]int]*Block
	Position [2]int

	VAO      uint32
	VBO      uint32
	EBO      uint32
	Vertices []float32
	Indices  []uint32

	World *World
}

func (c *Chunk) GetBlock(x, y, z int) *Block {
	return c.Blocks[[3]int{x, y, z}]
}

func (c *Chunk) HasRightNeighbors() bool {
	if len(c.World.chunks) <= 1 {
		return false
	}

	return c.World.chunks[[2]int{c.Position[0] + 1, c.Position[1]}] != nil
}

func (c *Chunk) HasLeftNeighbors() bool {
	if len(c.World.chunks) <= 1 {
		return false
	}
	return c.World.chunks[[2]int{c.Position[0] - 1, c.Position[1]}] != nil
}

func (c *Chunk) HasFrontNeighbors() bool {
	if len(c.World.chunks) <= 1 {
		return false
	}
	return c.World.chunks[[2]int{c.Position[0], c.Position[1] + 1}] != nil
}

func (c *Chunk) HasBackNeighbors() bool {
	if len(c.World.chunks) <= 1 {
		return false
	}
	return c.World.chunks[[2]int{c.Position[0], c.Position[1] - 1}] != nil
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

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 11*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 11*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 11*4, gl.PtrOffset(8*4))
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(3, 1, gl.FLOAT, false, 11*4, gl.PtrOffset(10*4))
	gl.EnableVertexAttribArray(3)

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
				if block.Type != Air {
					block.CullFaces(chunk)
					for direction, face := range block.Faces {
						if face.Visible {
							faceVertices, faceIndices := face.GetVerticesAndIndices(x, y, z, Direction(direction), indexOffset)
							vertices = append(vertices, faceVertices...)
							indices = append(indices, faceIndices...)
							indexOffset += 4
						}
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

func NewChunk(world *World, chunkX, chunkZ, size int) *Chunk {
	chunk := &Chunk{
		Position: [2]int{chunkX, chunkZ},
		Blocks:   make(map[[3]int]*Block),
	}

	for x := 0; x < size; x++ {
		for z := 0; z < size; z++ {
			for y := 0; y < WORLD_HEIGHT; y++ {
				blockType := Air
				if y <= SEA_LEVEL {
					blockType = Grass
				}
				block := world.NewBlock(float32(x), float32(y), float32(z), blockType)
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
	return minemath.GetTranslationMatrix(float32(c.Position[0]*16), 0, float32(c.Position[1]*16))
	// return minemath.GetTranslationMatrix(float32(c.Position[0]*16), 0, float32(c.Position[1]*16))
}

func (c Chunk) isInFrustum(frustum *engine.Frustum, viewMatrix minemath.Mat4) bool {
	chunkSize := float32(16)
	chunkHeight := float32(WORLD_HEIGHT)

	corners := [8]minemath.Vec3{
		{float32(c.Position[0]), 0, float32(c.Position[1])},
		{float32(c.Position[0]) + chunkSize, 0, float32(c.Position[1])},
		{float32(c.Position[0]), chunkHeight, float32(c.Position[1])},
		{float32(c.Position[0]) + chunkSize, chunkHeight, float32(c.Position[1])},
		{float32(c.Position[0]), 0, float32(c.Position[1]) + chunkSize},
		{float32(c.Position[0]) + chunkSize, 0, float32(c.Position[1]) + chunkSize},
		{float32(c.Position[0]), chunkHeight, float32(c.Position[1]) + chunkSize},
		{float32(c.Position[0]) + chunkSize, chunkHeight, float32(c.Position[1]) + chunkSize},
	}

	for i := 0; i < 8; i++ {
		corners[i] = minemath.TransformVec3(viewMatrix, corners[i])
	}

	planes := frustum.GetPlanes()

	for i := 0; i < 6; i++ {
		plane := planes[i]
		outside := true
		for j := 0; j < 8; j++ {
			if plane.DistanceToPoint(corners[j]) >= 0 {
				outside = false
				break
			}
		}
		if outside {
			return false
		}
	}
	return true

}
