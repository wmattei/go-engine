package engine

import (
	"math"

	minemath "github.com/wmattei/minceraft/math"
)

type Plane struct {
	Normal   minemath.Vec3
	Distance float32
}

func (p *Plane) DistanceToPoint(point minemath.Vec3) float32 {
	return p.Normal[0]*point[0] + p.Normal[1]*point[1] + p.Normal[2]*point[2] + p.Distance
}

func (p *Plane) Normalize() {
	length := float32(math.Sqrt(float64(p.Normal[0]*p.Normal[0] +
		p.Normal[1]*p.Normal[1] +
		p.Normal[2]*p.Normal[2])))
	p.Normal[0] /= length
	p.Normal[1] /= length
	p.Normal[2] /= length
	p.Distance /= length
}

type Frustum struct {
	topPlane    *Plane
	bottomPlane *Plane
	leftPlane   *Plane
	rightPlane  *Plane
	nearPlane   *Plane
	farPlane    *Plane

	camera *PerspectiveCamera
}

func (f *Frustum) GetPlanes() []Plane {
	return []Plane{*f.topPlane, *f.bottomPlane, *f.leftPlane, *f.rightPlane, *f.nearPlane, *f.farPlane}
}

func NewFrustum(camera *PerspectiveCamera) *Frustum {
	frustum := &Frustum{}
	frustum.camera = camera

	return frustum
}

func (f *Frustum) UpdateFrustum(viewProjMatrix minemath.Mat4) {
	var planes [6]Plane

	// Left plane
	planes[0] = Plane{
		Normal: [3]float32{
			viewProjMatrix[3][0] + viewProjMatrix[0][0],
			viewProjMatrix[3][1] + viewProjMatrix[0][1],
			viewProjMatrix[3][2] + viewProjMatrix[0][2],
		},
		Distance: viewProjMatrix[3][3] + viewProjMatrix[0][3],
	}

	// Right plane
	planes[1] = Plane{
		Normal: [3]float32{
			viewProjMatrix[3][0] - viewProjMatrix[0][0],
			viewProjMatrix[3][1] - viewProjMatrix[0][1],
			viewProjMatrix[3][2] - viewProjMatrix[0][2],
		},
		Distance: viewProjMatrix[3][3] - viewProjMatrix[0][3],
	}

	// Bottom plane
	planes[2] = Plane{
		Normal: [3]float32{
			viewProjMatrix[3][0] + viewProjMatrix[1][0],
			viewProjMatrix[3][1] + viewProjMatrix[1][1],
			viewProjMatrix[3][2] + viewProjMatrix[1][2],
		},
		Distance: viewProjMatrix[3][3] + viewProjMatrix[1][3],
	}

	// Top plane
	planes[3] = Plane{
		Normal: [3]float32{
			viewProjMatrix[3][0] - viewProjMatrix[1][0],
			viewProjMatrix[3][1] - viewProjMatrix[1][1],
			viewProjMatrix[3][2] - viewProjMatrix[1][2],
		},
		Distance: viewProjMatrix[3][3] - viewProjMatrix[1][3],
	}

	// Near plane
	planes[4] = Plane{
		Normal: [3]float32{
			viewProjMatrix[3][0] + viewProjMatrix[2][0],
			viewProjMatrix[3][1] + viewProjMatrix[2][1],
			viewProjMatrix[3][2] + viewProjMatrix[2][2],
		},
		Distance: viewProjMatrix[3][3] + viewProjMatrix[2][3],
	}

	// Far plane
	planes[5] = Plane{
		Normal: [3]float32{
			viewProjMatrix[3][0] - viewProjMatrix[2][0],
			viewProjMatrix[3][1] - viewProjMatrix[2][1],
			viewProjMatrix[3][2] - viewProjMatrix[2][2],
		},
		Distance: viewProjMatrix[3][3] - viewProjMatrix[2][3],
	}
	// // Normalize the planes
	for i := 0; i < 6; i++ {
		planes[i].Normalize()
	}

	f.leftPlane = &planes[0]
	f.rightPlane = &planes[1]
	f.bottomPlane = &planes[2]
	f.topPlane = &planes[3]
	f.nearPlane = &planes[4]
	f.farPlane = &planes[5]

	// fmt.Println(f.leftPlane)

}
