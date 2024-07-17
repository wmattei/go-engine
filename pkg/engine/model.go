package engine

import (
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/google/uuid"
	minemath "github.com/wmattei/minceraft/math"
)

type Model struct {
	ID          string
	Meshes      []*Mesh
	Position    minemath.Vec3
	Rotation    *minemath.Vec3
	Scale       minemath.Vec3
	NeedsRender bool

	vao            uint32
	trianglesCount int
}

func NewModel(meshes []*Mesh, position minemath.Vec3) *Model {
	id := randomID()
	return &Model{
		ID:          id.String(),
		Meshes:      meshes,
		Position:    position,
		Rotation:    &minemath.Vec3{0, 0, 0},
		Scale:       minemath.Vec3{1, 1, 1},
		NeedsRender: false,
	}
}

func (m *Model) SetVao() {
	m.vao = getVao(m.GetTriangles())
	m.trianglesCount = len(m.GetTriangles()) * 3
}

func (m *Model) GetTriangles() []*Triangle {
	var triangles []*Triangle
	for _, mesh := range m.Meshes {
		triangles = append(triangles, mesh.Triangles...)
	}
	return triangles

}

func (m *Model) GetModelMatrix() minemath.Mat4 {
	translation := minemath.Mat4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{m.Position.X(), m.Position.Y(), m.Position.Z(), 1},
	}

	angleX := m.Rotation.X() * (math.Pi / 180.0)
	cosX := float32(math.Cos(float64(angleX)))
	sinX := float32(math.Sin(float64(angleX)))
	rotationX := minemath.Mat4{
		{1, 0, 0, 0},
		{0, cosX, -sinX, 0},
		{0, sinX, cosX, 0},
		{0, 0, 0, 1},
	}

	angleY := m.Rotation.Y() * (math.Pi / 180.0)
	cosY := float32(math.Cos(float64(angleY)))
	sinY := float32(math.Sin(float64(angleY)))
	rotationY := minemath.Mat4{
		{cosY, 0, sinY, 0},
		{0, 1, 0, 0},
		{-sinY, 0, cosY, 0},
		{0, 0, 0, 1},
	}

	angleZ := m.Rotation.Z() * (math.Pi / 180.0)
	cosZ := float32(math.Cos(float64(angleZ)))
	sinZ := float32(math.Sin(float64(angleZ)))
	rotationZ := minemath.Mat4{
		{cosZ, -sinZ, 0, 0},
		{sinZ, cosZ, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}

	scale := minemath.Mat4{
		{m.Scale.X(), 0, 0, 0},
		{0, m.Scale.Y(), 0, 0},
		{0, 0, m.Scale.Z(), 0},
		{0, 0, 0, 1},
	}

	rotation := minemath.MultiplyMatrices(rotationX, minemath.MultiplyMatrices(rotationY, rotationZ))
	modelMatrix := minemath.MultiplyMatrices(translation, minemath.MultiplyMatrices(rotation, scale))
	return modelMatrix
}

func randomID() uuid.UUID {
	return uuid.New()
}

func getVao(triangles []*Triangle) uint32 {
	var vao, vbo, ebo uint32

	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)

	gl.BindVertexArray(vao)

	vertices, indices := unpackTriangles(triangles)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	// Position attribute
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 7*4, gl.PtrOffset(0))

	// Color attribute
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 7*4, gl.PtrOffset(3*4))

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	return vao
}

func unpackTriangles(triangles []*Triangle) ([]float32, []uint32) {
	var vertices []float32
	var indices []uint32
	var index uint32

	for _, triangle := range triangles {
		for _, vertex := range triangle.Vertices {
			vertices = append(vertices,
				vertex.Position.X(), vertex.Position.Y(), vertex.Position.Z(),
				float32(vertex.Color.R)/255, float32(vertex.Color.G)/255, float32(vertex.Color.B)/255, 1.0,
			)
			indices = append(indices, index)
			index++
		}
	}

	return vertices, indices
}
