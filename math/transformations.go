package minemath

import (
	"math"
)

func GetPerspectiveProjectionMatrix(fov, aspect, near, far float32) Mat4 {
	f := 1.0 / float32(math.Tan(float64(fov/2)))
	nf := 1 / (near - far)

	return Mat4{
		{f / aspect, 0, 0, 0},
		{0, f, 0, 0},
		{0, 0, (far + near) * nf, (2 * far * near) * nf},
		{0, 0, -1, 0},
	}
}

func GetIdentityMatrix() Mat4 {
	return Mat4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func GetTranslationMatrix(x, y, z float32) Mat4 {
	return Mat4{
		{1, 0, 0, x},
		{0, 1, 0, y},
		{0, 0, 1, z},
		{0, 0, 0, 1},
	}
}
