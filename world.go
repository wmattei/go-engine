package main

type World struct {
	chunks map[[2]int]*Chunk
}

func (w *World) Update() {
	for _, chunk := range w.chunks {
		chunk.Update()
	}
}

func NewSingleBlockWorld() *World {
	chunk := NewChunk(0, 0, 1)
	chunk.Initialize()
	world := &World{
		chunks: map[[2]int]*Chunk{
			{0, 0}: chunk,
		},
	}

	return world
}

func NewSingleChunkWorld() *World {
	world := &World{
		chunks: make(map[[2]int]*Chunk),
	}

	chunk := NewChunk(0, 0, 16)
	chunk.Initialize()
	world.chunks[[2]int{0, 0}] = chunk

	return world

}

func NewWorld(size int) *World {
	world := &World{
		chunks: make(map[[2]int]*Chunk),
	}

	for x := -size; x < size; x++ {
		for z := -size; z < size; z++ {
			chunk := NewChunk(x, z, 16)
			chunk.World = world
			world.chunks[[2]int{x, z}] = chunk
		}
	}

	for _, chunk := range world.chunks {
		chunk.Initialize()
	}

	return world

}
