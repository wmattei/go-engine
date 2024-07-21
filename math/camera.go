package minemath

func LookAt(eye, center, up Vec3) Mat4 {
	zaxis := Normalize(Subtract(eye, center)) // The forward vector
	xaxis := Normalize(Cross(up, zaxis))      // The right vector
	yaxis := Cross(zaxis, xaxis)              // The up vector

	return Mat4{
		{xaxis[0], xaxis[1], xaxis[2], -Dot(xaxis, eye)},
		{yaxis[0], yaxis[1], yaxis[2], -Dot(yaxis, eye)},
		{zaxis[0], zaxis[1], zaxis[2], -Dot(zaxis, eye)},
		{0, 0, 0, 1},
	}
}

func Normalize(v Vec3) Vec3 {
	length := v.Len()
	return Vec3{v.X() / length, v.Y() / length, v.Z() / length}
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
