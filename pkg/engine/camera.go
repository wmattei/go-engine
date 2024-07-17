package engine

import (
	minemath "github.com/wmattei/minceraft/math"
)

type Camera interface {
	GetProjectionMatrix() minemath.Mat4
	GetViewMatrix() minemath.Mat4
	// ProcessMouseMovement(xoffset, yoffset float32)
	Move(x, y, z float32)
	Rotate(pitch, yaw float32)
	ProcessKeyboard(direction string)
}
