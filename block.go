package main

type Face struct {
	Visible bool
	Mesh    *Mesh
}

type Block struct {
	Faces [6]*Face
	Color *Color
}

const (
	Right = iota
	Left
	Top
	Bottom
	Front
	Back
)

func (b *Block) Update(w *Chunk) {
	// position := b.Position
	// neighbors := [6][3]float32{
	// 	{position.X() + 1, position.Y(), position.Z()}, // right
	// 	{position.X() - 1, position.Y(), position.Z()}, // left
	// 	{position.X(), position.Y() + 1, position.Z()}, // top
	// 	{position.X(), position.Y() - 1, position.Z()}, // bottom
	// 	{position.X(), position.Y(), position.Z() + 1}, // front
	// 	{position.X(), position.Y(), position.Z() - 1}, // back
	// }
	// for i, n := range neighbors {
	// 	b.Faces[i].Visible = w.At(int(n[0]), int(n[1]), int(n[2])) != nil
	// }
}
func (b *Block) getFaceVerticesAndIndices(x, y, z int, direction int, indexOffset uint32) ([]float32, []uint32) {
	var faceVertices []float32
	var faceIndices []uint32

	color := b.Color.ToVec4()

	switch direction {
	case Right:
		faceVertices = []float32{
			float32(x + 1), float32(y), float32(z), color[0], color[1], color[2], color[3], // Bottom-right
			float32(x + 1), float32(y + 1), float32(z), color[0], color[1], color[2], color[3], // Top-right
			float32(x + 1), float32(y + 1), float32(z + 1), color[0], color[1], color[2], color[3], // Top-left
			float32(x + 1), float32(y), float32(z + 1), color[0], color[1], color[2], color[3], // Bottom-left
		}
	case Left:
		faceVertices = []float32{
			float32(x), float32(y), float32(z + 1), color[0], color[1], color[2], color[3], // Bottom-right
			float32(x), float32(y + 1), float32(z + 1), color[0], color[1], color[2], color[3], // Top-right
			float32(x), float32(y + 1), float32(z), color[0], color[1], color[2], color[3], // Top-left
			float32(x), float32(y), float32(z), color[0], color[1], color[2], color[3], // Bottom-left
		}
	case Top:
		faceVertices = []float32{
			float32(x), float32(y + 1), float32(z), color[0], color[1], color[2], color[3], // Bottom-left
			float32(x + 1), float32(y + 1), float32(z), color[0], color[1], color[2], color[3], // Bottom-right
			float32(x + 1), float32(y + 1), float32(z + 1), color[0], color[1], color[2], color[3], // Top-right
			float32(x), float32(y + 1), float32(z + 1), color[0], color[1], color[2], color[3], // Top-left
		}
	case Bottom:
		faceVertices = []float32{
			float32(x), float32(y), float32(z), color[0], color[1], color[2], color[3], // Bottom-left
			float32(x + 1), float32(y), float32(z), color[0], color[1], color[2], color[3], // Bottom-right
			float32(x + 1), float32(y), float32(z + 1), color[0], color[1], color[2], color[3], // Top-right
			float32(x), float32(y), float32(z + 1), color[0], color[1], color[2], color[3], // Top-left
		}
	case Front:
		faceVertices = []float32{
			float32(x), float32(y), float32(z + 1), color[0], color[1], color[2], color[3], // Bottom-left
			float32(x + 1), float32(y), float32(z + 1), color[0], color[1], color[2], color[3], // Bottom-right
			float32(x + 1), float32(y + 1), float32(z + 1), color[0], color[1], color[2], color[3], // Top-right
			float32(x), float32(y + 1), float32(z + 1), color[0], color[1], color[2], color[3], // Top-left
		}
	case Back:
		faceVertices = []float32{
			float32(x), float32(y), float32(z), color[0], color[1], color[2], color[3], // Bottom-left
			float32(x + 1), float32(y), float32(z), color[0], color[1], color[2], color[3], // Bottom-right
			float32(x + 1), float32(y + 1), float32(z), color[0], color[1], color[2], color[3], // Top-right
			float32(x), float32(y + 1), float32(z), color[0], color[1], color[2], color[3], // Top-left
		}
	}

	faceIndices = []uint32{
		indexOffset, indexOffset + 1, indexOffset + 2,
		indexOffset, indexOffset + 2, indexOffset + 3,
	}

	return faceVertices, faceIndices
}
