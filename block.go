package main

import (
	minemath "github.com/wmattei/minceraft/math"
	"github.com/wmattei/minceraft/pkg/engine"
)

type Block struct {
	*engine.Model
	Name string
}

func (b *Block) Render(scene *engine.Scene) {
	scene.AddModel(b.Model)
}

func NewBlock(name string, at minemath.Vec3) *Block {
	return &Block{
		Name: name,
		Model: engine.NewModel(
			[]*engine.Mesh{
				{
					Triangles: []*engine.Triangle{
						// FRONT
						{
							Vertices: [3]engine.Vertex{
								{Position: &minemath.Vec3{-.5, -.5, .5}, Color: engine.WHITE},
								{Position: &minemath.Vec3{.5, -.5, .5}, Color: engine.WHITE},
								{Position: &minemath.Vec3{.5, .5, .5}, Color: engine.WHITE},
							},
						},
						{
							Vertices: [3]engine.Vertex{
								{Position: &minemath.Vec3{-.5, -.5, .5}, Color: engine.WHITE},
								{Position: &minemath.Vec3{.5, .5, .5}, Color: engine.WHITE},
								{Position: &minemath.Vec3{-.5, .5, .5}, Color: engine.WHITE},
							},
						},
						// BACK
						{
							Vertices: [3]engine.Vertex{
								{Position: &minemath.Vec3{-.5, -.5, -.5}, Color: engine.WHITE},
								{Position: &minemath.Vec3{.5, -.5, -.5}, Color: engine.WHITE},
								{Position: &minemath.Vec3{.5, .5, -.5}, Color: engine.WHITE},
							},
						},
						{
							Vertices: [3]engine.Vertex{
								{Position: &minemath.Vec3{-.5, -.5, -.5}, Color: engine.WHITE},
								{Position: &minemath.Vec3{.5, .5, -.5}, Color: engine.WHITE},
								{Position: &minemath.Vec3{-.5, .5, -.5}, Color: engine.WHITE},
							},
						},
						// LEFT
						{
							Vertices: [3]engine.Vertex{
								{Position: &minemath.Vec3{-.5, -.5, .5}, Color: engine.WHITE},
								{Position: &minemath.Vec3{-.5, -.5, -.5}, Color: engine.WHITE},
								{Position: &minemath.Vec3{-.5, .5, -.5}, Color: engine.WHITE},
							},
						},
						{
							Vertices: [3]engine.Vertex{
								{Position: &minemath.Vec3{-.5, -.5, .5}, Color: engine.WHITE},
								{Position: &minemath.Vec3{-.5, .5, -.5}, Color: engine.WHITE},
								{Position: &minemath.Vec3{-.5, .5, .5}, Color: engine.WHITE},
							},
						},
						// RIGHT
						{
							Vertices: [3]engine.Vertex{
								{Position: &minemath.Vec3{.5, -.5, .5}, Color: engine.WHITE},
								{Position: &minemath.Vec3{.5, -.5, -.5}, Color: engine.WHITE},
								{Position: &minemath.Vec3{.5, .5, -.5}, Color: engine.WHITE},
							},
						},
						{
							Vertices: [3]engine.Vertex{
								{Position: &minemath.Vec3{.5, -.5, .5}, Color: engine.WHITE},
								{Position: &minemath.Vec3{.5, .5, -.5}, Color: engine.WHITE},
								{Position: &minemath.Vec3{.5, .5, .5}, Color: engine.WHITE},
							},
						},
						// TOP
						{
							Vertices: [3]engine.Vertex{
								{Position: &minemath.Vec3{-.5, .5, .5}, Color: engine.GREEN},
								{Position: &minemath.Vec3{.5, .5, .5}, Color: engine.GREEN},
								{Position: &minemath.Vec3{.5, .5, -.5}, Color: engine.GREEN},
							},
						},
						{
							Vertices: [3]engine.Vertex{
								{Position: &minemath.Vec3{-.5, .5, .5}, Color: engine.GREEN},
								{Position: &minemath.Vec3{.5, .5, -.5}, Color: engine.GREEN},
								{Position: &minemath.Vec3{-.5, .5, -.5}, Color: engine.GREEN},
							},
						},
						// BOTTOM
						{
							Vertices: [3]engine.Vertex{
								{Position: &minemath.Vec3{-.5, -.5, .5}, Color: engine.WHITE},
								{Position: &minemath.Vec3{.5, -.5, .5}, Color: engine.WHITE},
								{Position: &minemath.Vec3{.5, -.5, -.5}, Color: engine.WHITE},
							},
						},
						{
							Vertices: [3]engine.Vertex{
								{Position: &minemath.Vec3{-.5, -.5, .5}, Color: engine.WHITE},
								{Position: &minemath.Vec3{.5, -.5, -.5}, Color: engine.WHITE},
								{Position: &minemath.Vec3{-.5, -.5, -.5}, Color: engine.WHITE},
							},
						},
					},
				},
			},
			at,
		),
	}
}
