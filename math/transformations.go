package minemath

import "math"

func GetRotateXMatrix(angle float32) Mat3 {
	cos := float32(math.Cos(float64(angle)))
	sin := float32(math.Sin(float64(angle)))
	return Mat3{
		{1, 0, 0},
		{0, cos, -sin},
		{0, sin, cos},
	}
}

func GetRotateYMatrix(angle float32) Mat3 {
	cos := float32(math.Cos(float64(angle)))
	sin := float32(math.Sin(float64(angle)))
	return Mat3{
		{cos, 0, sin},
		{0, 1, 0},
		{-sin, 0, cos},
	}
}

func GetRotateZMatrix(angle float32) Mat3 {
	cos := float32(math.Cos(float64(angle)))
	sin := float32(math.Sin(float64(angle)))
	return Mat3{
		{cos, -sin, 0},
		{sin, cos, 0},
		{0, 0, 1},
	}
}

func GetOrthoGraphicProjectionMatrix(left, right, bottom, top, near, far float32) Mat4 {
	return Mat4{
		{2 / (right - left), 0, 0, 0},
		{0, 2 / (top - bottom), 0, 0},
		{0, 0, -2 / (far - near), 0},
		{-(right + left) / (right - left), -(top + bottom) / (top - bottom), -(far + near) / (far - near), 1},
	}
}

func GetPerspectiveProjectionMatrix(fov, aspect, near, far float32) Mat4 {
	f := 1.0 / float32(math.Tan(float64(fov)/2.0))
	nf := 1 / (near - far)

	return Mat4{
		{f / aspect, 0, 0, 0},
		{0, f, 0, 0},
		{0, 0, (far + near) * nf, -1},
		{0, 0, (2 * far * near) * nf, 0},
	}
}
