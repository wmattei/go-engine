package minemath

import "math"

type Vec3 [3]float32
type Vec2 [2]float32

func (v *Vec3) X() float32 {
	return v[0]
}
func (v *Vec3) Y() float32 {
	return v[1]
}
func (v *Vec3) Z() float32 {
	return v[2]
}

type Vec4 [4]float32

type Mat3 [3][3]float32
type Mat4 [4][4]float32

func MultiplyMatrices(a, b Mat4) Mat4 {
	var result Mat4
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			result[i][j] = a[i][0]*b[0][j] + a[i][1]*b[1][j] + a[i][2]*b[2][j] + a[i][3]*b[3][j]
		}
	}
	return result
}

func SphericalToCartesian(radius, yaw, pitch float32) [3]float32 {
	yawRad := yaw * (math.Pi / 180)
	pitchRad := pitch * (math.Pi / 180)
	x := radius * float32(math.Cos(float64(pitchRad))*math.Cos(float64(yawRad)))
	y := radius * float32(math.Sin(float64(pitchRad)))
	z := radius * float32(math.Cos(float64(pitchRad))*math.Sin(float64(yawRad)))
	return [3]float32{x, y, z}
}

func (m *Mat4) Flatten() [16]float32 {
	var result [16]float32
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			result[i*4+j] = m[i][j]
		}
	}
	return result
}

func DegreesToRadians(degrees float32) float32 {
	return degrees * (math.Pi / 180)
}
