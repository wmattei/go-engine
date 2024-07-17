package main

import (
	minemath "github.com/wmattei/minceraft/math"
	"github.com/wmattei/minceraft/pkg/engine"
)

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
	return 1000, 1000
}

func main() {
	cam := engine.NewPerspectiveCamera(
		[3]float32{0, 6, 0},
		[3]float32{0, 1, 0},
		0,
		0,
		minemath.DegreesToRadians(90),
		1,
		0.1,
		100.0,
	)

	chunk := NewSampleChunk()
	game := &Minceraft{
		camera: cam,
		chunks: []*Chunk{chunk},
	}

	eng := engine.NewEngine(game)
	defer engine.Destroy()

	eng.Start()
}
