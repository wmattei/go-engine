package minemath

import "math"

type Vec3 [3]float32
type Vec2 [2]float32

func (v Vec3) Negate() Vec3 {
	return Vec3{-v.X(), -v.Y(), -v.Z()}
}

func (v Vec3) Len() float32 {
	return float32(math.Sqrt(float64(v.X()*v.X() + v.Y()*v.Y() + v.Z()*v.Z())))
}

func (v Vec3) Mul(scalar float32) Vec3 {
	return Vec3{v.X() * scalar, v.Y() * scalar, v.Z() * scalar}
}

func (v Vec3) Dot(other Vec3) float32 {
	return v.X()*other.X() + v.Y()*other.Y() + v.Z()*other.Z()
}

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

func (m *Mat4) Flatten() [16]float32 {
	var result [16]float32
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			result[i*4+j] = m[j][i]
		}
	}
	return result
}
func DegreesToRadians(degrees float32) float32 {
	return degrees * (math.Pi / 180)
}

func TransformVec3(matrix Mat4, vec Vec3) Vec3 {
	x := matrix[0][0]*vec.X() + matrix[0][1]*vec.Y() + matrix[0][2]*vec.Z() + matrix[0][3]
	y := matrix[1][0]*vec.X() + matrix[1][1]*vec.Y() + matrix[1][2]*vec.Z() + matrix[1][3]
	z := matrix[2][0]*vec.X() + matrix[2][1]*vec.Y() + matrix[2][2]*vec.Z() + matrix[2][3]
	w := matrix[3][0]*vec.X() + matrix[3][1]*vec.Y() + matrix[3][2]*vec.Z() + matrix[3][3]

	if w != 0.0 {
		return Vec3{x / w, y / w, z / w}
	}
	return Vec3{x, y, z}
}

func MultiplyMatrices(a, b Mat4) Mat4 {
	var result Mat4
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			result[i][j] = a[i][0]*b[0][j] + a[i][1]*b[1][j] + a[i][2]*b[2][j] + a[i][3]*b[3][j]
		}
	}
	return result
}

func CalculateLightIntensity(normal, lightDir Vec3) float32 {
	intensity := Dot(normal, lightDir)
	if intensity < 0.4 {
		intensity = 0.4
	}
	return intensity
}
