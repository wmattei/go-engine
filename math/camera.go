package minemath

import (
	"math"
)

func LookAt(eye, center, up Vec3) Mat4 {
	f := Normalize(Subtract(center, eye))
	s := Normalize(Cross(f, up))
	u := Cross(s, f)

	return Mat4{
		{s.X(), u.X(), -f.X(), 0},
		{s.Y(), u.Y(), -f.Y(), 0},
		{s.Z(), u.Z(), -f.Z(), 0},
		{-Dot(s, eye), -Dot(u, eye), Dot(f, eye), 1},
	}
}

func Normalize(v Vec3) Vec3 {
	length := float32(math.Sqrt(float64(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])))
	if length == 0 {
		return Vec3{0, 0, 0}
	}
	return Vec3{v[0] / length, v[1] / length, v[2] / length}

}

func Subtract(a, b Vec3) Vec3 {
	return [3]float32{a[0] - b[0], a[1] - b[1], a[2] - b[2]}
}

func Cross(a, b Vec3) Vec3 {
	return [3]float32{
		a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - a[0]*b[2],
		a[0]*b[1] - a[1]*b[0],
	}
}

func Dot(a, b Vec3) float32 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

func Add(a, b Vec3) Vec3 {
	return [3]float32{a[0] + b[0], a[1] + b[1], a[2] + b[2]}
}

func Scale(v Vec3, s float32) Vec3 {
	return [3]float32{v[0] * s, v[1] * s, v[2] * s}
}
