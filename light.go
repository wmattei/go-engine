package main

import minemath "github.com/wmattei/minceraft/math"

type Light struct {
	Direction *minemath.Vec3
}

func (l *Light) Move() {
	l.Direction[0] *= 0.1
}
