package main

import (
	"sync"

	"github.com/go-gl/gl/v4.1-core/gl"
	minemath "github.com/wmattei/minceraft/math"
	"github.com/wmattei/minceraft/pkg/engine"
)

type World struct {
	chunks      map[[2]int]*Chunk
	textures    map[string]Texture
	noise       Noise
	light       Light
	activeChunk [2]int
	renderDist  int

	loadedChunks map[[2]int]struct{}
}

func (w *World) Update(camera *engine.PerspectiveCamera) {
	cameraChunkX := int(camera.Position[0] / 16)
	cameraChunkZ := int(camera.Position[2] / 16)

	if cameraChunkX != w.activeChunk[0] || cameraChunkZ != w.activeChunk[1] {
		w.activeChunk = [2]int{cameraChunkX, cameraChunkZ}
		w.LoadChunks()
	}

}

func (w *World) LoadChunks() {
	newLoadedChunks := make(map[[2]int]struct{})
	lastAddedChunks := make(map[[2]int]bool)
	activeX, activeZ := w.activeChunk[0], w.activeChunk[1]
	renderDist := w.renderDist

	var wg sync.WaitGroup
	mu := &sync.Mutex{}

	for x := activeX - renderDist; x < activeX+renderDist; x++ {
		for z := activeZ - renderDist; z < activeZ+renderDist; z++ {
			wg.Add(1)
			go func(x, z int) {
				defer wg.Done()
				chunkPos := [2]int{x, z}
				mu.Lock()
				newLoadedChunks[chunkPos] = struct{}{}
				mu.Unlock()

				if _, ok := w.chunks[chunkPos]; !ok {
					chunk := NewChunk(w, x, z, 16)
					chunk.World = w

					mu.Lock()
					w.chunks[chunkPos] = chunk
					lastAddedChunks[chunkPos] = true
					mu.Unlock()
				}
			}(x, z)
		}
	}

	wg.Wait()

	// Unload chunks that are no longer within the render distance
	for pos := range w.loadedChunks {
		if _, exists := newLoadedChunks[pos]; !exists {
			w.chunks[pos].Delete()
			delete(w.chunks, pos)
			delete(w.loadedChunks, pos)
		}
	}

	for pos := range lastAddedChunks {
		chunk := w.chunks[pos]
		chunk.CullBlocksFaces()
		chunk.Initialize()
		w.updateNeighborChunks(chunk, lastAddedChunks)
	}

	// Update loaded chunks
	w.loadedChunks = newLoadedChunks
}

func (w *World) updateNeighborChunks(chunk *Chunk, lastAddedChunks map[[2]int]bool) {
	neighbors := chunk.GetAllNeighbors()

	neighborFuncs := []func(*Chunk){
		func(c *Chunk) { c.CullLeftEdgeBlocksFaces() },
		func(c *Chunk) { c.CullRightEdgeBlocksFaces() },
		func(c *Chunk) { c.CullBackEdgeBlocksFaces() },
		func(c *Chunk) { c.CullFrontEdgeBlocksFaces() },
	}

	for i, n := range neighbors {
		if n == nil {
			continue
		}
		nPos := n.Position
		if neighborChunk, ok := w.chunks[[2]int{nPos[0], nPos[1]}]; ok && !lastAddedChunks[nPos] {
			neighborFuncs[i](neighborChunk)
			neighborChunk.GenerateMesh()
			neighborChunk.UpdateBuffers()
		}
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
		// if !chunk.isInFrustum(frustum, model) {
		// 	continue
		// }

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
	world := &World{
		chunks:       make(map[[2]int]*Chunk, size*size),
		loadedChunks: make(map[[2]int]struct{}, size*size),
		textures:     map[string]Texture{},
		noise:        Noise{},
		light:        Light{Direction: &minemath.Vec3{0.9, 1, 0.5}},
		renderDist:   size,
	}

	world.LoadTextures()

	for x := -size; x < size; x++ {
		for z := -size; z < size; z++ {
			c := NewChunk(world, x, z, 16)
			world.chunks[[2]int{x, z}] = c
			world.loadedChunks[[2]int{x, z}] = struct{}{}
		}
	}

	for _, chunk := range world.chunks {
		chunk.CullBlocksFaces()
		chunk.Initialize()
	}

	return world

}
