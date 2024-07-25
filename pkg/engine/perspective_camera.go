package engine

import (
	"math"

	minemath "github.com/wmattei/minceraft/math"
)

type PerspectiveCamera struct {
	Position *minemath.Vec3
	front    minemath.Vec3
	worldUp  minemath.Vec3

	yaw   float32
	pitch float32

	fov    float32
	aspect float32
	near   float32
	far    float32
}

func NewPerspectiveCamera(position, up minemath.Vec3, yaw, pitch, fov, aspect, near, far float32) *PerspectiveCamera {
	camera := &PerspectiveCamera{
		Position: &position,
		worldUp:  up,
		yaw:      yaw,
		pitch:    pitch,
		fov:      fov,
		aspect:   aspect,
		near:     near,
		far:      far,
		front:    minemath.Vec3{0, 0, -1},
	}
	camera.updateCameraVectors()

	return camera
}

func (cam *PerspectiveCamera) GetViewMatrix() minemath.Mat4 {
	return minemath.LookAt(*cam.Position, minemath.Add(*cam.Position, cam.front), cam.worldUp)
}

func (cam *PerspectiveCamera) GetProjectionMatrix() minemath.Mat4 {
	return minemath.GetPerspectiveProjectionMatrix(cam.fov, cam.aspect, cam.near, cam.far)
}

func (cam *PerspectiveCamera) Move(x, y, z float32) {
	cam.Position[0] += x
	cam.Position[1] += y
	cam.Position[2] += z
	cam.updateCameraVectors()
}

func (cam *PerspectiveCamera) Rotate(pitch, yaw float32) {
	cam.yaw += yaw
	cam.pitch += pitch

	if cam.pitch > 89.0 {
		cam.pitch = 89.0
	}
	if cam.pitch < -89.0 {
		cam.pitch = -89.0
	}

	cam.updateCameraVectors()
}

func (cam *PerspectiveCamera) updateCameraVectors() {
	front := minemath.Vec3{
		float32(math.Cos(float64(minemath.DegreesToRadians(cam.pitch))) * math.Cos(float64(minemath.DegreesToRadians(cam.yaw)))),
		float32(math.Sin(float64(minemath.DegreesToRadians((cam.pitch))))),
		float32(math.Cos(float64(minemath.DegreesToRadians(cam.pitch))) * math.Sin(float64(minemath.DegreesToRadians(cam.yaw)))),
	}

	cam.front = front
}

func (cam *PerspectiveCamera) ProcessKeyboard(direction string, dt float32) {
	velocity := float32(dt * 10)

	switch direction {
	case "FORWARD":
		forward := minemath.Vec3{cam.front.X(), 0, cam.front.Z()}
		forward = minemath.Normalize(forward)
		cam.Position[0] += forward.X() * velocity
		cam.Position[2] += forward.Z() * velocity

	case "BACKWARD":
		backward := minemath.Vec3{cam.front.X(), 0, cam.front.Z()}
		backward = minemath.Normalize(backward)
		cam.Position[0] -= backward.X() * velocity
		cam.Position[2] -= backward.Z() * velocity

	case "LEFT":
		right := minemath.Normalize(minemath.Cross(cam.front, cam.worldUp))
		left := minemath.Vec3{right.X(), 0, right.Z()}
		left = minemath.Normalize(left)
		cam.Position[0] -= left.X() * velocity
		cam.Position[2] -= left.Z() * velocity

	case "RIGHT":
		right := minemath.Normalize(minemath.Cross(cam.front, cam.worldUp))
		newRight := minemath.Vec3{right.X(), 0, right.Z()}
		newRight = minemath.Normalize(newRight)
		cam.Position[0] += newRight.X() * velocity
		cam.Position[2] += newRight.Z() * velocity

	case "UP":
		cam.Position[1] += velocity

	case "DOWN":
		cam.Position[1] -= velocity
	}

	cam.updateCameraVectors()
}
