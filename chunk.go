package main

import (
	"github.com/wmattei/minceraft/pkg/engine"
)

type Chunk struct {
	blocks []*Block
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

func NewSampleChunk() *Chunk {
	chunk := &Chunk{}
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			block := NewBlock("dirt", [3]float32{float32(i), 0, float32(j)})
			chunk.blocks = append(chunk.blocks, block)
		}
	}

	// chunk.blocks = append(chunk.blocks, NewBlock("dirt", [3]float32{0, 0, -1}))
	return chunk
}
