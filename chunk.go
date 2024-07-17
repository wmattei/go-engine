package main

import (
	"github.com/wmattei/minceraft/pkg/engine"
)

const WORLD_HEIGHT = 1

type Chunk struct {
	blocks   []*Block
	position [2]int
}

func (c *Chunk) Update(dt float32, camera *engine.PerspectiveCamera) {
	for _, block := range c.blocks {
		block.NeedsRender = true
	}
}

func (c *Chunk) Render(scene *engine.Scene) {
	for _, block := range c.blocks {
		if block.NeedsRender {
			block.Render(scene)
		}
	}
}

func NewChunk(x, z int) *Chunk {
	chunk := &Chunk{
		position: [2]int{x, z},
	}
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			for k := -WORLD_HEIGHT; k < WORLD_HEIGHT; k++ {
				x := float32(i) + float32(x)*16
				y := float32(k)
				z := float32(j) + float32(z)*16

				block := NewBlock("dirt", [3]float32{x, y, z})
				chunk.blocks = append(chunk.blocks, block)
			}
		}
	}

	// chunk.blocks = append(chunk.blocks, NewBlock("dirt", [3]float32{0, 0, -1}))
	return chunk
}
