package main

import (
	"github.com/wmattei/minceraft/pkg/engine"
)

const RENDER_DISTANCE = 2

type Minceraft struct {
	chunks []*Chunk
	camera *engine.PerspectiveCamera
}

func (m *Minceraft) Update(dt float32) {
	for _, chunk := range m.chunks {
		chunk.Update(dt, m.camera)
	}
}

func (m *Minceraft) Render(scene *engine.Scene) {
	scene.SetCamera(m.camera)

	for _, chunk := range m.chunks {
		chunk.Render(scene)
	}

}

func (m *Minceraft) GetScreenSize() (int, int) {
	return 1200, 980
}

func main() {
	cam := engine.NewPerspectiveCamera(
		[3]float32{0, 6, 0},
		[3]float32{0, 1, 0},
		0,
		0,
		90,
		1200.0/980.0,
		0.1,
		200.0,
	)

	game := &Minceraft{
		camera: cam,
		chunks: []*Chunk{},
	}
	for i := -RENDER_DISTANCE; i < RENDER_DISTANCE; i++ {
		for j := -RENDER_DISTANCE; j < RENDER_DISTANCE; j++ {
			chunk := NewChunk(i, j)
			game.chunks = append(game.chunks, chunk)
		}
	}

	eng := engine.NewEngine(game)
	defer engine.Destroy()

	eng.Start()
}
