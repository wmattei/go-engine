package engine

type Game interface {
	Update(dt float32)
	Render(scene *Scene)
	GetScreenSize() (int, int)
}
