package main

import "testing"

// Benchmark test for NewChunk function
func BenchmarkNewChunk(b *testing.B) {
	world := &World{
		chunks:       make(map[[2]int]*Chunk),
		textures:     make(map[string]Texture),
		noise:        Noise{},
		light:        Light{},
		activeChunk:  [2]int{0, 0},
		renderDist:   8,
		loadedChunks: make(map[[2]int]struct{}),
	}

	for x := -world.renderDist; x < world.renderDist; x++ {
		for z := -world.renderDist; z < world.renderDist; z++ {
			c := NewChunk(world, x, z, 16)
			world.chunks[[2]int{x, z}] = c
			world.loadedChunks[[2]int{x, z}] = struct{}{}
		}
	}

	b.ResetTimer()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = NewChunk(world, i, i, 16)
	}
}
