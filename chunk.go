package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	minemath "github.com/wmattei/minceraft/math"
	"github.com/wmattei/minceraft/pkg/engine"
)

const WORLD_HEIGHT = 164
const SEA_LEVEL = 64
const MaxBlocksPerChunk = 16 * 16 * WORLD_HEIGHT

type Chunk struct {
	Blocks   [][][]*Block
	Position [2]int

	VAO      uint32
	VBO      uint32
	EBO      uint32
	Vertices []float32
	Indices  []uint32

	World *World

	SolidBlocks map[[3]int]struct{}
}

func (c *Chunk) RightNeighbor() *Chunk {
	if len(c.World.chunks) <= 1 {
		return nil
	}

	return c.World.chunks[[2]int{c.Position[0] + 1, c.Position[1]}]
}

func (c *Chunk) LeftNeighbor() *Chunk {
	if len(c.World.chunks) <= 1 {
		return nil
	}
	return c.World.chunks[[2]int{c.Position[0] - 1, c.Position[1]}]
}

func (c *Chunk) FrontNeighbor() *Chunk {
	if len(c.World.chunks) <= 1 {
		return nil
	}
	return c.World.chunks[[2]int{c.Position[0], c.Position[1] + 1}]
}

func (c *Chunk) BackNeighbor() *Chunk {
	if len(c.World.chunks) <= 1 {
		return nil
	}
	return c.World.chunks[[2]int{c.Position[0], c.Position[1] - 1}]
}

func (chunk *Chunk) GenerateMesh() {
	chunk.Vertices, chunk.Indices = chunk.generateMeshData()

}

func (chunk *Chunk) Initialize() {

	gl.GenVertexArrays(1, &chunk.VAO)
	gl.GenBuffers(1, &chunk.VBO)
	gl.GenBuffers(1, &chunk.EBO)
	chunk.GenerateMesh()
	chunk.UpdateBuffers()
}

func (chunk *Chunk) Delete() {
	gl.DeleteVertexArrays(1, &chunk.VAO)
	gl.DeleteBuffers(1, &chunk.VBO)
	gl.DeleteBuffers(1, &chunk.EBO)
}

func (chunk *Chunk) CullBlocksFaces() {
	for pos := range chunk.SolidBlocks {
		block := chunk.At(pos[0], pos[1], pos[2])
		if block.NeedsCulling {
			block.CullFaces(chunk, pos)
		}
	}
}
func (chunk *Chunk) CullRightEdgeBlocksFaces() {
	for pos := range chunk.SolidBlocks {
		if pos[0] != 15 {
			continue
		}
		block := chunk.At(pos[0], pos[1], pos[2])
		block.CullFaces(chunk, pos)
	}
}

func (chunk *Chunk) CullLeftEdgeBlocksFaces() {
	for pos := range chunk.SolidBlocks {
		if pos[0] != 0 {
			continue
		}
		block := chunk.At(pos[0], pos[1], pos[2])
		block.CullFaces(chunk, pos)
	}
}

func (chunk *Chunk) CullFrontEdgeBlocksFaces() {
	for pos := range chunk.SolidBlocks {
		if pos[2] != 15 {
			continue
		}
		block := chunk.At(pos[0], pos[1], pos[2])
		block.CullFaces(chunk, pos)
	}
}

func (chunk *Chunk) CullBackEdgeBlocksFaces() {
	for pos := range chunk.SolidBlocks {
		if pos[2] != 0 {
			continue
		}
		block := chunk.At(pos[0], pos[1], pos[2])
		block.CullFaces(chunk, pos)
	}
}

func (chunk *Chunk) GetAllNeighbors() []*Chunk {
	neighbors := []*Chunk{
		chunk.RightNeighbor(),
		chunk.LeftNeighbor(),
		chunk.FrontNeighbor(),
		chunk.BackNeighbor(),
	}

	return neighbors
}

func (chunk *Chunk) UpdateBuffers() {
	gl.BindVertexArray(chunk.VAO)

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

	lightDirection := chunk.World.light.Direction

	for x := 0; x < 16; x++ {
		for y := 0; y < WORLD_HEIGHT; y++ {
			for z := 0; z < 16; z++ {
				block := chunk.At(x, y, z)
				// if x == 15 {
				// 	continue
				// }
				if block != nil && block.Type != Air {
					for direction, face := range block.Faces {
						if face.Visible {
							faceVertices, faceIndices := face.GetVerticesAndIndices(x, y, z, Direction(direction), indexOffset, *lightDirection)
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

var airBlock = &Block{Type: Air}

func NewChunk(world *World, chunkX, chunkZ, size int) *Chunk {
	// startedAt := time.Now()
	// defer func() {
	// 	elapsed := time.Since(startedAt)
	// 	fmt.Printf("Generated chunk in %s\n", elapsed)
	// }()
	estimatedSize := size * size * WORLD_HEIGHT
	chunk := &Chunk{
		Position:    [2]int{chunkX, chunkZ},
		Blocks:      make([][][]*Block, size),
		SolidBlocks: make(map[[3]int]struct{}, estimatedSize),
		World:       world,
	}

	grassSide := world.textures["grassside"]
	grassTop := world.textures["grasstop"]

	for x := 0; x < size; x++ {
		chunk.Blocks[x] = make([][]*Block, size)
		for z := 0; z < size; z++ {
			chunk.Blocks[x][z] = make([]*Block, WORLD_HEIGHT)
			chunkXPos := x + chunkX*16
			chunkZPos := z + chunkZ*16

			height := int(world.noise.GetHeight(chunkXPos, chunkZPos))
			for y := 0; y < WORLD_HEIGHT; y++ {
				if y > (height + SEA_LEVEL) {
					chunk.Blocks[x][z][y] = airBlock
					continue
				}

				block := &Block{}
				block.Type = Grass
				block.Faces = [6]Face{
					{Texture: &grassSide, Normal: normalRight, Visible: false},
					{Texture: &grassSide, Normal: normalLeft, Visible: false},
					{Texture: &grassTop, Normal: normalTop, Visible: false},
					{Texture: &grassSide, Normal: normalBottom, Visible: false},
					{Texture: &grassSide, Normal: normalFront, Visible: false},
					{Texture: &grassSide, Normal: normalBack, Visible: false},
				}

				block.NeedsCulling = y >= SEA_LEVEL+height-3
				pos := [3]int{x, y, z}
				chunk.SolidBlocks[pos] = struct{}{}
				chunk.Blocks[x][z][y] = block

			}
		}
	}

	return chunk
}

func (c *Chunk) At(x, y, z int) *Block {
	if x < 0 || x >= 16 || y < 0 || y >= WORLD_HEIGHT || z < 0 || z >= 16 {
		return nil
	}
	return c.Blocks[x][z][y]
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
