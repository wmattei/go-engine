package main

import (
	"fmt"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	minemath "github.com/wmattei/minceraft/math"
	"github.com/wmattei/minceraft/pkg/engine"
)

type World struct {
	chunks   map[[2]int]*Chunk
	textures map[string]Texture
	noise    *Noise
}

func (w *World) Update() {
	for _, chunk := range w.chunks {
		chunk.Update()
	}
}

func (w *World) LoadTextures() {
	texturesFile := LoadTextures()
	w.textures = texturesFile
}

func (w *World) BindTextures(program uint32) {
	for _, texture := range w.textures {
		gl.ActiveTexture(gl.TEXTURE0 + uint32(texture.Index))
		gl.BindTexture(gl.TEXTURE_2D, texture.ref)
		uniformName := texture.UniformName
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str(uniformName)), int32(texture.Index))
	}
}

func (w *World) Render(program uint32, frustum *engine.Frustum) {
	modelLoc := gl.GetUniformLocation(program, gl.Str("model\x00"))

	for _, chunk := range w.chunks {
		model := chunk.GetModelMatrix()
		if !chunk.isInFrustum(frustum, model) {
			continue
		}

		w.BindTextures(program)

		flattenModel := model.Flatten()

		gl.UniformMatrix4fv(modelLoc, 1, false, &flattenModel[0])

		chunk.Render()
	}
}

func NewSingleBlockWorld() *World {
	world := NewWorld(0)

	chunk := NewChunk(world, 0, 0, 1)
	chunk.World = world
	chunk.Initialize()

	world.chunks[[2]int{0, 0}] = chunk

	return world
}

func NewSingleChunkWorld() *World {
	world := NewWorld(0)

	chunk := NewChunk(world, 0, 0, 16)
	chunk.World = world
	chunk.Initialize()
	world.chunks[[2]int{0, 0}] = chunk

	return world

}

func NewWorld(size int) *World {
	startedAt := time.Now()
	defer func() {
		elapsed := time.Since(startedAt)
		fmt.Println(elapsed.Milliseconds())
	}()
	world := &World{
		chunks:   make(map[[2]int]*Chunk),
		textures: map[string]Texture{},
		noise:    &Noise{},
	}

	world.LoadTextures()

	for x := -size; x < size; x++ {
		for z := -size; z < size; z++ {
			c := NewChunk(world, x, z, 16)
			world.chunks[[2]int{x, z}] = c
		}
	}

	for _, chunk := range world.chunks {
		chunk.Initialize()
	}

	return world

}

func (w *World) NewBlock(x, y, z float32, blockType BlockType) *Block {
	top := w.textures[string(blockType)+"top"]
	side := w.textures[string(blockType)+"side"]

	return &Block{
		Type:     blockType,
		Position: &minemath.Vec3{x, y, z},
		Faces: [6]*Face{
			{Texture: &side, Visible: true},
			{Texture: &side, Visible: true},
			{Texture: &top, Visible: true},
			{Texture: &side, Visible: true},
			{Texture: &side, Visible: true},
			{Texture: &side, Visible: true},
		},
	}
}
